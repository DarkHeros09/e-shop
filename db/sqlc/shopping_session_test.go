package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DarkHeros09/e-shop/v2/util"
	"github.com/stretchr/testify/require"
)

func createRandomShoppingSession(t *testing.T) ShoppingSession {
	user := createRandomUser(t)
	arg := CreateShoppingSessionParams{
		UserID: user.ID,
		Total:  fmt.Sprint(util.RandomMoney()),
	}

	shoppingSession, err := testQueires.CreateShoppingSession(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, shoppingSession)

	require.Equal(t, arg.UserID, shoppingSession.UserID)
	require.Equal(t, arg.Total, shoppingSession.Total)
	require.NotEmpty(t, shoppingSession.ID)
	require.NotEmpty(t, shoppingSession.CreatedAt)
	require.True(t, shoppingSession.UpdatedAt.IsZero())

	return shoppingSession

}
func TestCreateShoppingSession(t *testing.T) {
	createRandomShoppingSession(t)
}

func TestGetShoppingSession(t *testing.T) {
	shoppingSession1 := createRandomShoppingSession(t)
	shoppingSession2, err := testQueires.GetShoppingSession(context.Background(), shoppingSession1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, shoppingSession2)

	require.Equal(t, shoppingSession1.ID, shoppingSession2.ID)
	require.Equal(t, shoppingSession1.UserID, shoppingSession2.UserID)
	require.Equal(t, shoppingSession1.Total, shoppingSession2.Total)
	require.Equal(t, shoppingSession1.CreatedAt, shoppingSession2.CreatedAt)
	require.Equal(t, shoppingSession1.UpdatedAt, shoppingSession2.UpdatedAt)

}

func TestUpdateShoppingSession(t *testing.T) {
	shoppingSession1 := createRandomShoppingSession(t)
	arg := UpdateShoppingSessionParams{
		ID:    shoppingSession1.ID,
		Total: fmt.Sprint(util.RandomMoney()),
	}
	shoppingSession2, err := testQueires.UpdateShoppingSession(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, shoppingSession2)

	require.Equal(t, shoppingSession1.ID, shoppingSession2.ID)
	require.Equal(t, shoppingSession1.UserID, shoppingSession2.UserID)
	require.Equal(t, arg.Total, shoppingSession2.Total)
	require.Equal(t, shoppingSession1.CreatedAt, shoppingSession2.CreatedAt)
	require.NotEqual(t, shoppingSession1.UpdatedAt, shoppingSession2.UpdatedAt)
	require.NotEqual(t, shoppingSession1.CreatedAt, shoppingSession2.UpdatedAt, time.Second)
}

func TestDeleteShoppingSession(t *testing.T) {
	shoppingSession1 := createRandomShoppingSession(t)
	err := testQueires.DeleteShoppingSession(context.Background(), shoppingSession1.ID)

	require.NoError(t, err)

	shoppingSession2, err := testQueires.GetShoppingSession(context.Background(), shoppingSession1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, shoppingSession2)

}

func TestListShoppingSessions(t *testing.T) {
	var lastShoppingSession ShoppingSession
	for i := 0; i < 10; i++ {
		lastShoppingSession = createRandomShoppingSession(t)
	}

	arg := ListShoppingSessionsParams{
		UserID: lastShoppingSession.UserID,
		Limit:  5,
		Offset: 0,
	}

	shoppingSessions, err := testQueires.ListShoppingSessions(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, shoppingSessions)

	for _, shoppingSession := range shoppingSessions {
		require.NotEmpty(t, shoppingSession)
		require.Equal(t, lastShoppingSession.ID, shoppingSession.ID)
		require.Equal(t, lastShoppingSession.UserID, shoppingSession.UserID)
		require.Equal(t, lastShoppingSession.Total, shoppingSession.Total)
	}
}
