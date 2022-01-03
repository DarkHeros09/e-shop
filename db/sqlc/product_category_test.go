package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DarkHeros09/e-shop/v2/util"
	"github.com/stretchr/testify/require"
)

func createRandomProductCategory(t *testing.T) ProductCategory {
	arg := CreateProductCategoryParams{
		Name:        util.RandomString(5),
		Description: util.RandomString(5),
	}

	productCategory, err := testQueires.CreateProductCategory(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, productCategory)

	require.Equal(t, arg.Name, productCategory.Name)
	require.Equal(t, arg.Description, productCategory.Description)
	require.False(t, productCategory.Active)

	require.NotEmpty(t, productCategory.ID)
	require.NotEmpty(t, productCategory.CreatedAt)
	require.NotEmpty(t, productCategory.UpdatedAt)
	require.Equal(t, productCategory.CreatedAt, productCategory.UpdatedAt)

	return productCategory

}
func TestCreateProductCategory(t *testing.T) {
	createRandomProductCategory(t)
}

func TestGetProductCategory(t *testing.T) {
	productCategory1 := createRandomProductCategory(t)
	productCategory2, err := testQueires.GetProductCategory(context.Background(), productCategory1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, productCategory2)

	require.Equal(t, productCategory1.ID, productCategory2.ID)
	require.Equal(t, productCategory1.Name, productCategory2.Name)
	require.Equal(t, productCategory1.Description, productCategory2.Description)
	require.Equal(t, productCategory1.Active, productCategory2.Active)
	require.Equal(t, productCategory1.CreatedAt, productCategory2.CreatedAt, time.Second)
	require.Equal(t, productCategory1.UpdatedAt, productCategory2.UpdatedAt, time.Second)

	require.False(t, productCategory2.Active)

}

func TestUpdateProductCategory(t *testing.T) {
	productCategory1 := createRandomProductCategory(t)
	arg := UpdateProductCategoryParams{
		ID:     productCategory1.ID,
		Active: true,
	}

	productCategory2, err := testQueires.UpdateProductCategory(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, productCategory2)

	require.Equal(t, productCategory1.ID, productCategory2.ID)
	require.Equal(t, productCategory1.Name, productCategory2.Name)
	require.Equal(t, productCategory1.Description, productCategory2.Description)
	require.Equal(t, arg.Active, productCategory2.Active)
	require.Equal(t, productCategory1.CreatedAt, productCategory2.CreatedAt, time.Second)
	require.NotEqual(t, productCategory1.UpdatedAt, productCategory2.UpdatedAt, time.Second)

	require.True(t, productCategory2.Active)
}

func TestDeleteProductCategory(t *testing.T) {
	productCategory1 := createRandomProductCategory(t)
	err := testQueires.DeleteProductCategory(context.Background(), productCategory1.ID)

	require.NoError(t, err)

	productCategory2, err := testQueires.GetProductCategory(context.Background(), productCategory1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, productCategory2)

}

func TestListProductCategories(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomProductCategory(t)
	}
	arg := ListProductCategoriesParams{
		Limit:  5,
		Offset: 5,
	}

	userCategories, err := testQueires.ListProductCategories(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, userCategories, 5)

	for _, userCategory := range userCategories {
		require.NotEmpty(t, userCategory)

	}
}
