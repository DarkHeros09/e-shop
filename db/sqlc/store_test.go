package db

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFinishedPurchaseTx(t *testing.T) {
	store := NewStore(testDB)

	cartItem1 := createRandomCartItem(t)
	shoppingSession1, err := store.GetShoppingSession(context.Background(), cartItem1.SessionID)
	if err != nil {
		log.Fatal("err is: ", err)
	}
	product, err := store.GetProduct(context.Background(), cartItem1.ProductID)
	if err != nil {
		log.Fatal("err is: ", err)
	}
	productInventory, err := store.GetProductInventory(context.Background(), product.InventoryID)
	if err != nil {
		log.Fatal("err is: ", err)
	}

	// run a concurrent FinishedPurchase
	n := 2

	errs := make(chan error)
	results := make(chan FinishedPurchaseTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.FinishedPurchaseTx(context.Background(), FinishedPurchaseTxParams{
				ShoppingSession:  shoppingSession1,
				CartItem:         cartItem1,
				ProductInventory: productInventory,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check finishedPurchase
		finishedPurchase := result.PaymentDetail
		// remainingQuantity := result.RemainingQuantity.Quantity
		require.NotEmpty(t, finishedPurchase)
		require.Equal(t, cartItem1.Quantity, finishedPurchase.Amount)

		orderDetail, err := store.GetOrderDetail(context.Background(), finishedPurchase.OrderID)
		require.NoError(t, err)
		require.NotEmpty(t, orderDetail)

		// arg := GetOrderItemParams{
		// 	ID:     ,
		// 	UserID: orderDetail.UserID,
		// }

		// orderItem, err := store.GetOrderItem(context.Background(), finishedPurchase.OrderID)
		// require.NoError(t, err)
		// require.NotEmpty(t, orderItem)

		// require.Equal(t, finishedPurchase.Amount, orderItem.Quantity)

		//update product quantity

		// diff1 := remainingQuantity
		// diff2 := productInventory.Quantity - orderItem.Quantity
		// // diff3 := productInventory.Quantity
		// // diff4 := remainingQuantity + orderItem.Quantity

		// require.NoError(t, err)
		// require.Equal(t, diff1, diff2)
		// // require.Equal(t, diff3, diff4)
		// require.True(t, diff1 > 0)

		// check if shopping session, and cartItem is deleted
		shoppingSession2, SSerr := store.GetShoppingSession(context.Background(), shoppingSession1.ID)
		require.Error(t, SSerr)
		require.Empty(t, shoppingSession2)

		cartItem2, CIerr := store.GetShoppingSession(context.Background(), shoppingSession1.ID)
		require.Error(t, CIerr)
		require.Empty(t, cartItem2)

		// arg := GetPaymentDetailParams{
		// 	ID:     finishedPurchase.ID,
		// 	UserID: shoppingSession2.UserID,
		// }
		// _, err = store.GetPaymentDetail(context.Background(), arg)
		// require.NoError(t, err)
	}

}
