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
WHERE id = $1 
AND user_id = $2
LIMIT 1;

-- name: ListUserPayments :many
SELECT * FROM "user_payment"
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateUserPayment :one
UPDATE "user_payment"
SET payment_type = $3
WHERE id = $1
AND user_id = $2
RETURNING *;

-- name: DeleteUserPayment :exec
DELETE FROM "user_payment"
WHERE id = $1
AND user_id = $2;