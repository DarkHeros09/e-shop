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
WHERE id = $1 LIMIT 1;

-- name: ListOrderDetails :many
SELECT * FROM "order_detail"
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateOrderDetail :one
UPDATE "order_detail"
SET total = $2
WHERE id = $1
RETURNING *;

-- name: DeleteOrderDetail :exec
DELETE FROM "order_detail"
WHERE id = $1;