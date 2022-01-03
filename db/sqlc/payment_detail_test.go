package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomPaymentDetail(t *testing.T) PaymentDetail {
	arg := CreatePaymentDetailParams{
		OrderID:  0,
		Amount:   0,
		Provider: "cash",
		Status:   "pending",
	}

	paymentDetail, err := testQueires.CreatePaymentDetail(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, paymentDetail)

	require.Equal(t, arg.OrderID, paymentDetail.OrderID)
	require.Equal(t, arg.Amount, paymentDetail.Amount)
	require.Equal(t, arg.Provider, paymentDetail.Provider)
	require.Equal(t, arg.Status, paymentDetail.Status)

	require.NotEmpty(t, paymentDetail.ID)
	require.NotEmpty(t, paymentDetail.CreatedAt)
	require.NotEmpty(t, paymentDetail.UpdatedAt)

	require.Equal(t, paymentDetail.CreatedAt, paymentDetail.UpdatedAt)

	return paymentDetail

}
func TestCreatePaymentDetail(t *testing.T) {
	createRandomPaymentDetail(t)
}

func TestGetPaymentDetail(t *testing.T) {
	paymentDetail1 := createRandomPaymentDetail(t)
	paymentDetail2, err := testQueires.GetPaymentDetail(context.Background(), paymentDetail1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, paymentDetail2)

	require.Equal(t, paymentDetail1.ID, paymentDetail2.ID)
	require.Equal(t, paymentDetail1.OrderID, paymentDetail2.OrderID)
	require.Equal(t, paymentDetail1.Amount, paymentDetail2.Amount)
	require.Equal(t, paymentDetail1.Provider, paymentDetail2.Provider)
	require.Equal(t, paymentDetail1.Status, paymentDetail2.Status)
	require.Equal(t, paymentDetail1.CreatedAt, paymentDetail2.CreatedAt)
	require.Equal(t, paymentDetail1.UpdatedAt, paymentDetail2.UpdatedAt)

}

func TestUpdatePaymentDetail(t *testing.T) {
	paymentDetail1 := createRandomPaymentDetail(t)
	arg := UpdatePaymentDetailParams{
		ID:       paymentDetail1.ID,
		Amount:   0,
		Provider: "payme",
		Status:   "approved",
	}
	paymentDetail2, err := testQueires.UpdatePaymentDetail(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, paymentDetail2)

	require.Equal(t, paymentDetail1.ID, paymentDetail2.ID)
	require.Equal(t, paymentDetail1.OrderID, paymentDetail2.OrderID)
	require.Equal(t, paymentDetail1.Amount, paymentDetail2.Amount)
	require.Equal(t, arg.Provider, paymentDetail2.Provider)
	require.Equal(t, arg.Status, paymentDetail2.Status)
	require.Equal(t, paymentDetail1.CreatedAt, paymentDetail2.CreatedAt)
	require.NotEqual(t, paymentDetail1.UpdatedAt, paymentDetail2.UpdatedAt)

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
