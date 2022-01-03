-- name: CreateUserPayment :one
INSERT INTO "user_payment" (
  user_id,
  payment_type,
  provider,
  account_no,
  expiry
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUserPayment :one
SELECT * FROM "user_payment"
WHERE id = $1 LIMIT 1;

-- name: ListUserPayments :many
SELECT * FROM "user_payment"
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUserPayment :one
UPDATE "user_payment"
SET payment_type = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUserPayment :exec
DELETE FROM "user_payment"
WHERE id = $1;