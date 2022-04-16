package token

import (
	"testing"
	"time"

	"github.com/skamranahmed/banking-system/utils"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)

	userID := uint(utils.RandomInt(1, 1000))
	duration := time.Minute

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(userID, duration)
	require.NotNil(t, payload)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, userID, payload.UserID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiresAt, payload.ExpiresAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)

	userID := uint(utils.RandomInt(1, 1000))
	duration := time.Minute

	token, payload, err := maker.CreateToken(userID, -duration)
	require.NotNil(t, payload)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err = maker.VerifyToken(token)
	require.Equal(t, err.Error(), ErrExpiredToken.Error())
	require.Nil(t, payload)
}
