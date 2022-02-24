package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DarkHeros09/e-shop/v2/util"
	"github.com/stretchr/testify/require"
)

func createRandomProduct(t *testing.T) Product {
	category := createRandomProductCategory(t)
	inventory := createRandomProductInventory(t)
	discount := createRandomDiscount(t)
	arg := CreateProductParams{
		Name:        util.RandomUser(),
		Description: util.RandomUser(),
		Sku:         util.RandomString(6),
		CategoryID:  category.ID,
		InventoryID: inventory.ID,
		Price:       util.RandomDecimal(1, 100),
		DiscountID:  discount.ID,
	}

	product, err := testQueires.CreateProduct(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, product)

	require.Equal(t, arg.Name, product.Name)
	require.Equal(t, arg.Description, product.Description)
	require.Equal(t, arg.CategoryID, product.CategoryID)
	require.Equal(t, arg.DiscountID, product.DiscountID)
	require.Equal(t, arg.InventoryID, product.InventoryID)
	require.Equal(t, arg.Price, product.Price)
	require.Equal(t, arg.Sku, product.Sku)

	require.NotEmpty(t, product.CreatedAt)
	require.True(t, product.UpdatedAt.IsZero())
	require.False(t, product.Active)

	return product
}
func TestCreateProduct(t *testing.T) {
	createRandomProduct(t)
}

func TestGetProduct(t *testing.T) {
	product1 := createRandomProduct(t)
	product2, err := testQueires.GetProduct(context.Background(), product1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, product2)

	require.Equal(t, product1.ID, product2.ID)
	require.Equal(t, product1.Name, product2.Name)
	require.Equal(t, product1.Description, product2.Description)
	require.Equal(t, product1.Price, product2.Price)
	require.Equal(t, product1.Sku, product2.Sku)
	require.Equal(t, product1.CategoryID, product2.CategoryID)
	require.Equal(t, product1.InventoryID, product2.InventoryID)
	require.Equal(t, product1.DiscountID, product2.DiscountID)
	require.Equal(t, product1.Active, product2.Active)
	require.Equal(t, product1.CreatedAt, product2.CreatedAt)
	require.Equal(t, product1.UpdatedAt, product2.UpdatedAt)

	require.False(t, product2.Active)

}

func TestUpdateProduct(t *testing.T) {
	product1 := createRandomProduct(t)
	arg := UpdateProductParams{
		ID:          product1.ID,
		Name:        product1.Name,
		Description: "new discription",
		CategoryID:  product1.CategoryID,
		Price:       fmt.Sprint(util.RandomMoney()),
		Active:      true,
	}

	product2, err := testQueires.UpdateProduct(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, product2)

	require.Equal(t, product1.ID, product2.ID)
	require.Equal(t, product1.Name, product2.Name)
	require.Equal(t, arg.Description, product2.Description)
	require.Equal(t, arg.Price, product2.Price)
	require.Equal(t, product1.Sku, product2.Sku)
	require.Equal(t, product1.CategoryID, product2.CategoryID)
	require.Equal(t, product1.InventoryID, product2.InventoryID)
	require.Equal(t, arg.Active, product2.Active)
	require.WithinDuration(t, product1.CreatedAt, product2.CreatedAt, time.Second)
	require.NotEqual(t, product1.UpdatedAt, product2.UpdatedAt)

	require.True(t, product2.Active)

}

func TestDeleteProduct(t *testing.T) {
	product1 := createRandomProduct(t)
	err := testQueires.DeleteProduct(context.Background(), product1.ID)

	require.NoError(t, err)

	product2, err := testQueires.GetProduct(context.Background(), product1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, product2)

}

func TestListProducts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomProduct(t)
	}
	arg := ListProductsParams{
		Limit:  5,
		Offset: 5,
	}

	products, err := testQueires.ListProducts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, products)

	for _, product := range products {
		require.NotEmpty(t, product)
	}

}
