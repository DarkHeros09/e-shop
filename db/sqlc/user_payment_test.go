package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DarkHeros09/e-shop/v2/util"
	"github.com/stretchr/testify/require"
)

func createRandomUserPayment(t *testing.T) UserPayment {
	user1 := createRandomUser(t)
	arg := CreateUserPaymentParams{
		UserID:      user1.ID,
		PaymentType: util.RandomUser(),
		Provider:    util.RandomUser(),
		AccountNo:   int32(util.RandomInt(0, 10)),
		Expiry:      time.Date(2022, 01, 01, 1, 1, 1, 1, time.Local),
	}

	userPayment, err := testQueires.CreateUserPayment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, userPayment)

	require.Equal(t, arg.UserID, userPayment.UserID)
	require.Equal(t, arg.PaymentType, userPayment.PaymentType)
	require.Equal(t, arg.Provider, userPayment.Provider)
	require.Equal(t, arg.AccountNo, userPayment.AccountNo)

	require.NotEmpty(t, userPayment.Expiry)

	return userPayment
}

func TestCreateUserPayment(t *testing.T) {
	createRandomUserPayment(t)
}

func TestGetUserPayment(t *testing.T) {
	arg := createRandomUserPayment(t)
	userPayment, err := testQueires.GetUserPayment(context.Background(), arg.ID)

	require.NoError(t, err)
	require.NotEmpty(t, userPayment)

	require.Equal(t, arg.ID, userPayment.ID)
	require.Equal(t, arg.UserID, userPayment.UserID)
	require.Equal(t, arg.PaymentType, userPayment.PaymentType)
	require.Equal(t, arg.Provider, userPayment.Provider)
	require.Equal(t, arg.AccountNo, userPayment.AccountNo)
}

func TestUpdateUserPayment(t *testing.T) {
	userPayment1 := createRandomUserPayment(t)

	arg := UpdateUserPaymentParams{
		ID:          userPayment1.ID,
		PaymentType: "edfa3le",
	}

	userPayment2, err := testQueires.UpdateUserPayment(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, userPayment2)

	require.Equal(t, userPayment1.ID, userPayment2.ID)
	require.Equal(t, userPayment1.UserID, userPayment2.UserID)
	require.Equal(t, arg.PaymentType, userPayment2.PaymentType)
	require.Equal(t, userPayment1.Provider, userPayment2.Provider)
	require.Equal(t, userPayment1.AccountNo, userPayment2.AccountNo)
}
func TestDeleteUserPayment(t *testing.T) {
	userPayment1 := createRandomUserPayment(t)

	err := testQueires.DeleteUserPayment(context.Background(), userPayment1.ID)

	require.NoError(t, err)

	userPayment2, err := testQueires.GetUserPayment(context.Background(), userPayment1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, userPayment2)
}

func TestListUserPayments(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUserPayment(t)
	}
	arg := ListUserPaymentsParams{
		Limit:  5,
		Offset: 5,
	}

	userPayments, err := testQueires.ListUserPayments(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, userPayments, 5)

	for _, userPayment := range userPayments {
		require.NotEmpty(t, userPayment)

	}
}
