package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomPaymentDetail(t *testing.T) PaymentDetail {
	arg := CreatePaymentDetailParams{
		Amount:   0,
		Provider: "cash",
		Status:   "pending",
	}

	paymentDetail, err := testQueires.CreatePaymentDetail(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, paymentDetail)

	require.Equal(t, arg.Amount, paymentDetail.Amount)
	require.Equal(t, arg.Provider, paymentDetail.Provider)
	require.Equal(t, arg.Status, paymentDetail.Status)

	require.NotEmpty(t, paymentDetail.ID)
	require.NotEmpty(t, paymentDetail.CreatedAt)
	require.True(t, paymentDetail.UpdatedAt.IsZero())

	// require.True(t, paymentDetail.ID)

	return paymentDetail

}
func TestCreatePaymentDetail(t *testing.T) {
	createRandomPaymentDetail(t)
}

func TestGetPaymentDetail(t *testing.T) {
	orderDetail1 := createRandomOrderDetail(t)

	paymentDetail1, err := testQueires.GetPaymentDetail(context.Background(), orderDetail1.PaymentID)

	require.NoError(t, err)
	require.NotEmpty(t, paymentDetail1)

	require.Equal(t, orderDetail1.PaymentID, paymentDetail1.ID)
	// require.Equal(t, orderDetail1.ID, paymentDetail1.OrderID)
	// require.Equal(t, orderDetail1.Amount, paymentDetail1.Amount)
	// require.Equal(t, orderDetail1.Provider, paymentDetail1.Provider)
	// require.Equal(t, orderDetail1.Status, paymentDetail1.Status)
	// require.Equal(t, orderDetail1.CreatedAt, paymentDetail1.CreatedAt)
	// require.Equal(t, orderDetail1.UpdatedAt, paymentDetail1.UpdatedAt)

}

func TestUpdatePaymentDetail(t *testing.T) {
	paymentDetail1 := createRandomPaymentDetail(t)
	orderItem := createRandomOrderItem(t)
	arg := UpdatePaymentDetailParams{
		ID:       paymentDetail1.ID,
		OrderID:  orderItem.OrderID,
		Amount:   0,
		Provider: "payme",
		Status:   "approved",
	}
	paymentDetail2, err := testQueires.UpdatePaymentDetail(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, paymentDetail2)

	require.Equal(t, paymentDetail1.ID, paymentDetail2.ID)
	require.Equal(t, arg.OrderID, orderItem.OrderID)
	require.Equal(t, paymentDetail1.Amount, paymentDetail2.Amount)
	require.Equal(t, arg.Provider, paymentDetail2.Provider)
	require.Equal(t, arg.Status, paymentDetail2.Status)
	require.Equal(t, paymentDetail1.CreatedAt, paymentDetail2.CreatedAt)
	require.NotEqual(t, paymentDetail1.UpdatedAt, paymentDetail2.UpdatedAt)

}

func TestUpdatePaymentDetailOrderID(t *testing.T) {
	orderDetail := createRandomOrderDetailAndPaymentDetail(t)
	arg := UpdatePaymentDetailParams{
		ID:       orderDetail.PaymentID,
		OrderID:  orderDetail.ID,
		Amount:   10,
		Provider: "payme",
		Status:   "approved",
	}
	paymentDetail, err := testQueires.UpdatePaymentDetail(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, paymentDetail)

	require.Equal(t, orderDetail.PaymentID, paymentDetail.ID)
	require.Equal(t, orderDetail.ID, paymentDetail.OrderID)
	require.Equal(t, paymentDetail.CreatedAt, paymentDetail.CreatedAt)
	require.NotEqual(t, paymentDetail.CreatedAt, paymentDetail.UpdatedAt)

}

func TestDeletePaymentDetail(t *testing.T) {
	paymentDetail1 := createRandomPaymentDetail(t)
	err := testQueires.DeletePaymentDetail(context.Background(), paymentDetail1.ID)
	require.NoError(t, err)

	paymentDetail2, err := testQueires.GetPaymentDetail(context.Background(), paymentDetail1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, paymentDetail2)

}

func TestListPaymentDetails(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomPaymentDetail(t)

	}
	arg := ListPaymentDetailsParams{
		Limit:  5,
		Offset: 5,
	}

	paymentDetails, err := testQueires.ListPaymentDetails(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, paymentDetails)

	for _, paymentDetail := range paymentDetails {
		require.NotEmpty(t, paymentDetail)
	}

}
