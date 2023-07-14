package sqlc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	arg := CreateUserParams{
		Username:       "user1",
		FullName:       "user1",
		HashedPassword: "123",
		Email:          "user1@gmai.com",
	}

	username, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, username)

	require.Equal(t, arg.Username, username.Username)
	require.Equal(t, arg.FullName, username.FullName)
	require.Equal(t, arg.HashedPassword, username.HashedPassword)
	require.Equal(t, arg.Email, username.Email)
}
