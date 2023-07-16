-- name: GetUser :one
SELECT * FROM users
WHERE username = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY username
LIMIT $1
OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (
  username, full_name, hashed_password, email
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: CreateAdmin :one
INSERT INTO users (
  username, full_name, hashed_password, email, level
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE username = $1;

-- name: UpdateHashedPasswordOfUser :one
UPDATE users
  set hashed_password = $2,
      password_changed_at = $3
WHERE username = $1
RETURNING *;