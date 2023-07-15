package sqlc

import (
	"context"
)

type PurchaseTransactionPagrams struct {
	UsernameBuyer  string  `json:"username_buyer"`
	UsernameSeller string  `json:"username_seller"`
	Product        Product `json:"product"`
}

type ChangedBalance struct {
	Username       string `json:"username"`
	ChangedBalance int32  `json:"changed_balance"`
}

type PurchaseTransactionResult struct {
	PurchaseHistory         PurchaseHistory `json:"purchase_history"`
	BalanceOfUsernameBuyer  ChangedBalance  `json:"balance_of_username_buyer"`
	BalanceOfUsernameSeller ChangedBalance  `json:"balance_of_username_seller"`
}

// PurchaseTransaction mades a purchase from userA to userB
func (store *Store) PurchaseTransaction(ctx context.Context, arg PurchaseTransactionPagrams) (
	PurchaseTransactionResult, error) {
	var result PurchaseTransactionResult
	err := store.execTx(context.Background(), func(q *Queries) error {
		argCardNumber := GetCardNumberByUserNameAndCurrencyParams{
			Username: arg.UsernameBuyer,
			Currency: "USD",
		}

		bankAccountOfUsernameBuyer, err := q.GetCardNumberByUserNameAndCurrency(context.Background(), argCardNumber)
		if err != nil {
			return err
		}

		result.PurchaseHistory, err = q.CreatePuschaseHistory(context.Background(), CreatePuschaseHistoryParams{
			IDProduct:  arg.Product.IDProduct,
			Buyer:      arg.UsernameBuyer,
			CardNumber: bankAccountOfUsernameBuyer.CardNumber,
		})
		if err != nil {
			return err
		}

		product, err := q.GetProduct(context.Background(), result.PurchaseHistory.IDProduct)
		if err != nil {
			return err
		}

		result.BalanceOfUsernameBuyer = ChangedBalance{
			Username:       arg.UsernameBuyer,
			ChangedBalance: -product.Price,
		}

		result.BalanceOfUsernameSeller = ChangedBalance{
			Username:       product.Owner,
			ChangedBalance: product.Price,
		}

		return nil
	})
	return result, err
}
