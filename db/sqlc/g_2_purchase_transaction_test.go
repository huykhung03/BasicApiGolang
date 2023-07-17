package sqlc

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPurchaseTransactionDeadLock(t *testing.T) {
	store := NewStore(testBD)

	productOfBuyer := createRandomProduct(t)

	bankAccountOfBuyer, err := store.GetBankAccountByUserNameAndCurrency(context.Background(),
		GetBankAccountByUserNameAndCurrencyParams{
			Username: productOfBuyer.Owner,
			Currency: "USD",
		})
	require.NoError(t, err)

	// * run n concurrent puschase transactions
	n := 10
	productOfSeller := createRandomProduct(t)

	bankAccountOfSeller, err := store.GetBankAccountByUserNameAndCurrency(context.Background(),
		GetBankAccountByUserNameAndCurrencyParams{
			Username: productOfSeller.Owner,
			Currency: "USD",
		})
	require.NoError(t, err)

	fmt.Println(">>Before: Buyer: ", bankAccountOfBuyer.Balance, " Seller: ", bankAccountOfSeller.Balance)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		buyerTmp := productOfBuyer.Owner
		product := productOfBuyer

		if i%2 == 1 {
			buyerTmp = bankAccountOfSeller.Username
			product = productOfSeller
		}

		go func() {
			_, err := store.PurchaseTransaction(context.Background(), PurchaseTransactionPagrams{
				Product: product,
				Buyer:   buyerTmp,
			})
			errs <- err
		}()
	}

	// * check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedBankAccountBuyer, err := store.GetBankAccountByUserNameAndCurrency(context.Background(),
		GetBankAccountByUserNameAndCurrencyParams{
			Username: productOfBuyer.Owner,
			Currency: "USD",
		})
	require.NoError(t, err)

	updatedBankAccountSeller, err := store.GetBankAccountByUserNameAndCurrency(context.Background(),
		GetBankAccountByUserNameAndCurrencyParams{
			Username: bankAccountOfSeller.Username,
			Currency: "USD",
		})
	require.NoError(t, err)

	fmt.Println(">>After_1: Buyer: ", bankAccountOfBuyer.Balance, " Seller: ", bankAccountOfSeller.Balance)
	fmt.Println(">>After_2: Buyer: ", updatedBankAccountBuyer.Balance, " Seller: ", updatedBankAccountSeller.Balance)
	require.Equal(t, bankAccountOfBuyer.Balance, updatedBankAccountBuyer.Balance)
	require.Equal(t, bankAccountOfSeller.Balance, updatedBankAccountSeller.Balance)
}
