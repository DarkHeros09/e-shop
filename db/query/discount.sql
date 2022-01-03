-- name: CreateDiscount :one
INSERT INTO "discount" (
  name,
  description,
  discount_percent
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetDiscount :one
SELECT * FROM "discount"
WHERE id = $1 LIMIT 1;

-- name: ListDiscounts :many
SELECT * FROM "discount"
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateDiscount :one
UPDATE "discount"
SET active = $2
WHERE id = $1
RETURNING *;

-- name: DeleteDiscount :exec
DELETE FROM "discount"
WHERE id = $1;