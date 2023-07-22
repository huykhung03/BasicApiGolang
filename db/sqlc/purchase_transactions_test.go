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

	// var purchaseQuantity uint16 = uint16(util.RandomIntNumber(0, uint32(product.Quantity)))

	var purchaseQuantity uint16 = 2

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.PurchaseTransaction(context.Background(), PurchaseTransactionPagrams{
				Product:           product,
				PurchaseQuantity:  purchaseQuantity,
				Buyer:             buyer.Username,
				CardNumberOfBuyer: buyer.CardNumber,
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
		entryOfBuyer := result.EntryOfBuyer
		require.NotEmpty(t, entryOfBuyer)
		require.Equal(t, entryOfBuyer.BankAccount.Username, buyer.Username)
		require.Equal(t, entryOfBuyer.BankAccount.CardNumber, buyer.CardNumber)
		require.Equal(t, entryOfBuyer.BankAccount.Currency, buyer.Currency)

		entryOfSeller := result.EntryOfSeller
		require.NotEmpty(t, entryOfSeller)
		require.Equal(t, entryOfSeller.BankAccount.Username, seller.Username)
		require.Equal(t, entryOfSeller.BankAccount.CardNumber, seller.CardNumber)
		require.Equal(t, entryOfSeller.BankAccount.Currency, seller.Currency)

		fmt.Println("Buyer: ", entryOfBuyer.BankAccount.Balance, " Seller: ", entryOfSeller.BankAccount.Balance)

		// * check purchase history
		purchaseHistory := result.PurchaseHistory
		require.NotEmpty(t, purchaseHistory)
		require.Equal(t, purchaseHistory.IDProduct, product.IDProduct)
		require.Equal(t, purchaseHistory.Buyer, buyer.Username)
		require.Equal(t, purchaseHistory.CardNumberOfBuyer, entryOfBuyer.BankAccount.CardNumber)
		require.NotZero(t, purchaseHistory.CreatedAt)

		_, err = store.GetPurchaseHistory(context.Background(), purchaseHistory.IDPurchaseHistory)
		require.NoError(t, err)

		// * check balance
		diff_1 := buyer.Balance - entryOfBuyer.BankAccount.Balance
		diff_2 := entryOfSeller.BankAccount.Balance - seller.Balance
		require.Equal(t, diff_1, diff_2)
		require.True(t, diff_1 >= 0)
		require.True(t, diff_1%product.Price == 0)

		k := int(diff_1 / (product.Price * int32(purchaseQuantity)))
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
	require.Equal(t, buyer.Balance-(int32(n)*int32(product.Price*int32(purchaseQuantity))), updatedBankAccountBuyer.Balance)
	require.Equal(t, seller.Balance+(int32(n)*int32(product.Price*int32(purchaseQuantity))), updatedBankAccountSeller.Balance)
}
