package sqlc

import (
	"context"
	"simple_shop/db/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomBankAccount(t *testing.T) BankAccount {
	username := createRandomUser(t)

	arg := CreateBankAccountParams{
		Username:   username.Username,
		CardNumber: util.RandomStringNumber(12),
		Currency:   "USD",
		Balance:    int32(util.RandomIntNumber(1000, 10000)),
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
func TestCreateBankAccount(t *testing.T) {
	createRandomBankAccount(t)
}
