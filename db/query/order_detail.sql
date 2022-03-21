-- name: CreateOrderDetail :one
INSERT INTO "order_detail" (
  user_id,
  total,
  payment_id
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetOrderDetail :one
SELECT * FROM "order_detail"
WHERE id = $1 
LIMIT 1;

-- name: ListOrderDetails :many
SELECT * FROM "order_detail"
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateOrderDetail :one
UPDATE "order_detail"
SET total = $2
WHERE id = $1
RETURNING *;

-- name: DeleteOrderDetail :exec
DELETE FROM "order_detail"
WHERE id = $1;