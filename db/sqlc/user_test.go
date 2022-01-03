package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DarkHeros09/e-shop/v2/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {

	arg := CreateUserParams{
		Username:  util.RandomUser(),
		Email:     util.RandomEmail(),
		Password:  util.RandomString(5),
		Telephone: int32(util.RandomInt(0, 7)),
	}
	user, err := testQueires.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Password, user.Password)
	require.Equal(t, arg.Telephone, user.Telephone)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)

	return user

}
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueires.GetUser(context.Background(), user1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.Telephone, user2.Telephone)
	require.Equal(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.Equal(t, user1.UpdatedAt, user2.UpdatedAt, time.Second)
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t)

	arg := UpdateUserParams{
		ID:        user1.ID,
		Telephone: int32(util.RandomInt(0, 10)),
	}

	user2, err := testQueires.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, arg.Telephone, user2.Telephone)
	require.Equal(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.NotEqual(t, user1.UpdatedAt, user2.UpdatedAt, time.Second)
}
func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)

	err := testQueires.DeleteUser(context.Background(), user1.ID)

	require.NoError(t, err)

	user2, err := testQueires.GetUser(context.Background(), user1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}
	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testQueires.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)

	}
}
