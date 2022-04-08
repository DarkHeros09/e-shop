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
AND "order_detail".user_id = $2 
LIMIT 1;

-- name: ListPaymentDetails :many
SELECT * FROM "payment_detail"
LEFT JOIN "order_detail" ON "order_detail".id = "payment_detail".order_id
WHERE "order_detail".user_id = $1
ORDER BY "payment_detail".id
LIMIT $2
OFFSET $3;

-- name: UpdatePaymentDetail :one
WITH t1 AS (
SELECT pd.* 
FROM "payment_detail" AS pd 
LEFT JOIN "order_detail" ON "order_detail".payment_id = pd.id 
WHERE pd.id = $1
And user_id= $2
)

UPDATE "payment_detail"
SET order_id = $3,
amount = $4,
provider = $5,
status = $6 
WHERE "payment_detail".id = (SELECT id FROM t1)
RETURNING *;

-- name: DeletePaymentDetail :exec
DELETE FROM "payment_detail"
WHERE id = $1;