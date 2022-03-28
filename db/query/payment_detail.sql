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
SELECT "payment_detail".id, "payment_detail".order_id, "payment_detail".amount, 
"payment_detail".provider, "payment_detail".status, "payment_detail".created_at, 
"payment_detail".updated_at
FROM "payment_detail"
LEFT JOIN "order_detail" ON "order_detail".id = "payment_detail".order_id
WHERE "payment_detail".id = $1
-- AND "order_detail".user_id = $2 
LIMIT 1;

-- name: ListPaymentDetails :many
SELECT * FROM "payment_detail"
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdatePaymentDetail :one
UPDATE "payment_detail"
SET order_id = $2,
amount = $3,
provider = $4,
status = $5
WHERE id = $1
RETURNING *;

-- name: DeletePaymentDetail :exec
DELETE FROM "payment_detail"
WHERE id = $1;