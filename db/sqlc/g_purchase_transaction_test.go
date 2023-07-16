package sqlc

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPurchaseTransaction(t *testing.T) {
	store := NewStore(testBD)

	buyer := createRandomBankAccount(t)

	// * run n concurrent puschase transactions
	n := 5
	product := createRandomProduct(t)

	seller, err := store.GetBankAccountByUserNameAndCurrency(context.Background(),
		GetBankAccountByUserNameAndCurrencyParams{
			Username: product.Owner,
			Currency: "USD",
		})
	require.NoError(t, err)

	fmt.Println(">>Before: Buyer: ", buyer.Balance, " Seller: ", seller.Balance)

	errs := make(chan error)
	results := make(chan PurchaseTransactionResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.PurchaseTransaction(context.Background(), PurchaseTransactionPagrams{
				Product: product,
				Buyer:   buyer.Username,
				Amount:  product.Price,
			})
			errs <- err
			results <- result
		}()
	}

	// * check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// * check infomation of seller and buyer
		bankAccountOfBuyer := result.BankAccountOfBuyer
		require.NotEmpty(t, bankAccountOfBuyer)
		require.Equal(t, bankAccountOfBuyer.Username, buyer.Username)
		require.Equal(t, bankAccountOfBuyer.CardNumber, buyer.CardNumber)
		require.Equal(t, bankAccountOfBuyer.Currency, buyer.Currency)
		// require.Equal(t, bankAccountOfBuyer.Balance, buyer.Balance)

		bankAccountOfSeller := result.BankAccountOfSeller
		require.NotEmpty(t, bankAccountOfSeller)
		require.Equal(t, bankAccountOfSeller.Username, seller.Username)
		require.Equal(t, bankAccountOfSeller.CardNumber, seller.CardNumber)
		require.Equal(t, bankAccountOfSeller.Currency, seller.Currency)
		// require.Equal(t, bankAccountOfSeller.Balance, seller.Balance)

		fmt.Println("Buyer: ", bankAccountOfBuyer.Balance, " Seller: ", bankAccountOfSeller.Balance)

		// * check purchase history
		purchaseHistory := result.PurchaseHistory
		require.NotEmpty(t, purchaseHistory)
		require.Equal(t, purchaseHistory.IDProduct, product.IDProduct)
		require.Equal(t, purchaseHistory.Buyer, buyer.Username)
		require.Equal(t, purchaseHistory.CardNumberOfBuyer, bankAccountOfBuyer.CardNumber)
		require.NotZero(t, purchaseHistory.CreatedAt)

		_, err = store.GetPurchaseHistory(context.Background(), purchaseHistory.IDPurchaseHistory)
		require.NoError(t, err)

		// * check balance
		diff_1 := buyer.Balance - bankAccountOfBuyer.Balance
		diff_2 := bankAccountOfSeller.Balance - seller.Balance
		fmt.Println(diff_1, " ", diff_2)
		require.Equal(t, diff_1, diff_2)
		require.True(t, diff_1 > 0)
		require.True(t, diff_1%product.Price == 0)

		k := int(diff_1 / product.Price)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedBankAccountBuyer, err := store.GetBankAccountByUserNameAndCurrency(context.Background(),
		GetBankAccountByUserNameAndCurrencyParams{
			Username: buyer.Username,
			Currency: "USD",
		})
	require.NoError(t, err)

	updatedBankAccountSeller, err := store.GetBankAccountByUserNameAndCurrency(context.Background(),
		GetBankAccountByUserNameAndCurrencyParams{
			Username: seller.Username,
			Currency: "USD",
		})
	require.NoError(t, err)

	fmt.Println(">>After_1: Buyer: ", buyer.Balance, " Seller: ", seller.Balance)
	fmt.Println(">>After_2: Buyer: ", updatedBankAccountBuyer.Balance, " Seller: ", updatedBankAccountSeller.Balance)
	require.Equal(t, buyer.Balance-int32(n)*product.Price, updatedBankAccountBuyer.Balance)
	require.Equal(t, seller.Balance-int32(n)*product.Price, updatedBankAccountSeller.Balance)
}
