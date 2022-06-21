-- name: GetUser :one
SELECT * FROM users
WHERE id = sqlc.arg(id) LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY name;

-- name: CreateUser :execresult
INSERT INTO users (
  id, name
) VALUES (
  sqlc.arg(id), sqlc.arg(name)
);

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = sqlc.arg(id);