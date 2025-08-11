-- name: CreateUser :one
INSERT INTO users(
	id, username, hashed_password, created_at, updated_at
) VALUES (
	$1, $2, $3, $4, $5
) RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: GetUsers :many
SELECT * FROM users
ORDER BY created_at DESC;


