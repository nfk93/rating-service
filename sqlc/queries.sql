-- noinspection SqlNoDataSourceInspectionForFile

-- ########################################
-- # Users
-- ########################################

-- name: GetUser :one
SELECT * FROM users
WHERE id = sqlc.arg(id) LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (
  id, name
) VALUES (
  gen_random_uuid(), $1
)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = sqlc.arg(id);

-- ########################################
-- # GAMES
-- ########################################

-- name: CreateGame :one
INSERT INTO games (
    id, name, rating_system
) VALUES (
    gen_random_uuid(), $1, $2
)
RETURNING *;

-- name: DeleteGame :exec
DELETE FROM games
WHERE id = sqlc.arg(id);

-- ########################################
-- # ELO
-- ########################################

-- name: GetEloRating :one
SELECT rating FROM elo_rating
WHERE 
  user_id = sqlc.arg(userId) AND 
  game_id = sqlc.arg(gameId);

-- name: GetEloRatingForUpdate :one
SELECT rating FROM elo_rating
WHERE 
  user_id = sqlc.arg(userId) AND 
  game_id = sqlc.arg(gameId)
FOR UPDATE;

-- name: GetEloRatings :many
SELECT rating FROM elo_rating
WHERE 
  user_id = ANY($2::int[]) AND
  game_id = $1;

-- name: GetEloRatingsForUpdate :many
SELECT rating FROM elo_rating
WHERE 
  user_id = ANY($2::int[]) AND
  game_id = $1
FOR UPDATE;

-- name: UpsertEloRating :exec
INSERT INTO elo_rating (
  user_id, game_id, rating
) VALUES (
  $1, $2, $3
)
ON CONFLICT DO 
UPDATE SET rating = $3;

-- name: ApplyRatingDiff :exec
UPDATE elo_rating 
SET rating = rating + sqlc.arg(ratingDiff)
WHERE user_id = $1 AND game_id = $2;

-- ########################################
-- # MATCHES
-- ########################################

-- name: CreateMatch :one
INSERT INTO matches (
  id, game_id, happened_at
) VALUES (
  gen_random_uuid(), $1, $2
)
RETURNING *;

-- name: AddPlayerToMatch :one
INSERT INTO match_player (
  match_id, user_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetMatchForUpdate :one
SELECT * FROM matches WHERE id = $1 FOR UPDATE;

-- name: GetMatchPlayers :many
SELECT * FROM match_player WHERE match_id = $1;

-- name: SetMatchFinished :exec
UPDATE matches SET finished = true WHERE id = $1;

-- name: GetMatchResult :many
SELECT * FROM match_player WHERE match_id = $1;

-- name: UpdateMatchPlayer :exec
UPDATE match_player
SET rating = $3, score = $4
WHERE match_id = $1 AND user_id = $2;

-- name: GetEloMatchResult :many
SELECT B.user_id, B.score, C.rating, A.happened_at
FROM ((
  (SELECT id, happened_at FROM matches WHERE id = $1) as A
  INNER JOIN match_player as B ON A.id = B.match_id)
  INNER JOIN elo_rating as C ON B.user_id = C.user_id
);

-- name: GetMatches :many
SELECT A.id 
FROM (
  (SELECT id, happened_at FROM matches WHERE game_id = $2) as A
  INNER JOIN 
  (SELECT match_id FROM match_player WHERE user_id = $1) AS B 
  ON A.id = B.match_id
);