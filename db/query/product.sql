-- name: CreateProduct :one
INSERT INTO "product" (
  name,
  description,
  sku,
  category_id,
  inventory_id,
  price,
  discount_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetProduct :one
SELECT * FROM "product"
WHERE id = $1 LIMIT 1;

-- name: ListProducts :many
SELECT * FROM "product"
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateProduct :one
UPDATE "product"
SET price = $2,
discount_id = $3,
active = $4
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM "product"
WHERE id = $1;