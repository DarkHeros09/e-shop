-- name: CreateProductInventory :one
INSERT INTO "product_inventory" (
  quantity
) VALUES (
  $1
)
RETURNING *;

-- name: GetProductInventory :one
SELECT * FROM "product_inventory"
WHERE id = $1 LIMIT 1;

-- name: ListProductInventories :many
SELECT * FROM "product_inventory"
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateProductInventory :one
UPDATE "product_inventory"
SET quantity = $2,
active= $3
WHERE id = $1
RETURNING *;

-- name: DeleteProductInventory :exec
DELETE FROM "product_inventory"
WHERE id = $1;