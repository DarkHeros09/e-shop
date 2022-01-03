// Code generated by sqlc. DO NOT EDIT.
// source: cart_item.sql

package db

import (
	"context"
)

const createCartItem = `-- name: CreateCartItem :one
INSERT INTO "cart_item" (
  session_id,
  product_id,
  quantity
) VALUES (
  $1, $2, $3
)
RETURNING id, session_id, product_id, quantity, created_at, updated_at
`

type CreateCartItemParams struct {
	SessionID int64 `json:"session_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

func (q *Queries) CreateCartItem(ctx context.Context, arg CreateCartItemParams) (CartItem, error) {
	row := q.db.QueryRowContext(ctx, createCartItem, arg.SessionID, arg.ProductID, arg.Quantity)
	var i CartItem
	err := row.Scan(
		&i.ID,
		&i.SessionID,
		&i.ProductID,
		&i.Quantity,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteCartItem = `-- name: DeleteCartItem :exec
DELETE FROM "cart_item"
WHERE id = $1
`

func (q *Queries) DeleteCartItem(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteCartItem, id)
	return err
}

const getCartItem = `-- name: GetCartItem :one
SELECT id, session_id, product_id, quantity, created_at, updated_at FROM "cart_item"
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetCartItem(ctx context.Context, id int64) (CartItem, error) {
	row := q.db.QueryRowContext(ctx, getCartItem, id)
	var i CartItem
	err := row.Scan(
		&i.ID,
		&i.SessionID,
		&i.ProductID,
		&i.Quantity,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listCartItem = `-- name: ListCartItem :many
SELECT id, session_id, product_id, quantity, created_at, updated_at FROM "cart_item"
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListCartItemParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListCartItem(ctx context.Context, arg ListCartItemParams) ([]CartItem, error) {
	rows, err := q.db.QueryContext(ctx, listCartItem, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CartItem
	for rows.Next() {
		var i CartItem
		if err := rows.Scan(
			&i.ID,
			&i.SessionID,
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

const updateCartItem = `-- name: UpdateCartItem :one
UPDATE "cart_item"
SET quantity = $2
WHERE id = $1
RETURNING id, session_id, product_id, quantity, created_at, updated_at
`

type UpdateCartItemParams struct {
	ID       int64 `json:"id"`
	Quantity int32 `json:"quantity"`
}

func (q *Queries) UpdateCartItem(ctx context.Context, arg UpdateCartItemParams) (CartItem, error) {
	row := q.db.QueryRowContext(ctx, updateCartItem, arg.ID, arg.Quantity)
	var i CartItem
	err := row.Scan(
		&i.ID,
		&i.SessionID,
		&i.ProductID,
		&i.Quantity,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
