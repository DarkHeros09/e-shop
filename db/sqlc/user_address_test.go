package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DarkHeros09/e-shop/v2/util"
	"github.com/stretchr/testify/require"
)

func createRandomUserAddress(t *testing.T) UserAddress {
	user1 := createRandomUser(t)

	arg := CreateUserAddressParams{
		UserID:      user1.ID,
		AddressLine: util.RandomString(5),
		City:        util.RandomString(5),
		Telephone:   user1.Telephone,
	}

	userAddress, err := testQueires.CreateUserAddress(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, userAddress)

	require.Equal(t, arg.UserID, userAddress.UserID)
	require.Equal(t, arg.AddressLine, userAddress.AddressLine)
	require.Equal(t, arg.City, userAddress.City)
	require.Equal(t, arg.Telephone, userAddress.Telephone)

	return userAddress

}

func TestCreateUserAddress(t *testing.T) {
	createRandomUserAddress(t)
}

func TestGetUserAddress(t *testing.T) {
	userAdress1 := createRandomUserAddress(t)
	userAdress2, err := testQueires.GetUserAddressByUserID(context.Background(), userAdress1.UserID)

	require.NoError(t, err)
	require.NotEmpty(t, userAdress2)

	require.Equal(t, userAdress1.ID, userAdress2.ID)
	require.Equal(t, userAdress1.UserID, userAdress2.UserID)
	require.Equal(t, userAdress1.AddressLine, userAdress2.AddressLine)
	require.Equal(t, userAdress1.City, userAdress2.City)
	require.Equal(t, userAdress1.Telephone, userAdress2.Telephone)

}

func TestGetUserAddressByUserID(t *testing.T) {
	userAdress1 := createRandomUserAddress(t)
	userAdress2, err := testQueires.GetUserAddress(context.Background(), userAdress1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, userAdress2)

	require.Equal(t, userAdress1.ID, userAdress2.ID)
	require.Equal(t, userAdress1.UserID, userAdress2.UserID)
	require.Equal(t, userAdress1.AddressLine, userAdress2.AddressLine)
	require.Equal(t, userAdress1.City, userAdress2.City)
	require.Equal(t, userAdress1.Telephone, userAdress2.Telephone)

}

func TestUpdateUserAddress(t *testing.T) {
	userAddress1 := createRandomUserAddress(t)
	arg := UpdateUserAddressParams{
		ID:          userAddress1.ID,
		AddressLine: "NewAddress",
		City:        "Benghazi",
		Telephone:   123456,
	}

	userAddress2, err := testQueires.UpdateUserAddress(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, userAddress2)

	require.Equal(t, userAddress1.ID, userAddress2.ID)
	require.Equal(t, userAddress1.UserID, userAddress2.UserID)
	require.Equal(t, arg.AddressLine, userAddress2.AddressLine)
	require.Equal(t, arg.City, userAddress2.City)
	require.Equal(t, arg.Telephone, userAddress2.Telephone)

}

func TestDeleteUserAddress(t *testing.T) {
	useraddress1 := createRandomUserAddress(t)
	err := testQueires.DeleteUserAddress(context.Background(), useraddress1.ID)

	require.NoError(t, err)

	useraddress2, err := testQueires.GetUserAddress(context.Background(), useraddress1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, useraddress2)

}

func TestListUserAddresses(t *testing.T) {
	var lastUserAddress UserAddress
	for i := 0; i < 10; i++ {
		lastUserAddress = createRandomUserAddress(t)
	}
	arg := ListUserAddressesParams{
		UserID: lastUserAddress.UserID,
		Limit:  5,
		Offset: 0,
	}

	userAddresses, err := testQueires.ListUserAddresses(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, userAddresses, 1)

	for _, userAddress := range userAddresses {
		require.NotEmpty(t, userAddress)
		require.Equal(t, lastUserAddress.ID, userAddress.ID)
		require.Equal(t, lastUserAddress.UserID, userAddress.UserID)
		require.Equal(t, lastUserAddress.AddressLine, userAddress.AddressLine)
		require.Equal(t, lastUserAddress.City, userAddress.City)
		require.Equal(t, lastUserAddress.Telephone, userAddress.Telephone)

	}
}
