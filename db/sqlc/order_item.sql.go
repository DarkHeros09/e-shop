// Code generated by sqlc. DO NOT EDIT.
// source: order_item.sql

package db

import (
	"context"
)

const createOrderItem = `-- name: CreateOrderItem :one
INSERT INTO "order_item" (
  order_id,
  product_id,
  quantity
) VALUES (
  $1, $2, $3
)
RETURNING id, order_id, product_id, quantity, created_at, updated_at
`

type CreateOrderItemParams struct {
	OrderID   int64 `json:"order_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

func (q *Queries) CreateOrderItem(ctx context.Context, arg CreateOrderItemParams) (OrderItem, error) {
	row := q.db.QueryRowContext(ctx, createOrderItem, arg.OrderID, arg.ProductID, arg.Quantity)
	var i OrderItem
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.ProductID,
		&i.Quantity,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteOrderItem = `-- name: DeleteOrderItem :exec
DELETE FROM "order_item"
WHERE id = $1
`

func (q *Queries) DeleteOrderItem(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteOrderItem, id)
	return err
}

const getOrderItemByID = `-- name: GetOrderItemByID :one
SELECT id, order_id, product_id, quantity, created_at, updated_at FROM "order_item"
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetOrderItemByID(ctx context.Context, id int64) (OrderItem, error) {
	row := q.db.QueryRowContext(ctx, getOrderItemByID, id)
	var i OrderItem
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.ProductID,
		&i.Quantity,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getOrderItemByOrderDetailID = `-- name: GetOrderItemByOrderDetailID :one
SELECT id, order_id, product_id, quantity, created_at, updated_at FROM "order_item"
WHERE order_id = $1 LIMIT 1
`

func (q *Queries) GetOrderItemByOrderDetailID(ctx context.Context, orderID int64) (OrderItem, error) {
	row := q.db.QueryRowContext(ctx, getOrderItemByOrderDetailID, orderID)
	var i OrderItem
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.ProductID,
		&i.Quantity,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listOrderItems = `-- name: ListOrderItems :many
SELECT id, order_id, product_id, quantity, created_at, updated_at FROM "order_item"
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListOrderItemsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListOrderItems(ctx context.Context, arg ListOrderItemsParams) ([]OrderItem, error) {
	rows, err := q.db.QueryContext(ctx, listOrderItems, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []OrderItem{}
	for rows.Next() {
		var i OrderItem
		if err := rows.Scan(
			&i.ID,
			&i.OrderID,
			&i.ProductID,
			&i.Quantity,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateOrderItem = `-- name: UpdateOrderItem :one
UPDATE "order_item"
SET quantity = $2
WHERE id = $1
RETURNING id, order_id, product_id, quantity, created_at, updated_at
`

type UpdateOrderItemParams struct {
	ID       int64 `json:"id"`
	Quantity int32 `json:"quantity"`
}

func (q *Queries) UpdateOrderItem(ctx context.Context, arg UpdateOrderItemParams) (OrderItem, error) {
	row := q.db.QueryRowContext(ctx, updateOrderItem, arg.ID, arg.Quantity)
	var i OrderItem
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.ProductID,
		&i.Quantity,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
