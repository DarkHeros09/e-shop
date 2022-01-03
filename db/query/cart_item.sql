-- name: CreateCartItem :one
INSERT INTO "cart_item" (
  session_id,
  product_id,
  quantity
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetCartItem :one
SELECT * FROM "cart_item"
WHERE id = $1 LIMIT 1;

-- name: ListCartItem :many
SELECT * FROM "cart_item"
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateCartItem :one
UPDATE "cart_item"
SET quantity = $2
WHERE id = $1
RETURNING *;

-- name: DeleteCartItem :exec
DELETE FROM "cart_item"
WHERE id = $1;