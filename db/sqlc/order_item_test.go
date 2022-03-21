package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DarkHeros09/e-shop/v2/util"
	"github.com/stretchr/testify/require"
)

func createRandomOrderItem(t *testing.T) OrderItem {
	orderDetail := createRandomOrderDetail(t)
	product := createRandomProduct(t)
	arg := CreateOrderItemParams{
		OrderID:   orderDetail.ID,
		ProductID: product.ID,
		Quantity:  int32(util.RandomMoney()),
	}

	orderItem, err := testQueires.CreateOrderItem(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, orderItem)

	require.Equal(t, arg.OrderID, orderItem.OrderID)
	require.Equal(t, arg.ProductID, orderItem.ProductID)
	require.Equal(t, arg.Quantity, orderItem.Quantity)

	require.NotEmpty(t, orderItem.ID)
	require.NotEmpty(t, orderItem.CreatedAt)
	require.True(t, orderItem.UpdatedAt.IsZero())

	return orderItem

}
func TestCreateOrderItem(t *testing.T) {
	createRandomOrderItem(t)
}

func TestGetOrderItem(t *testing.T) {
	orderItem1 := createRandomOrderItem(t)

	orderDetail1, err := testQueires.GetOrderDetail(context.Background(), orderItem1.OrderID)

	arg := GetOrderItemParams{
		ID:     orderItem1.ID,
		UserID: orderDetail1.UserID,
	}

	orderItem2, err := testQueires.GetOrderItem(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, orderItem2)

	require.Equal(t, orderItem1.ID, orderItem2.ID)
	require.Equal(t, orderItem1.OrderID, orderItem2.OrderID)
	require.Equal(t, orderItem1.ProductID, orderItem2.ProductID)
	require.Equal(t, orderItem1.Quantity, orderItem2.Quantity)
	require.Equal(t, orderItem1.CreatedAt, orderItem2.CreatedAt)
	require.Equal(t, orderItem1.UpdatedAt, orderItem2.UpdatedAt)

}

// func TestGetOrderItemByOrderID(t *testing.T) {
// 	orderItem1 := createRandomOrderItem(t)
// 	orderID := orderItem1.OrderID

// 	orderItem2, err := testQueires.GetOrderItemByOrderDetailID(context.Background(), orderID)

// 	require.NoError(t, err)
// 	require.NotEmpty(t, orderItem2)

// 	require.Equal(t, orderItem1.ID, orderItem2.ID)
// 	require.Equal(t, orderItem1.OrderID, orderItem2.OrderID)
// 	require.Equal(t, orderItem1.ProductID, orderItem2.ProductID)
// 	require.Equal(t, orderItem1.Quantity, orderItem2.Quantity)
// 	require.Equal(t, orderItem1.CreatedAt, orderItem2.CreatedAt)
// 	require.Equal(t, orderItem1.UpdatedAt, orderItem2.UpdatedAt)

// }

func TestUpdateOrderItem(t *testing.T) {
	orderItem1 := createRandomOrderItem(t)
	arg := UpdateOrderItemParams{
		ID:       orderItem1.ID,
		Quantity: int32(util.RandomMoney()),
	}
	orderItem2, err := testQueires.UpdateOrderItem(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, orderItem2)

	require.Equal(t, orderItem1.ID, orderItem2.ID)
	require.Equal(t, orderItem1.OrderID, orderItem2.OrderID)
	require.Equal(t, orderItem1.ProductID, orderItem2.ProductID)
	require.Equal(t, arg.Quantity, orderItem2.Quantity)
	require.Equal(t, orderItem1.CreatedAt, orderItem2.CreatedAt)
	require.NotEqual(t, orderItem1.UpdatedAt, orderItem2.UpdatedAt)

}

func TestDeleteOrderItem(t *testing.T) {
	orderItem1 := createRandomOrderItem(t)
	err := testQueires.DeleteOrderItem(context.Background(), orderItem1.ID)

	require.NoError(t, err)

	orderDetail1, err := testQueires.GetOrderDetail(context.Background(), orderItem1.OrderID)

	arg := GetOrderItemParams{
		ID:     orderItem1.ID,
		UserID: orderDetail1.UserID,
	}

	orderItem2, err := testQueires.GetOrderItem(context.Background(), arg)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, orderItem2)
}

func TestListOrderItems(t *testing.T) {
	var lastOrderItem OrderItem
	for i := 0; i < 10; i++ {
		lastOrderItem = createRandomOrderItem(t)

	}

	orderDetail1, err := testQueires.GetOrderDetail(context.Background(), lastOrderItem.OrderID)

	arg := ListOrderItemsParams{
		UserID: orderDetail1.UserID,
		Limit:  5,
		Offset: 0,
	}

	orderItems, err := testQueires.ListOrderItems(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, orderItems)

	for _, orderItem := range orderItems {
		require.NotEmpty(t, orderItem)
		require.Equal(t, lastOrderItem.ID, orderItem.ID)
		require.Equal(t, lastOrderItem.OrderID, orderItem.OrderID)
		require.Equal(t, lastOrderItem.ProductID, orderItem.ProductID)

	}

}
