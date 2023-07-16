package sqlc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPurchaseTransaction(t *testing.T) {
	store := NewStore(testBD)

	usernameBuyer := createRandomBankAccount(t)

	// * run n concurrent puschase transactions
	n := 5
	product := createRandomProduct(t)

	errs := make(chan error)
	results := make(chan PurchaseTransactionResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.PurchaseTransaction(context.Background(), PurchaseTransactionPagrams{
				Product:       product,
				UsernameBuyer: usernameBuyer.Username,
			})
			errs <- err
			results <- result
		}()
	}

	// * check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// * check purchase history
		purchaseHistory := result.PurchaseHistory
		require.NotEmpty(t, purchaseHistory.IDPurchaseHistory, result.PurchaseHistory.IDPurchaseHistory)
		require.Equal(t, purchaseHistory.IDProduct, product.IDProduct)
		require.Equal(t, purchaseHistory.Buyer, usernameBuyer.Username)
		require.Equal(t, purchaseHistory.CardNumber, usernameBuyer.CardNumber)
		require.NotZero(t, purchaseHistory.CreatedAt)

		_, err = store.GetPurchaseHistory(context.Background(), purchaseHistory.IDPurchaseHistory)
		require.NoError(t, err)

		// * check balance of seller and buyer
		balanceOfBuyer := result.BalanceOfUsernameBuyer
		require.NotEmpty(t, balanceOfBuyer)
		require.Equal(t, balanceOfBuyer.Username, usernameBuyer.Username)
		require.Equal(t, balanceOfBuyer.ChangedBalance, -product.Price)

		balanceOfSeller := result.BalanceOfUsernameSeller
		require.NotEmpty(t, balanceOfSeller)
		require.Equal(t, balanceOfSeller.Username, product.Owner)
		require.Equal(t, balanceOfSeller.ChangedBalance, product.Price)

	}

}
