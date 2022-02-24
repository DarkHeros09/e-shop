package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DarkHeros09/e-shop/v2/util"
	"github.com/stretchr/testify/require"
)

func createRandomDiscount(t *testing.T) Discount {

	arg := CreateDiscountParams{
		Name:            util.RandomUser(),
		Description:     util.RandomUser(),
		DiscountPercent: util.RandomDecimal(1, 100),
	}
	discount, err := testQueires.CreateDiscount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, discount)

	require.Equal(t, arg.Name, discount.Name)
	require.Equal(t, arg.Description, discount.Description)
	require.Equal(t, arg.DiscountPercent, discount.DiscountPercent)

	require.NotZero(t, discount.ID)
	require.NotZero(t, discount.CreatedAt)
	require.True(t, discount.UpdatedAt.IsZero())

	require.False(t, discount.Active)

	return discount

}
func TestCreateDiscount(t *testing.T) {
	createRandomDiscount(t)
}

func TestGetDiscount(t *testing.T) {
	discount1 := createRandomDiscount(t)
	discount2, err := testQueires.GetDiscount(context.Background(), discount1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, discount2)

	require.Equal(t, discount1.ID, discount2.ID)
	require.Equal(t, discount1.Name, discount2.Name)
	require.Equal(t, discount1.Description, discount2.Description)
	require.Equal(t, discount1.DiscountPercent, discount2.DiscountPercent)
	require.Equal(t, discount1.Active, discount2.Active)
	require.Equal(t, discount1.CreatedAt, discount2.CreatedAt, time.Second)
	require.Equal(t, discount1.UpdatedAt, discount2.UpdatedAt, time.Second)

	require.False(t, discount2.Active)
}

func TestUpdateDiscount(t *testing.T) {
	discount1 := createRandomDiscount(t)

	arg := UpdateDiscountParams{
		ID:     discount1.ID,
		Active: true,
	}

	discount2, err := testQueires.UpdateDiscount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, discount2)

	require.Equal(t, discount1.ID, discount2.ID)
	require.Equal(t, discount1.Name, discount2.Name)
	require.Equal(t, discount1.Description, discount2.Description)
	require.Equal(t, discount1.DiscountPercent, discount2.DiscountPercent)
	require.Equal(t, arg.Active, discount2.Active)
	require.Equal(t, discount1.CreatedAt, discount2.CreatedAt, time.Second)
	require.NotEqual(t, discount1.UpdatedAt, discount2.UpdatedAt, time.Second)
}
func TestDeleteDiscount(t *testing.T) {
	discount1 := createRandomDiscount(t)

	err := testQueires.DeleteDiscount(context.Background(), discount1.ID)

	require.NoError(t, err)

	discount2, err := testQueires.GetDiscount(context.Background(), discount1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, discount2)
}

func TestListDiscounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomDiscount(t)
	}
	arg := ListDiscountsParams{
		Limit:  5,
		Offset: 5,
	}

	Discounts, err := testQueires.ListDiscounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, Discounts, 5)

	for _, discount := range Discounts {
		require.NotEmpty(t, discount)

	}
}
