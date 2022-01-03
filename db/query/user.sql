-- name: CreateUser :one
INSERT INTO "user" (
  username,
  email,
  password,
  telephone
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM "user"
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM "user"
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE "user"
SET telephone = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM "user"
WHERE id = $1;