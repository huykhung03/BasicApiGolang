package sqlc

import (
	"context"
)

type PurchaseTransactionPagrams struct {
	Product Product `json:"product"`
	Buyer   string  `json:"username_buyer"`
}

type PurchaseTransactionResult struct {
	PurchaseHistory     PurchaseHistory `json:"purchase_history"`
	BankAccountOfBuyer  BankAccount     `json:"bank_acount_of_buyer"`
	BankAccountOfSeller BankAccount     `json:"bank_acount_of_seller"`
}

// PurchaseTransaction mades a purchase from Buyer to Seller
// Step 1: Create the purchase history
// Step 2: Deduct money from account buyer
// Step 3: Add money from account seller
func (store *SQLStore) PurchaseTransaction(ctx context.Context, arg PurchaseTransactionPagrams) (
	PurchaseTransactionResult, error) {
	var result PurchaseTransactionResult

	err := store.execTx(context.Background(), func(q *Queries) error {
		// * get bank account of buyer with username and currency
		bankAccountOfBuyer, err := q.GetBankAccountByUserNameAndCurrencyForUpdate(context.Background(),
			GetBankAccountByUserNameAndCurrencyForUpdateParams{
				Username: arg.Buyer,
				Currency: "USD",
			})
		if err != nil {
			return err
		}

		// * step 1
		result.PurchaseHistory, err = q.CreatePuschaseHistory(context.Background(),
			CreatePuschaseHistoryParams{
				IDProduct:         arg.Product.IDProduct,
				Buyer:             arg.Buyer,
				CardNumberOfBuyer: bankAccountOfBuyer.CardNumber,
			})
		if err != nil {
			return err
		}

		// * step 2
		result.BankAccountOfBuyer, err = q.AddBankAccountBalance(context.Background(),
			AddBankAccountBalanceParams{
				Currency:   bankAccountOfBuyer.Currency,
				Amount:     -arg.Product.Price,
				CardNumber: bankAccountOfBuyer.CardNumber,
			})
		if err != nil {
			return err
		}

		// * get bank account of seller with username and currency
		bankAccountOfSeller, err := q.GetBankAccountByUserNameAndCurrencyForUpdate(context.Background(),
			GetBankAccountByUserNameAndCurrencyForUpdateParams{
				Username: arg.Product.Owner,
				Currency: "USD",
			})
		if err != nil {
			return err
		}

		// * step 3
		result.BankAccountOfSeller, err = q.AddBankAccountBalance(context.Background(),
			AddBankAccountBalanceParams{
				Currency:   "USD",
				Amount:     arg.Product.Price,
				CardNumber: bankAccountOfSeller.CardNumber,
			})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
