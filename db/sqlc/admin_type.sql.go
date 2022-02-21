// Code generated by sqlc. DO NOT EDIT.
// source: admin_type.sql

package db

import (
	"context"
)

const createAdminType = `-- name: CreateAdminType :one
INSERT INTO "admin_type" (
  admin_type
) VALUES ( $1 )
RETURNING id, admin_type, created_at, updated_at
`

func (q *Queries) CreateAdminType(ctx context.Context, adminType string) (AdminType, error) {
	row := q.db.QueryRowContext(ctx, createAdminType, adminType)
	var i AdminType
	err := row.Scan(
		&i.ID,
		&i.AdminType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAdminTypeByID = `-- name: DeleteAdminTypeByID :exec
DELETE FROM "admin_type"
WHERE id = $1
`

func (q *Queries) DeleteAdminTypeByID(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAdminTypeByID, id)
	return err
}

const deleteAdminTypeByType = `-- name: DeleteAdminTypeByType :exec
DELETE FROM "admin_type"
WHERE admin_type = $1
`

func (q *Queries) DeleteAdminTypeByType(ctx context.Context, adminType string) error {
	_, err := q.db.ExecContext(ctx, deleteAdminTypeByType, adminType)
	return err
}

const getAdminType = `-- name: GetAdminType :one
SELECT id, admin_type, created_at, updated_at FROM "admin_type"
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAdminType(ctx context.Context, id int64) (AdminType, error) {
	row := q.db.QueryRowContext(ctx, getAdminType, id)
	var i AdminType
	err := row.Scan(
		&i.ID,
		&i.AdminType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listAdminTypes = `-- name: ListAdminTypes :many
SELECT id, admin_type, created_at, updated_at FROM "admin_type"
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListAdminTypesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListAdminTypes(ctx context.Context, arg ListAdminTypesParams) ([]AdminType, error) {
	rows, err := q.db.QueryContext(ctx, listAdminTypes, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AdminType{}
	for rows.Next() {
		var i AdminType
		if err := rows.Scan(
			&i.ID,
			&i.AdminType,
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

const updateAdminType = `-- name: UpdateAdminType :one
UPDATE "admin_type"
SET admin_type = $2
WHERE id = $1
RETURNING id, admin_type, created_at, updated_at
`

type UpdateAdminTypeParams struct {
	ID        int64  `json:"id"`
	AdminType string `json:"admin_type"`
}

func (q *Queries) UpdateAdminType(ctx context.Context, arg UpdateAdminTypeParams) (AdminType, error) {
	row := q.db.QueryRowContext(ctx, updateAdminType, arg.ID, arg.AdminType)
	var i AdminType
	err := row.Scan(
		&i.ID,
		&i.AdminType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}