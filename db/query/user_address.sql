-- name: CreateUserAddress :one
INSERT INTO "user_address" (
  user_id,
  address_line,
  city,
  telephone
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUserAddress :one
SELECT * FROM "user_address"
WHERE id = $1 
And user_id = $2
LIMIT 1;

-- name: ListUserAddresses :many
SELECT * FROM "user_address"
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateUserAddress :one
UPDATE "user_address"
SET address_line = $2,
city = $3,
telephone = $4
WHERE id = $1
RETURNING *;

-- name: UpdateUserAddressByUserID :one
UPDATE "user_address"
SET address_line = $3,
city = $4,
telephone = $5
WHERE user_id = $1
AND id = $2
RETURNING *;

-- name: DeleteUserAddress :exec
DELETE FROM "user_address"
WHERE id = $1;