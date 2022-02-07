-- name: CreateShoppingSession :one
INSERT INTO "shopping_session" (
  user_id,
  total
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetShoppingSession :one
SELECT * FROM "shopping_session"
WHERE id = $1
LIMIT 1;

-- name: ListShoppingSessions :many
SELECT * FROM "shopping_session"
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateShoppingSession :one
UPDATE "shopping_session"
SET total = $2
WHERE id = $1
RETURNING *;

-- name: DeleteShoppingSession :exec
DELETE FROM "shopping_session"
WHERE id = $1;