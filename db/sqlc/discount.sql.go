// Code generated by sqlc. DO NOT EDIT.
// source: discount.sql

package db

import (
	"context"
)

const createDiscount = `-- name: CreateDiscount :one
INSERT INTO "discount" (
  name,
  description,
  discount_percent
) VALUES (
  $1, $2, $3
)
RETURNING id, name, description, discount_percent, active, created_at, updated_at
`

type CreateDiscountParams struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	DiscountPercent string `json:"discount_percent"`
}

func (q *Queries) CreateDiscount(ctx context.Context, arg CreateDiscountParams) (Discount, error) {
	row := q.db.QueryRowContext(ctx, createDiscount, arg.Name, arg.Description, arg.DiscountPercent)
	var i Discount
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.DiscountPercent,
		&i.Active,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteDiscount = `-- name: DeleteDiscount :exec
DELETE FROM "discount"
WHERE id = $1
`

func (q *Queries) DeleteDiscount(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteDiscount, id)
	return err
}

const getDiscount = `-- name: GetDiscount :one
SELECT id, name, description, discount_percent, active, created_at, updated_at FROM "discount"
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetDiscount(ctx context.Context, id int64) (Discount, error) {
	row := q.db.QueryRowContext(ctx, getDiscount, id)
	var i Discount
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.DiscountPercent,
		&i.Active,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listDiscounts = `-- name: ListDiscounts :many
SELECT id, name, description, discount_percent, active, created_at, updated_at FROM "discount"
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListDiscountsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListDiscounts(ctx context.Context, arg ListDiscountsParams) ([]Discount, error) {
	rows, err := q.db.QueryContext(ctx, listDiscounts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Discount
	for rows.Next() {
		var i Discount
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.DiscountPercent,
			&i.Active,
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

const updateDiscount = `-- name: UpdateDiscount :one
UPDATE "discount"
SET active = $2
WHERE id = $1
RETURNING id, name, description, discount_percent, active, created_at, updated_at
`

type UpdateDiscountParams struct {
	ID     int64 `json:"id"`
	Active bool  `json:"active"`
}

func (q *Queries) UpdateDiscount(ctx context.Context, arg UpdateDiscountParams) (Discount, error) {
	row := q.db.QueryRowContext(ctx, updateDiscount, arg.ID, arg.Active)
	var i Discount
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.DiscountPercent,
		&i.Active,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}