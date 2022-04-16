package db

import (
	"context"
	"testing"
	"time"

	"github.com/DarkHeros09/e-shop/v2/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomUserSession(t *testing.T) UserSession {
	user1 := createRandomUser(t)
	arg := CreateUserSessionParams{
		ID:           uuid.New(),
		UserID:       user1.ID,
		RefreshToken: util.RandomString(100),
		UserAgent:    util.RandomString(6),
		ClientIp:     util.RandomString(10),
		IsBlocked:    false,
		ExpiresAt:    time.Now().Local().UTC(),
	}

	userSession, err := testQueires.CreateUserSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, userSession)

	require.Equal(t, arg.ID, userSession.ID)
	require.Equal(t, arg.UserID, userSession.UserID)
	require.Equal(t, arg.RefreshToken, userSession.RefreshToken)
	require.Equal(t, arg.UserAgent, userSession.UserAgent)
	require.Equal(t, arg.ClientIp, userSession.ClientIp)
	require.Equal(t, arg.IsBlocked, userSession.IsBlocked)
	require.WithinDuration(t, arg.ExpiresAt, userSession.ExpiresAt, time.Second)

	require.NotEmpty(t, userSession.CreatedAt)

	return userSession
}

func TestCreateUserSession(t *testing.T) {
	createRandomUserSession(t)
}

func TestGetUserSession(t *testing.T) {
	userSession1 := createRandomUserSession(t)

	userSession2, err := testQueires.GetUserSession(context.Background(), userSession1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, userSession2)

	require.Equal(t, userSession1.ID, userSession2.ID)
	require.Equal(t, userSession1.UserID, userSession2.UserID)
	require.Equal(t, userSession1.RefreshToken, userSession2.RefreshToken)
	require.Equal(t, userSession1.UserAgent, userSession2.UserAgent)
	require.Equal(t, userSession1.ClientIp, userSession2.ClientIp)
	require.Equal(t, userSession1.IsBlocked, userSession2.IsBlocked)
	require.Equal(t, userSession1.CreatedAt, userSession2.CreatedAt)
	require.Equal(t, userSession1.ExpiresAt, userSession2.ExpiresAt)
}
