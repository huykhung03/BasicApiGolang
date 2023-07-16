package sqlc

import (
	"context"
	"fmt"
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
func (store *Store) PurchaseTransaction(ctx context.Context, arg PurchaseTransactionPagrams) (
	PurchaseTransactionResult, error) {
	var result PurchaseTransactionResult

	err := store.execTx(context.Background(), func(q *Queries) error {

		txName := ctx.Value(txKey)

		// * get bank account of buyer with username and currency
		fmt.Println(txName, "get account buyer")
		bankAccountOfBuyer, err := q.GetBankAccountByUserNameAndCurrency(context.Background(),
			GetBankAccountByUserNameAndCurrencyParams{
				Username: arg.Buyer,
				Currency: "USD",
			})
		if err != nil {
			return err
		}

		// * step 1
		fmt.Println(txName, "create purchase history")
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
		fmt.Println(txName, "update account buyer")
		result.BankAccountOfBuyer, err = q.AddBankAccountBalance(context.Background(),
			AddBankAccountBalanceParams{
				Currency:   bankAccountOfBuyer.Currency,
				Balance:    bankAccountOfBuyer.Balance - arg.Product.Price,
				CardNumber: bankAccountOfBuyer.CardNumber,
			})
		if err != nil {
			return err
		}

		// * get bank account of seller with username and currency
		fmt.Println(txName, "get account seller")
		bankAccountOfSeller, err := q.GetBankAccountByUserNameAndCurrency(context.Background(),
			GetBankAccountByUserNameAndCurrencyParams{
				Username: arg.Product.Owner,
				Currency: "USD",
			})
		if err != nil {
			return err
		}

		// * step 3
		fmt.Println(txName, "update account seller")
		result.BankAccountOfSeller, err = q.AddBankAccountBalance(context.Background(),
			AddBankAccountBalanceParams{
				Currency:   "USD",
				Balance:    bankAccountOfSeller.Balance + arg.Product.Price,
				CardNumber: bankAccountOfSeller.CardNumber,
			})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
