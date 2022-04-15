package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	plainTextPassword := RandomString(6)

	hashedPassword1, err := HashPassword(plainTextPassword)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	err = CheckPassword(plainTextPassword, hashedPassword1)
	require.NoError(t, err)

	wrongPlainTextPassword := RandomString(6)
	err = CheckPassword(wrongPlainTextPassword, hashedPassword1)
	require.Equal(t, err.Error(), bcrypt.ErrMismatchedHashAndPassword.Error())

	// if the same plain text password is hashed twice then the hash value in both the cases should be different
	hashedPassword2, err := HashPassword(plainTextPassword)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
