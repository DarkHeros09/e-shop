// Code generated by sqlc. DO NOT EDIT.
// source: payment_detail.sql

package db

import (
	"context"
)

const createPaymentDetail = `-- name: CreatePaymentDetail :one
INSERT INTO "payment_detail" (
  amount,
  provider,
  status
) VALUES (
  $1, $2, $3
)
RETURNING id, order_id, amount, provider, status, created_at, updated_at
`

type CreatePaymentDetailParams struct {
	Amount   int32  `json:"amount"`
	Provider string `json:"provider"`
	Status   string `json:"status"`
}

func (q *Queries) CreatePaymentDetail(ctx context.Context, arg CreatePaymentDetailParams) (PaymentDetail, error) {
	row := q.db.QueryRowContext(ctx, createPaymentDetail, arg.Amount, arg.Provider, arg.Status)
	var i PaymentDetail
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.Amount,
		&i.Provider,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deletePaymentDetail = `-- name: DeletePaymentDetail :exec
DELETE FROM "payment_detail"
WHERE id = $1
`

func (q *Queries) DeletePaymentDetail(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deletePaymentDetail, id)
	return err
}

const getPaymentDetail = `-- name: GetPaymentDetail :one
SELECT id, order_id, amount, provider, status, created_at, updated_at FROM "payment_detail"
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetPaymentDetail(ctx context.Context, id int64) (PaymentDetail, error) {
	row := q.db.QueryRowContext(ctx, getPaymentDetail, id)
	var i PaymentDetail
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.Amount,
		&i.Provider,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listPaymentDetails = `-- name: ListPaymentDetails :many
SELECT id, order_id, amount, provider, status, created_at, updated_at FROM "payment_detail"
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListPaymentDetailsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListPaymentDetails(ctx context.Context, arg ListPaymentDetailsParams) ([]PaymentDetail, error) {
	rows, err := q.db.QueryContext(ctx, listPaymentDetails, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PaymentDetail
	for rows.Next() {
		var i PaymentDetail
		if err := rows.Scan(
			&i.ID,
			&i.OrderID,
			&i.Amount,
			&i.Provider,
			&i.Status,
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

const updatePaymentDetail = `-- name: UpdatePaymentDetail :one
UPDATE "payment_detail"
SET order_id = $2,
amount = $3,
provider = $4,
status = $5
WHERE id = $1
RETURNING id, order_id, amount, provider, status, created_at, updated_at
`

type UpdatePaymentDetailParams struct {
	ID       int64  `json:"id"`
	OrderID  int64  `json:"order_id"`
	Amount   int32  `json:"amount"`
	Provider string `json:"provider"`
	Status   string `json:"status"`
}

func (q *Queries) UpdatePaymentDetail(ctx context.Context, arg UpdatePaymentDetailParams) (PaymentDetail, error) {
	row := q.db.QueryRowContext(ctx, updatePaymentDetail,
		arg.ID,
		arg.OrderID,
		arg.Amount,
		arg.Provider,
		arg.Status,
	)
	var i PaymentDetail
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.Amount,
		&i.Provider,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}