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
		Expiry:      time.Now().Local(),
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
	userPayment1 := createRandomUserPayment(t)
	arg := GetUserPaymentParams{
		ID:     userPayment1.ID,
		UserID: userPayment1.UserID,
	}
	userPayment2, err := testQueires.GetUserPayment(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, userPayment2)

	require.Equal(t, userPayment1.ID, userPayment2.ID)
	require.Equal(t, userPayment1.UserID, userPayment2.UserID)
	require.Equal(t, userPayment1.PaymentType, userPayment2.PaymentType)
	require.Equal(t, userPayment1.Provider, userPayment2.Provider)
	require.Equal(t, userPayment1.AccountNo, userPayment2.AccountNo)
}

func TestUpdateUserPayment(t *testing.T) {
	userPayment1 := createRandomUserPayment(t)

	arg := UpdateUserPaymentParams{
		UserID:      userPayment1.UserID,
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

	arg := DeleteUserPaymentParams{
		ID:     userPayment1.ID,
		UserID: userPayment1.UserID,
	}

	err := testQueires.DeleteUserPayment(context.Background(), arg)

	require.NoError(t, err)

	arg1 := GetUserPaymentParams{
		ID:     userPayment1.ID,
		UserID: userPayment1.UserID,
	}
	userPayment2, err := testQueires.GetUserPayment(context.Background(), arg1)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, userPayment2)
}

func TestListUserPayments(t *testing.T) {
	var lastUserPayment UserPayment
	for i := 0; i < 10; i++ {
		lastUserPayment = createRandomUserPayment(t)
	}
	arg := ListUserPaymentsParams{
		UserID: lastUserPayment.UserID,
		Limit:  5,
		Offset: 0,
	}

	userPayments, err := testQueires.ListUserPayments(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, userPayments, 1)

	for _, userPayment := range userPayments {
		require.NotEmpty(t, userPayment)
		require.Equal(t, lastUserPayment.ID, userPayment.ID)
		require.Equal(t, lastUserPayment.UserID, userPayment.UserID)
		require.Equal(t, lastUserPayment.PaymentType, userPayment.PaymentType)
		require.Equal(t, lastUserPayment.Provider, userPayment.Provider)
		require.Equal(t, lastUserPayment.AccountNo, userPayment.AccountNo)
		require.Equal(t, lastUserPayment.Expiry, userPayment.Expiry)

	}
}
