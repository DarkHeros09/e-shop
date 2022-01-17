// Code generated by sqlc. DO NOT EDIT.
// source: product.sql

package db

import (
	"context"
)

const createProduct = `-- name: CreateProduct :one
INSERT INTO "product" (
  name,
  description,
  sku,
  category_id,
  inventory_id,
  price,
  discount_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, name, description, sku, category_id, inventory_id, price, active, discount_id, created_at, updated_at
`

type CreateProductParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Sku         string `json:"sku"`
	CategoryID  int64  `json:"category_id"`
	InventoryID int64  `json:"inventory_id"`
	Price       string `json:"price"`
	DiscountID  int64  `json:"discount_id"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, createProduct,
		arg.Name,
		arg.Description,
		arg.Sku,
		arg.CategoryID,
		arg.InventoryID,
		arg.Price,
		arg.DiscountID,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Sku,
		&i.CategoryID,
		&i.InventoryID,
		&i.Price,
		&i.Active,
		&i.DiscountID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteProduct = `-- name: DeleteProduct :exec
DELETE FROM "product"
WHERE id = $1
`

func (q *Queries) DeleteProduct(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteProduct, id)
	return err
}

const getProduct = `-- name: GetProduct :one
SELECT id, name, description, sku, category_id, inventory_id, price, active, discount_id, created_at, updated_at FROM "product"
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetProduct(ctx context.Context, id int64) (Product, error) {
	row := q.db.QueryRowContext(ctx, getProduct, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Sku,
		&i.CategoryID,
		&i.InventoryID,
		&i.Price,
		&i.Active,
		&i.DiscountID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listProducts = `-- name: ListProducts :many
SELECT id, name, description, sku, category_id, inventory_id, price, active, discount_id, created_at, updated_at FROM "product"
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListProductsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListProducts(ctx context.Context, arg ListProductsParams) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, listProducts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Sku,
			&i.CategoryID,
			&i.InventoryID,
			&i.Price,
			&i.Active,
			&i.DiscountID,
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

const updateProduct = `-- name: UpdateProduct :one
UPDATE "product"
SET name = $2,
description = $3,
sku = $4,
category_id = $5,
price = $6,
active = $7
WHERE id = $1
RETURNING id, name, description, sku, category_id, inventory_id, price, active, discount_id, created_at, updated_at
`

type UpdateProductParams struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Sku         string `json:"sku"`
	CategoryID  int64  `json:"category_id"`
	Price       string `json:"price"`
	Active      bool   `json:"active"`
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, updateProduct,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.Sku,
		arg.CategoryID,
		arg.Price,
		arg.Active,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Sku,
		&i.CategoryID,
		&i.InventoryID,
		&i.Price,
		&i.Active,
		&i.DiscountID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
