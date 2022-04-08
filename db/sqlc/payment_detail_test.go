package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomPaymentDetail(t *testing.T) PaymentDetail {
	_, paymentDetail := createRandomOrderDetailAndPaymentDetail(t)

	require.NotEmpty(t, paymentDetail)

	require.NotEmpty(t, paymentDetail.Amount)
	require.NotEmpty(t, paymentDetail.Provider)
	require.NotEmpty(t, paymentDetail.Status)

	require.NotEmpty(t, paymentDetail.ID)
	require.NotEmpty(t, paymentDetail.CreatedAt)
	require.False(t, paymentDetail.UpdatedAt.IsZero())

	return paymentDetail

}
func TestCreatePaymentDetail(t *testing.T) {
	createRandomPaymentDetail(t)
}

func TestGetPaymentDetail(t *testing.T) {
	orderDetail1, paymentDetail1 := createRandomOrderDetailAndPaymentDetail(t)

	arg := GetPaymentDetailParams{
		ID:     orderDetail1.PaymentID.Int64,
		UserID: orderDetail1.UserID,
	}
	paymentDetail2, err := testQueires.GetPaymentDetail(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, paymentDetail1)

	require.Equal(t, paymentDetail1.ID, paymentDetail2.ID)
	require.Equal(t, paymentDetail1.OrderID, paymentDetail2.OrderID)
	require.Equal(t, paymentDetail1.Amount, paymentDetail2.Amount)
	require.Equal(t, paymentDetail1.Provider, paymentDetail2.Provider)
	require.Equal(t, paymentDetail1.Status, paymentDetail2.Status)
	require.Equal(t, paymentDetail1.CreatedAt, paymentDetail2.CreatedAt)
	require.Equal(t, paymentDetail1.UpdatedAt, paymentDetail2.UpdatedAt)

}

func TestUpdatePaymentDetail(t *testing.T) {
	orderDetail1, paymentDetail1 := createRandomOrderDetailAndPaymentDetail(t)
	arg := UpdatePaymentDetailParams{
		ID:       paymentDetail1.ID,
		UserID:   orderDetail1.UserID,
		OrderID:  paymentDetail1.OrderID,
		Amount:   10,
		Provider: "payme",
		Status:   "approved",
	}
	paymentDetail2, err := testQueires.UpdatePaymentDetail(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, paymentDetail2)

	require.Equal(t, paymentDetail1.ID, paymentDetail2.ID)
	require.Equal(t, arg.OrderID, orderDetail1.ID)
	require.NotEqual(t, paymentDetail1.Amount, paymentDetail2.Amount)
	require.Equal(t, arg.Provider, paymentDetail2.Provider)
	require.Equal(t, arg.Status, paymentDetail2.Status)
	require.Equal(t, paymentDetail1.CreatedAt, paymentDetail2.CreatedAt)
	require.NotEqual(t, paymentDetail1.UpdatedAt, paymentDetail2.UpdatedAt)

}

func TestUpdatePaymentDetailOrderID(t *testing.T) {
	orderDetail, _ := createRandomOrderDetailAndPaymentDetail(t)
	arg := UpdatePaymentDetailParams{
		ID:       orderDetail.PaymentID.Int64,
		UserID:   orderDetail.UserID,
		OrderID:  orderDetail.ID,
		Amount:   10,
		Provider: "payme",
		Status:   "approved",
	}
	paymentDetail, err := testQueires.UpdatePaymentDetail(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, paymentDetail)

	require.Equal(t, orderDetail.PaymentID.Int64, paymentDetail.ID)
	require.Equal(t, orderDetail.ID, paymentDetail.OrderID)
	require.Equal(t, paymentDetail.CreatedAt, paymentDetail.CreatedAt)
	require.NotEqual(t, paymentDetail.CreatedAt, paymentDetail.UpdatedAt)

}

func TestDeletePaymentDetail(t *testing.T) {
	orderDetail, _ := createRandomOrderDetailAndPaymentDetail(t)

	err := testQueires.DeletePaymentDetail(context.Background(), orderDetail.PaymentID.Int64)
	require.NoError(t, err)
	arg := GetPaymentDetailParams{
		ID:     orderDetail.PaymentID.Int64,
		UserID: orderDetail.UserID,
	}
	paymentDetail, err := testQueires.GetPaymentDetail(context.Background(), arg)

	require.Empty(t, paymentDetail)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())

}

func TestListPaymentDetails(t *testing.T) {
	var lastOrderDetail OrderDetail
	for i := 0; i < 10; i++ {
		lastOrderDetail, _ = createRandomOrderDetailAndPaymentDetail(t)

	}
	arg := ListPaymentDetailsParams{
		UserID: lastOrderDetail.UserID,
		Limit:  5,
		Offset: 0,
	}

	paymentDetails, err := testQueires.ListPaymentDetails(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, paymentDetails)

	for _, paymentDetail := range paymentDetails {
		require.NotEmpty(t, paymentDetail)
	}

}
