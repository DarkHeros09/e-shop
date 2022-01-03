-- name: CreatePaymentDetail :one
INSERT INTO "payment_detail" (
  order_id,
  amount,
  provider,
  status
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetPaymentDetail :one
SELECT * FROM "payment_detail"
WHERE id = $1 LIMIT 1;

-- name: ListPaymentDetails :many
SELECT * FROM "payment_detail"
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdatePaymentDetail :one
UPDATE "payment_detail"
SET amount = $2,
provider = $3,
status = $4
WHERE id = $1
RETURNING *;

-- name: DeletePaymentDetail :exec
DELETE FROM "payment_detail"
WHERE id = $1;