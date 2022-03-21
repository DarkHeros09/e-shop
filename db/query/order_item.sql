-- name: CreateOrderItem :one
INSERT INTO "order_item" (
  order_id,
  product_id,
  quantity
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetOrderItem :one
SELECT "order_item".id, "order_item".order_id, "order_item".product_id, 
"order_item".quantity, "order_item".created_at, "order_item".updated_at 
FROM "order_item"
LEFT JOIN "order_detail" ON "order_detail".id = "order_item".order_id
WHERE "order_item".id = $1 
AND "order_detail".user_id = $2
LIMIT 1;

-- name: ListOrderItems :many
SELECT "order_item".id, "order_item".order_id, "order_item".product_id, 
"order_item".quantity, "order_item".created_at, "order_item".updated_at
FROM "order_item"
LEFT JOIN "order_detail" ON "order_detail".id = "order_item".order_id
WHERE "order_detail".user_id = $1
LIMIT $2
OFFSET $3;

-- name: UpdateOrderItem :one
UPDATE "order_item"
SET quantity = $2
WHERE id = $1
RETURNING *;

-- name: DeleteOrderItem :exec
DELETE FROM "order_item"
WHERE id = $1;