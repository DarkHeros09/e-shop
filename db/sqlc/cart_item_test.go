package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DarkHeros09/e-shop/v2/util"
	"github.com/stretchr/testify/require"
)

func createRandomCartItem(t *testing.T) CartItem {
	shoppingSession := createRandomShoppingSession(t)
	product := createRandomProduct(t)
	arg := CreateCartItemParams{
		SessionID: shoppingSession.ID,
		ProductID: product.ID,
		Quantity:  int32(util.RandomInt(0, 10)),
	}

	cartItem, err := testQueires.CreateCartItem(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, cartItem)

	require.Equal(t, arg.SessionID, cartItem.SessionID)
	require.Equal(t, arg.ProductID, cartItem.ProductID)
	require.Equal(t, arg.Quantity, cartItem.Quantity)

	require.NotEmpty(t, cartItem.ID)
	require.NotEmpty(t, cartItem.CreatedAt)
	require.NotEmpty(t, cartItem.UpdatedAt)
	require.Equal(t, cartItem.CreatedAt, cartItem.UpdatedAt, time.Second)

	return cartItem

}

func TestCreateCartItem(t *testing.T) {
	createRandomCartItem(t)
}

func TestGetCartItem(t *testing.T) {
	cartItem1 := createRandomCartItem(t)
	cartItem2, err := testQueires.GetCartItem(context.Background(), cartItem1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, cartItem2)

	require.Equal(t, cartItem1.ID, cartItem2.ID)
	require.Equal(t, cartItem1.ProductID, cartItem2.ProductID)
	require.Equal(t, cartItem1.SessionID, cartItem2.SessionID)
	require.Equal(t, cartItem1.Quantity, cartItem2.Quantity)
	require.Equal(t, cartItem1.CreatedAt, cartItem2.CreatedAt)
	require.Equal(t, cartItem1.UpdatedAt, cartItem2.UpdatedAt)

	require.Equal(t, cartItem2.CreatedAt, cartItem2.UpdatedAt)

}

func TestUpdateCartItem(t *testing.T) {
	cartItem1 := createRandomCartItem(t)
	arg := UpdateCartItemParams{
		ID:       cartItem1.ID,
		Quantity: int32(util.RandomInt(0, 10)),
	}

	cartItem2, err := testQueires.UpdateCartItem(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, cartItem2)

	require.Equal(t, cartItem1.ID, cartItem2.ID)
	require.Equal(t, cartItem1.ProductID, cartItem2.ProductID)
	require.Equal(t, cartItem1.SessionID, cartItem2.SessionID)
	require.Equal(t, arg.Quantity, cartItem2.Quantity)
	require.Equal(t, cartItem1.CreatedAt, cartItem2.CreatedAt, time.Second)
	require.NotEqual(t, cartItem1.UpdatedAt, cartItem2.UpdatedAt, time.Second)

	require.NotEqual(t, cartItem2.CreatedAt, cartItem2.UpdatedAt, time.Second)
}

func TestDeleteCartItem(t *testing.T) {
	cartItem1 := createRandomCartItem(t)
	err := testQueires.DeleteCartItem(context.Background(), cartItem1.ID)

	require.NoError(t, err)

	cartItem2, err := testQueires.GetCartItem(context.Background(), cartItem1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, cartItem2)
}

func TestListCartItem(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomCartItem(t)
	}

	arg := ListCartItemParams{
		Limit:  5,
		Offset: 5,
	}

	cartItems, err := testQueires.ListCartItem(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, cartItems)

	for _, cartItem := range cartItems {
		require.NotEmpty(t, cartItem)

	}

}
