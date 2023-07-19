package sqlc

import (
	"context"
	"simple_shop/db/util"
	"simple_shop/db/util/random"
	"testing"

	"github.com/stretchr/testify/require"
)

// Three below creating functions create admin
func createRandomAdmin(t *testing.T) User {
	hashedPassword, err := util.HashPassword(random.RandomHashedPassword())
	require.NoError(t, err)

	arg := CreateAdminParams{
		Username:       random.RandomUsername(),
		FullName:       random.RandomFullName(),
		HashedPassword: hashedPassword,
		Email:          random.RandomEmail(),
		Level:          true,
	}

	username, err := testQueries.CreateAdmin(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, username)
	require.Equal(t, arg.Username, username.Username)
	require.Equal(t, arg.FullName, username.FullName)
	require.Equal(t, arg.HashedPassword, username.HashedPassword)
	require.Equal(t, arg.Email, username.Email)
	require.Equal(t, arg.Level, username.Level)

	return username
}

func createRandomBankAccountAdmin(t *testing.T) BankAccount {
	admin := createRandomAdmin(t)

	arg := CreateBankAccountParams{
		Username:   admin.Username,
		CardNumber: util.RandomStringNumber(9),
		Currency:   "USD",
		Balance:    1000,
	}

	bankAccount, err := testQueries.CreateBankAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, bankAccount)

	require.Equal(t, arg.Username, bankAccount.Username)
	require.Equal(t, arg.CardNumber, bankAccount.CardNumber)
	require.Equal(t, arg.Currency, bankAccount.Currency)
	require.Equal(t, arg.Balance, bankAccount.Balance)

	return bankAccount
}

func createRandomProduct(t *testing.T) Product {
	username := createRandomBankAccountAdmin(t)

	arg := CreateProductParams{
		ProductName:   util.RandomString(8),
		KindOfProduct: util.RandomString(8),
		Owner:         username.Username,
		Currency:      "USD",
		Price:         10,
		Quantity:      int32(util.RandomIntNumber(20, 50)),
	}

	product, err := testQueries.CreateProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, product)

	require.Equal(t, arg.ProductName, product.ProductName)
	require.Equal(t, arg.KindOfProduct, product.KindOfProduct)
	require.Equal(t, arg.Owner, product.Owner)
	require.Equal(t, arg.Currency, product.Currency)
	require.Equal(t, arg.Price, product.Price)
	require.Equal(t, arg.Quantity, product.Quantity)

	return product
}

func TestCreateProduct(t *testing.T) {
	createRandomProduct(t)
}
