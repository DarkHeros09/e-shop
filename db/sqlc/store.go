package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store interface {
	Querier
	FinishedPurchaseTx(ctx context.Context, arg FinishedPurchaseTxParams) (FinishedPurchaseTxResult, error)
}

// Store provides all functions to execute db queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
// method starts with lower case to not be exported so external packages can't call it directly
// we will provide an exported function for each specific transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()

}

// OrderItemTx contains the input parameters of the transfer transaction
type FinishedPurchaseTxParams struct {
	ShoppingSession  ShoppingSession  `json:"shopping_session"`
	CartItem         CartItem         `json:"cart_item"`
	ProductInventory ProductInventory `json:"product_inventory"`
}

// OrderItemTxResult is the result of the transfer transaction
type FinishedPurchaseTxResult struct {
	OrderItem         OrderItem        `json:"order_item"`
	OrderDetail       OrderDetail      `json:"order_detail"`
	PaymentDetail     PaymentDetail    `json:"payment_item"`
	RemainingQuantity ProductInventory `json:"remiainig_quantity"`
}

/* OrderItemTx performs a product transfer from products DB to the user's order_item
once the payments is finished successfully it creates OrderItem record,
substract from/ update the product DB, adds the products to the users' order_item DB,
and update products quantity within a single database transaction.*/
func (store *SQLStore) FinishedPurchaseTx(ctx context.Context, arg FinishedPurchaseTxParams) (FinishedPurchaseTxResult, error) {
	var result FinishedPurchaseTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.PaymentDetail, err = q.CreatePaymentDetail(ctx, CreatePaymentDetailParams{
			Amount:   0,
			Provider: "Unknown",
			Status:   "Pending",
		})
		if err != nil {
			return err
		}

		result.OrderDetail, err = q.CreateOrderDetail(ctx, CreateOrderDetailParams{
			UserID:    arg.ShoppingSession.UserID,
			Total:     arg.ShoppingSession.Total,
			PaymentID: result.PaymentDetail.ID,
		})
		if err != nil {
			// log.Fatal("error1: ", err)
			return err
		}

		result.OrderItem, err = q.CreateOrderItem(ctx, CreateOrderItemParams{
			OrderID:   result.OrderDetail.ID,
			ProductID: arg.CartItem.ProductID,
			Quantity:  arg.CartItem.Quantity,
		})
		if err != nil {
			// log.Fatal("error2: ", err)
			return err
		}

		result.PaymentDetail, err = q.UpdatePaymentDetail(ctx, UpdatePaymentDetailParams{
			ID:       result.OrderDetail.PaymentID,
			OrderID:  result.OrderDetail.ID,
			Amount:   result.OrderItem.Quantity,
			Provider: "Cash",
			Status:   "Finished",
		})
		if err != nil {
			// log.Fatal("error3: ", err)
			return err
		}

		// update ProductInventory method

		// fmt.Println("before: ", arg.ProductInventory.Quantity)
		// fmt.Println("minus: ", result.OrderItem.Quantity)

		result.RemainingQuantity, err = q.UpdateProductQuantity(ctx, UpdateProductQuantityParams{
			ID:       arg.ProductInventory.ID,
			Quantity: -result.OrderItem.Quantity,
		})
		// fmt.Println("after: ", result.RemainingQuantity.Quantity)
		if err != nil {
			return err
		}

		err = q.DeleteCartItem(ctx, arg.CartItem.ID)
		if err != nil {
			return err
		}

		err = q.DeleteShoppingSession(ctx, arg.ShoppingSession.ID)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
