package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomProductInventory(t *testing.T) ProductInventory {
	arg := int32(10)

	productInventory, err := testQueires.CreateProductInventory(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, productInventory)

	require.Equal(t, arg, productInventory.Quantity)
	require.True(t, productInventory.Active)

	require.NotEmpty(t, productInventory.ID)
	require.NotEmpty(t, productInventory.CreatedAt)
	require.NotEmpty(t, productInventory.UpdatedAt)
	require.Equal(t, productInventory.CreatedAt, productInventory.UpdatedAt)

	return productInventory

}
func TestCreateProductInventory(t *testing.T) {
	createRandomProductCategory(t)
}

func TestGetProductInventory(t *testing.T) {
	productInventory1 := createRandomProductInventory(t)
	productInventory2, err := testQueires.GetProductInventory(context.Background(), productInventory1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, productInventory2)

	require.Equal(t, productInventory1.ID, productInventory2.ID)
	require.Equal(t, productInventory1.Quantity, productInventory2.Quantity)
	require.Equal(t, productInventory1.Active, productInventory2.Active)
	require.Equal(t, productInventory1.CreatedAt, productInventory2.CreatedAt, time.Second)
	require.Equal(t, productInventory1.UpdatedAt, productInventory2.UpdatedAt, time.Second)

	require.True(t, productInventory2.Active)

}

func TestUpdateProductInventory(t *testing.T) {
	productInventory1 := createRandomProductInventory(t)
	arg := UpdateProductInventoryParams{
		ID:       productInventory1.ID,
		Active:   false,
		Quantity: 20,
	}

	productInventory2, err := testQueires.UpdateProductInventory(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, productInventory2)

	require.Equal(t, productInventory1.ID, productInventory2.ID)
	require.Equal(t, arg.Quantity, productInventory2.Quantity)
	require.Equal(t, arg.Active, productInventory2.Active)
	require.Equal(t, productInventory1.CreatedAt, productInventory2.CreatedAt, time.Second)
	require.NotEqual(t, productInventory1.UpdatedAt, productInventory2.UpdatedAt, time.Second)

	require.False(t, productInventory2.Active)
}

func TestDeleteProductInventory(t *testing.T) {
	productInventory1 := createRandomProductInventory(t)
	err := testQueires.DeleteProductInventory(context.Background(), productInventory1.ID)

	require.NoError(t, err)

	productInventory2, err := testQueires.GetProductInventory(context.Background(), productInventory1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, productInventory2)

}

func TestListProductInventories(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomProductInventory(t)
	}
	arg := ListProductInventoriesParams{
		Limit:  5,
		Offset: 5,
	}

	productInventories, err := testQueires.ListProductInventories(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, productInventories, 5)

	for _, productInventory := range productInventories {
		require.NotEmpty(t, productInventory)

	}
}
