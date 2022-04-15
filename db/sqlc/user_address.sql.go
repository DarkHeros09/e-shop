// Code generated by sqlc. DO NOT EDIT.
// source: user_address.sql

package db

import (
	"context"
)

const createUserAddress = `-- name: CreateUserAddress :one
INSERT INTO "user_address" (
  user_id,
  address_line,
  city,
  telephone
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, user_id, address_line, city, telephone
`

type CreateUserAddressParams struct {
	UserID      int64  `json:"user_id"`
	AddressLine string `json:"address_line"`
	City        string `json:"city"`
	Telephone   int32  `json:"telephone"`
}

func (q *Queries) CreateUserAddress(ctx context.Context, arg CreateUserAddressParams) (UserAddress, error) {
	row := q.db.QueryRowContext(ctx, createUserAddress,
		arg.UserID,
		arg.AddressLine,
		arg.City,
		arg.Telephone,
	)
	var i UserAddress
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.AddressLine,
		&i.City,
		&i.Telephone,
	)
	return i, err
}

const deleteUserAddress = `-- name: DeleteUserAddress :exec
DELETE FROM "user_address"
WHERE id = $1
`

func (q *Queries) DeleteUserAddress(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUserAddress, id)
	return err
}

const getUserAddress = `-- name: GetUserAddress :one
SELECT id, user_id, address_line, city, telephone FROM "user_address"
WHERE id = $1 
And user_id = $2
LIMIT 1
`

type GetUserAddressParams struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
}

func (q *Queries) GetUserAddress(ctx context.Context, arg GetUserAddressParams) (UserAddress, error) {
	row := q.db.QueryRowContext(ctx, getUserAddress, arg.ID, arg.UserID)
	var i UserAddress
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.AddressLine,
		&i.City,
		&i.Telephone,
	)
	return i, err
}

const listUserAddresses = `-- name: ListUserAddresses :many
SELECT id, user_id, address_line, city, telephone FROM "user_address"
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListUserAddressesParams struct {
	UserID int64 `json:"user_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUserAddresses(ctx context.Context, arg ListUserAddressesParams) ([]UserAddress, error) {
	rows, err := q.db.QueryContext(ctx, listUserAddresses, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []UserAddress{}
	for rows.Next() {
		var i UserAddress
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.AddressLine,
			&i.City,
			&i.Telephone,
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

const updateUserAddress = `-- name: UpdateUserAddress :one
UPDATE "user_address"
SET address_line = $2,
city = $3,
telephone = $4
WHERE id = $1
RETURNING id, user_id, address_line, city, telephone
`

type UpdateUserAddressParams struct {
	ID          int64  `json:"id"`
	AddressLine string `json:"address_line"`
	City        string `json:"city"`
	Telephone   int32  `json:"telephone"`
}

func (q *Queries) UpdateUserAddress(ctx context.Context, arg UpdateUserAddressParams) (UserAddress, error) {
	row := q.db.QueryRowContext(ctx, updateUserAddress,
		arg.ID,
		arg.AddressLine,
		arg.City,
		arg.Telephone,
	)
	var i UserAddress
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.AddressLine,
		&i.City,
		&i.Telephone,
	)
	return i, err
}

const updateUserAddressByUserID = `-- name: UpdateUserAddressByUserID :one
UPDATE "user_address"
SET address_line = $3,
city = $4,
telephone = $5
WHERE user_id = $1
AND id = $2
RETURNING id, user_id, address_line, city, telephone
`

type UpdateUserAddressByUserIDParams struct {
	UserID      int64  `json:"user_id"`
	ID          int64  `json:"id"`
	AddressLine string `json:"address_line"`
	City        string `json:"city"`
	Telephone   int32  `json:"telephone"`
}

func (q *Queries) UpdateUserAddressByUserID(ctx context.Context, arg UpdateUserAddressByUserIDParams) (UserAddress, error) {
	row := q.db.QueryRowContext(ctx, updateUserAddressByUserID,
		arg.UserID,
		arg.ID,
		arg.AddressLine,
		arg.City,
		arg.Telephone,
	)
	var i UserAddress
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.AddressLine,
		&i.City,
		&i.Telephone,
	)
	return i, err
}
