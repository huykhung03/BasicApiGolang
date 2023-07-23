package sqlc

import (
	"context"
	"errors"
	"fmt"
)

type PurchaseTransactionPagrams struct {
	Product           Product `json:"product"`
	PurchaseQuantity  uint16  `json:"purchase_quantity"`
	Buyer             string  `json:"buyer"`
	CardNumberOfBuyer string  `json:"card_number_of_buyer"`
}

type Entry struct {
	BankAccount   BankAccount `json:"bank_account"`
	AmountOfMoney int32       `json:"amount_of_money"`
}

type PurchaseTransactionResult struct {
	PurchaseHistory PurchaseHistory `json:"purchase_history"`
	EntryOfBuyer    Entry           `json:"entry_of_buyer"`
	EntryOfSeller   Entry           `json:"entry_of_seller"`
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
				Currency: arg.Product.Currency,
			})
		if err != nil {
			return errors.New("The customer does not have a bank account that matches the currency of the goods")
		}

		if arg.CardNumberOfBuyer != bankAccountOfBuyer.CardNumber {
			return errors.New("The currency of the bank account does not match the currency of the product")
		}

		productBeforePurchase, err := q.UpdateQuantityOfProduct(context.Background(),
			UpdateQuantityOfProductParams{
				Amount:    0,
				IDProduct: arg.Product.IDProduct,
			})
		if err != nil {
			return err
		}

		if productBeforePurchase.Quantity < 0 {
			return errors.New("Goods are no longer available")
		}

		if arg.Product.Quantity < int32(arg.PurchaseQuantity) {
			return errors.New("The goods are not in the quantity you need")
		}

		if bankAccountOfBuyer.Balance < arg.Product.Price*int32(arg.PurchaseQuantity) {
			return errors.New("Customer's balance is not enough")
		}

		if arg.PurchaseQuantity > arg.PurchaseQuantity {
			return errors.New("The quantity of goods is not enough")
		}

		fmt.Println("Quantity of product: ", productBeforePurchase.Quantity, " Quantity of purchase: ", arg.PurchaseQuantity)

		productAfterPurchase, err := q.UpdateQuantityOfProduct(context.Background(),
			UpdateQuantityOfProductParams{
				Amount:    int32(arg.PurchaseQuantity),
				IDProduct: arg.Product.IDProduct,
			})
		if err != nil {
			return err
		}

		fmt.Println("Quantity of product after purchase: ", productAfterPurchase.Quantity)

		// * step 1
		result.PurchaseHistory, err = q.CreatePurchaseHistory(context.Background(),
			CreatePurchaseHistoryParams{
				IDProduct:         arg.Product.IDProduct,
				Buyer:             arg.Buyer,
				CardNumberOfBuyer: bankAccountOfBuyer.CardNumber,
			})
		if err != nil {
			return err
		}
		// *

		bankAccountOfBuyerAfterPurchase, err := q.AddBankAccountBalance(context.Background(),
			AddBankAccountBalanceParams{
				Currency:   bankAccountOfBuyer.Currency,
				Amount:     -arg.Product.Price * int32(arg.PurchaseQuantity),
				CardNumber: bankAccountOfBuyer.CardNumber,
			})
		if err != nil {
			return err
		}

		// * step 2
		result.EntryOfBuyer = Entry{
			BankAccount:   bankAccountOfBuyerAfterPurchase,
			AmountOfMoney: -arg.Product.Price * int32(arg.PurchaseQuantity),
		}
		if err != nil {
			return err
		}
		// *

		// * get bank account of seller with username and currency
		bankAccountOfSeller, err := q.GetBankAccountByUserNameAndCurrencyForUpdate(context.Background(),
			GetBankAccountByUserNameAndCurrencyForUpdateParams{
				Username: arg.Product.Owner,
				Currency: arg.Product.Currency,
			})
		if err != nil {
			return err
		}

		bankAccountOfSellerAfterPurchase, err := q.AddBankAccountBalance(context.Background(),
			AddBankAccountBalanceParams{
				Currency:   bankAccountOfSeller.Currency,
				Amount:     arg.Product.Price * int32(arg.PurchaseQuantity),
				CardNumber: bankAccountOfSeller.CardNumber,
			})
		if err != nil {
			return err
		}

		// * step 3
		result.EntryOfSeller = Entry{
			BankAccount:   bankAccountOfSellerAfterPurchase,
			AmountOfMoney: arg.Product.Price * int32(arg.PurchaseQuantity),
		}
		// *

		return nil
	})

	return result, err
}

// func (store *SQLStore) PurchaseTransaction(ctx context.Context, arg PurchaseTransactionPagrams) (
// 	PurchaseTransactionResult, error) {
// 	var result PurchaseTransactionResult

// 	err := store.execTx(context.Background(), func(q *Queries) error {
// 		// * get bank account of buyer with username and currency
// 		bankAccountOfBuyer, err := q.GetBankAccountByUserNameAndCurrencyForUpdate(context.Background(),
// 			GetBankAccountByUserNameAndCurrencyForUpdateParams{
// 				Username: arg.Buyer,
// 				Currency: "USD",
// 			})
// 		if err != nil {
// 			return err
// 		}

// 		// * step 1
// 		result.PurchaseHistory, err = q.CreatePuschaseHistory(context.Background(),
// 			CreatePuschaseHistoryParams{
// 				IDProduct:         arg.Product.IDProduct,
// 				Buyer:             arg.Buyer,
// 				CardNumberOfBuyer: bankAccountOfBuyer.CardNumber,
// 			})
// 		if err != nil {
// 			return err
// 		}

// 		// * step 2
// 		result.BankAccountOfBuyer, err = q.AddBankAccountBalance(context.Background(),
// 			AddBankAccountBalanceParams{
// 				Currency:   bankAccountOfBuyer.Currency,
// 				Amount:     -arg.Product.Price,
// 				CardNumber: bankAccountOfBuyer.CardNumber,
// 			})
// 		if err != nil {
// 			return err
// 		}

// 		// * get bank account of seller with username and currency
// 		bankAccountOfSeller, err := q.GetBankAccountByUserNameAndCurrencyForUpdate(context.Background(),
// 			GetBankAccountByUserNameAndCurrencyForUpdateParams{
// 				Username: arg.Product.Owner,
// 				Currency: "USD",
// 			})
// 		if err != nil {
// 			return err
// 		}

// 		// * step 3
// 		result.BankAccountOfSeller, err = q.AddBankAccountBalance(context.Background(),
// 			AddBankAccountBalanceParams{
// 				Currency:   "USD",
// 				Amount:     arg.Product.Price,
// 				CardNumber: bankAccountOfSeller.CardNumber,
// 			})
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})

// 	return result, err
// }
