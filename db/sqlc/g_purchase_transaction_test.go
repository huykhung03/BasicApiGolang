package sqlc

// func TestPurchaseTransaction(t *testing.T) {
// 	store := NewStore(testBD)

// 	usernameBuyer := createRandomBankAccount(t)

// 	// * run n concurrent puschase transactions
// 	n := 5
// 	product := createRandomProduct(t)

// 	errs := make(chan error)
// 	results := make(chan PurchaseTransactionResult)

// 	for i := 0; i < n; i++ {
// 		go func() {
// 			result, err := store.PurchaseTransaction(context.Background(), PurchaseTransactionPagrams{
// 				UsernameBuyer:  usernameBuyer.Username,
// 				UsernameSeller: product.Owner,
// 				Product:        product,
// 			})
// 			results <- result
// 			errs <- err
// 		}()
// 	}

// 	// * check results
// 	for i := 0; i < n; i++ {
// 		err := <-errs
// 		require.NoError(t, err)

// 		result := <-results
// 		require.NotEmpty(t, result)

// 		// * check purchase history
// 		purchaseHistory := result.PurchaseHistory
// 		require.NotEmpty(t, purchaseHistory)
// 		require.Equal(t, purchaseHistory.IDProduct, product.IDProduct)
// 		require.Equal(t, purchaseHistory.Buyer, usernameBuyer.Username)
// 		require.Equal(t, purchaseHistory.CardNumber, usernameBuyer.CardNumber)
// 		require.NotZero(t, purchaseHistory.CreatedAt)

// 		_, err = store.GetPurchaseHistory(context.Background(), purchaseHistory.IDProduct)
// 		require.NoError(t, err)

// 		// * check balance

// 	}

// }
