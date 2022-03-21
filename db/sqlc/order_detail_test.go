package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DarkHeros09/e-shop/v2/util"
	"github.com/stretchr/testify/require"
)

func createRandomOrderDetail(t *testing.T) OrderDetail {
	user := createRandomUser(t)
	paymentDetail := createRandomPaymentDetail(t)
	arg := CreateOrderDetailParams{
		UserID:    user.ID,
		Total:     util.RandomDecimal(1, 100),
		PaymentID: paymentDetail.ID,
	}

	orderDetail, err := testQueires.CreateOrderDetail(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, orderDetail)

	require.Equal(t, arg.UserID, orderDetail.UserID)
	require.Equal(t, arg.Total, orderDetail.Total)
	require.Equal(t, arg.PaymentID, orderDetail.PaymentID)

	require.NotEmpty(t, orderDetail.ID)
	require.NotEmpty(t, orderDetail.CreatedAt)
	require.True(t, orderDetail.UpdatedAt.IsZero())

	return orderDetail

}
func TestCreateOrderDetail(t *testing.T) {
	createRandomOrderDetail(t)
}

func TestGetOrderDetail(t *testing.T) {
	orderDetail1 := createRandomOrderDetail(t)
	orderDetail2, err := testQueires.GetOrderDetail(context.Background(), orderDetail1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, orderDetail2)

	require.Equal(t, orderDetail1.ID, orderDetail2.ID)
	require.Equal(t, orderDetail1.UserID, orderDetail2.UserID)
	require.Equal(t, orderDetail1.PaymentID, orderDetail2.PaymentID)
	require.Equal(t, orderDetail1.Total, orderDetail2.Total)
	require.Equal(t, orderDetail1.CreatedAt, orderDetail2.CreatedAt)
	require.Equal(t, orderDetail1.UpdatedAt, orderDetail2.UpdatedAt)

}

func TestUpdateOrderDetail(t *testing.T) {
	orderDetail1 := createRandomOrderDetail(t)
	arg := UpdateOrderDetailParams{
		ID:    orderDetail1.ID,
		Total: fmt.Sprint(util.RandomMoney()),
	}

	orderDetail2, err := testQueires.UpdateOrderDetail(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, orderDetail2)

	require.Equal(t, orderDetail1.ID, orderDetail2.ID)
	require.Equal(t, orderDetail1.UserID, orderDetail2.UserID)
	require.Equal(t, orderDetail1.PaymentID, orderDetail2.PaymentID)
	require.Equal(t, arg.Total, orderDetail2.Total)
	require.Equal(t, orderDetail1.CreatedAt, orderDetail2.CreatedAt)
	require.NotEqual(t, orderDetail1.UpdatedAt, orderDetail2.UpdatedAt)

}

func TestDeleteOrderDetail(t *testing.T) {
	orderDetail1 := createRandomOrderDetail(t)
	err := testQueires.DeleteOrderDetail(context.Background(), orderDetail1.ID)

	require.NoError(t, err)

	orderDetail2, err := testQueires.GetOrderDetail(context.Background(), orderDetail1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, orderDetail2)

}

func TestListOrderDetails(t *testing.T) {
	var lastOrderDetail OrderDetail
	for i := 0; i < 10; i++ {
		lastOrderDetail = createRandomOrderDetail(t)

	}
	arg := ListOrderDetailsParams{
		UserID: lastOrderDetail.UserID,
		Limit:  5,
		Offset: 0,
	}
	orderDetails, err := testQueires.ListOrderDetails(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, orderDetails)

	for _, orderDetail := range orderDetails {
		require.NotEmpty(t, orderDetail)
		require.Equal(t, lastOrderDetail.ID, orderDetail.ID)
		require.Equal(t, lastOrderDetail.UserID, orderDetail.UserID)
		require.Equal(t, lastOrderDetail.PaymentID, orderDetail.PaymentID)
		require.Equal(t, lastOrderDetail.Total, orderDetail.Total)
	}
}
