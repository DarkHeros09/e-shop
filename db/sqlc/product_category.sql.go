// Code generated by sqlc. DO NOT EDIT.
// source: product_category.sql

package db

import (
	"context"
)

const createProductCategory = `-- name: CreateProductCategory :one
INSERT INTO "product_category" (
  name,
  description
) VALUES (
  $1, $2
)
RETURNING id, name, description, active, created_at, updated_at
`

type CreateProductCategoryParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (q *Queries) CreateProductCategory(ctx context.Context, arg CreateProductCategoryParams) (ProductCategory, error) {
	row := q.db.QueryRowContext(ctx, createProductCategory, arg.Name, arg.Description)
	var i ProductCategory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Active,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteProductCategory = `-- name: DeleteProductCategory :exec
DELETE FROM "product_category"
WHERE id = $1
`

func (q *Queries) DeleteProductCategory(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteProductCategory, id)
	return err
}

const getProductCategory = `-- name: GetProductCategory :one
SELECT id, name, description, active, created_at, updated_at FROM "product_category"
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetProductCategory(ctx context.Context, id int64) (ProductCategory, error) {
	row := q.db.QueryRowContext(ctx, getProductCategory, id)
	var i ProductCategory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Active,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listProductCategories = `-- name: ListProductCategories :many
SELECT id, name, description, active, created_at, updated_at FROM "product_category"
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListProductCategoriesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListProductCategories(ctx context.Context, arg ListProductCategoriesParams) ([]ProductCategory, error) {
	rows, err := q.db.QueryContext(ctx, listProductCategories, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ProductCategory{}
	for rows.Next() {
		var i ProductCategory
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
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

const updateProductCategory = `-- name: UpdateProductCategory :one
UPDATE "product_category"
SET active = $2
WHERE id = $1
RETURNING id, name, description, active, created_at, updated_at
`

type UpdateProductCategoryParams struct {
	ID     int64 `json:"id"`
	Active bool  `json:"active"`
}

func (q *Queries) UpdateProductCategory(ctx context.Context, arg UpdateProductCategoryParams) (ProductCategory, error) {
	row := q.db.QueryRowContext(ctx, updateProductCategory, arg.ID, arg.Active)
	var i ProductCategory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Active,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
