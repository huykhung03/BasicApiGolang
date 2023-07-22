package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hashedPassword_1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword_1)

	err = CheckPassword(password, hashedPassword_1)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPassword_1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPassword_2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword_2)
	require.NotEqual(t, hashedPassword_1, hashedPassword_2)
}
