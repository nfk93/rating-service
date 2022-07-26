-- ############
-- # Users
-- ############

-- name: GetUser :one
SELECT * FROM users
WHERE id = sqlc.arg(id) LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (
  id, name
) VALUES (
  sqlc.arg(id), sqlc.arg(name)
)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = sqlc.arg(id);

-- ############
-- # Ratings
-- ############

-- name: CreateGame :one
INSERT INTO games (
  id, name
) VALUES (
  sqlc.arg(id), sqlc.arg(name)
)
RETURNING *;

-- name: DeleteGame :exec
DELETE FROM games
WHERE id = sqlc.arg(id);

-- ############
-- # Games
-- ############

-- name: GetRating :one
SELECT * FROM ratings
WHERE 
  user_id = sqlc.arg(userId) AND 
  game_id = sqlc.arg(gameId)
LIMIT 1;

-- name: UpsertRating :exec
INSERT INTO ratings (
  user_id, game_id, rating
) VALUES (
  $1, $2, $3
)
ON CONFLICT ON CONSTRAINT ratings_pkey DO 
UPDATE SET rating = $3;
