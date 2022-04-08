package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DarkHeros09/e-shop/v2/util"
	"github.com/stretchr/testify/require"
)

func createRandomOrderDetailAndPaymentDetail(t *testing.T) (OrderDetail, PaymentDetail) {
	user := createRandomUser(t)
	arg := CreateOrderDetailAndPaymentDetailParams{
		UserID: user.ID,
		Total:  util.RandomDecimal(1, 100),
	}

	orderDetail, err := testQueires.CreateOrderDetailAndPaymentDetail(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, orderDetail)
	require.NotEmpty(t, orderDetail.PaymentID.Int64)
	require.NotEmpty(t, orderDetail.UserID)

	arg1 := UpdatePaymentDetailParams{
		ID:       orderDetail.PaymentID.Int64,
		UserID:   orderDetail.UserID,
		OrderID:  orderDetail.ID,
		Amount:   int32(util.RandomMoney()),
		Provider: util.RandomUser(),
		Status:   util.RandomUser(),
	}

	updatedPaymentDetail, err := testQueires.UpdatePaymentDetail(context.Background(), arg1)

	require.NoError(t, err)
	require.NotEmpty(t, updatedPaymentDetail)

	arg2 := GetPaymentDetailParams{
		ID:     orderDetail.PaymentID.Int64,
		UserID: orderDetail.UserID,
	}

	paymentDetail, err := testQueires.GetPaymentDetail(context.Background(), arg2)

	require.NoError(t, err)
	require.NotEmpty(t, paymentDetail)

	require.Equal(t, arg.UserID, orderDetail.UserID)
	require.Equal(t, arg.Total, orderDetail.Total)
	require.Equal(t, paymentDetail.ID, orderDetail.PaymentID.Int64)

	require.NotEmpty(t, orderDetail.ID)
	require.NotEmpty(t, orderDetail.CreatedAt)
	require.True(t, orderDetail.UpdatedAt.IsZero())
	require.NotZero(t, orderDetail.PaymentID.Int64)

	return orderDetail, paymentDetail

}

func TestCreateOrderDetailAndPaymentDetail(t *testing.T) {
	createRandomOrderDetailAndPaymentDetail(t)
}

func TestGetOrderDetail(t *testing.T) {
	orderDetail1, _ := createRandomOrderDetailAndPaymentDetail(t)
	orderDetail2, err := testQueires.GetOrderDetail(context.Background(), orderDetail1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, orderDetail2)

	require.Equal(t, orderDetail1.ID, orderDetail2.ID)
	require.Equal(t, orderDetail1.UserID, orderDetail2.UserID)
	require.Equal(t, orderDetail1.PaymentID.Int64, orderDetail2.PaymentID.Int64)
	require.Equal(t, orderDetail1.Total, orderDetail2.Total)
	require.Equal(t, orderDetail1.CreatedAt, orderDetail2.CreatedAt)
	require.Equal(t, orderDetail1.UpdatedAt, orderDetail2.UpdatedAt)

}

func TestUpdateOrderDetail(t *testing.T) {
	orderDetail1, _ := createRandomOrderDetailAndPaymentDetail(t)
	arg := UpdateOrderDetailParams{
		ID:    orderDetail1.ID,
		Total: fmt.Sprint(util.RandomMoney()),
	}

	orderDetail2, err := testQueires.UpdateOrderDetail(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, orderDetail2)

	require.Equal(t, orderDetail1.ID, orderDetail2.ID)
	require.Equal(t, orderDetail1.UserID, orderDetail2.UserID)
	require.Equal(t, orderDetail1.PaymentID.Int64, orderDetail2.PaymentID.Int64)
	require.Equal(t, arg.Total, orderDetail2.Total)
	require.Equal(t, orderDetail1.CreatedAt, orderDetail2.CreatedAt)
	require.NotEqual(t, orderDetail1.UpdatedAt, orderDetail2.UpdatedAt)

}

func TestDeleteOrderDetail(t *testing.T) {
	orderDetail1, _ := createRandomOrderDetailAndPaymentDetail(t)
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
		lastOrderDetail, _ = createRandomOrderDetailAndPaymentDetail(t)

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
		require.Equal(t, lastOrderDetail.PaymentID.Int64, orderDetail.PaymentID.Int64)
		require.Equal(t, lastOrderDetail.Total, orderDetail.Total)
	}
}
