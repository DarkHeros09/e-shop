-- name: CreateOrderItem :one
INSERT INTO "order_item" (
  order_id,
  product_id,
  quantity
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetOrderItemByID :one
SELECT * FROM "order_item"
WHERE id = $1 LIMIT 1;

-- name: GetOrderItemByOrderDetailID :one
SELECT * FROM "order_item"
WHERE order_id = $1 LIMIT 1;

-- name: ListOrderItems :many
SELECT * FROM "order_item"
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateOrderItem :one
UPDATE "order_item"
SET quantity = $2
WHERE id = $1
RETURNING *;

-- name: DeleteOrderItem :exec
DELETE FROM "order_item"
WHERE id = $1;