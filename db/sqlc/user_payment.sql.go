// Code generated by sqlc. DO NOT EDIT.
// source: user_payment.sql

package db

import (
	"context"
	"time"
)

const createUserPayment = `-- name: CreateUserPayment :one
INSERT INTO "user_payment" (
  user_id,
  payment_type,
  provider,
  account_no,
  expiry
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id, user_id, payment_type, provider, account_no, expiry
`

type CreateUserPaymentParams struct {
	UserID      int64     `json:"user_id"`
	PaymentType string    `json:"payment_type"`
	Provider    string    `json:"provider"`
	AccountNo   int32     `json:"account_no"`
	Expiry      time.Time `json:"expiry"`
}

func (q *Queries) CreateUserPayment(ctx context.Context, arg CreateUserPaymentParams) (UserPayment, error) {
	row := q.db.QueryRowContext(ctx, createUserPayment,
		arg.UserID,
		arg.PaymentType,
		arg.Provider,
		arg.AccountNo,
		arg.Expiry,
	)
	var i UserPayment
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.PaymentType,
		&i.Provider,
		&i.AccountNo,
		&i.Expiry,
	)
	return i, err
}

const deleteUserPayment = `-- name: DeleteUserPayment :exec
DELETE FROM "user_payment"
WHERE id = $1
AND user_id = $2
`

type DeleteUserPaymentParams struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
}

func (q *Queries) DeleteUserPayment(ctx context.Context, arg DeleteUserPaymentParams) error {
	_, err := q.db.ExecContext(ctx, deleteUserPayment, arg.ID, arg.UserID)
	return err
}

const getUserPayment = `-- name: GetUserPayment :one
SELECT id, user_id, payment_type, provider, account_no, expiry FROM "user_payment"
WHERE id = $1 
AND user_id = $2
LIMIT 1
`

type GetUserPaymentParams struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
}

func (q *Queries) GetUserPayment(ctx context.Context, arg GetUserPaymentParams) (UserPayment, error) {
	row := q.db.QueryRowContext(ctx, getUserPayment, arg.ID, arg.UserID)
	var i UserPayment
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.PaymentType,
		&i.Provider,
		&i.AccountNo,
		&i.Expiry,
	)
	return i, err
}

const listUserPayments = `-- name: ListUserPayments :many
SELECT id, user_id, payment_type, provider, account_no, expiry FROM "user_payment"
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListUserPaymentsParams struct {
	UserID int64 `json:"user_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUserPayments(ctx context.Context, arg ListUserPaymentsParams) ([]UserPayment, error) {
	rows, err := q.db.QueryContext(ctx, listUserPayments, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []UserPayment{}
	for rows.Next() {
		var i UserPayment
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.PaymentType,
			&i.Provider,
			&i.AccountNo,
			&i.Expiry,
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

const updateUserPayment = `-- name: UpdateUserPayment :one
UPDATE "user_payment"
SET payment_type = $3
WHERE id = $1
AND user_id = $2
RETURNING id, user_id, payment_type, provider, account_no, expiry
`

type UpdateUserPaymentParams struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	PaymentType string `json:"payment_type"`
}

func (q *Queries) UpdateUserPayment(ctx context.Context, arg UpdateUserPaymentParams) (UserPayment, error) {
	row := q.db.QueryRowContext(ctx, updateUserPayment, arg.ID, arg.UserID, arg.PaymentType)
	var i UserPayment
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.PaymentType,
		&i.Provider,
		&i.AccountNo,
		&i.Expiry,
	)
	return i, err
}
