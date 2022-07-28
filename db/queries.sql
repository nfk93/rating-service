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
  sqlc.arg(id), sqlc.arg(name)
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
  id, name
) VALUES (
  sqlc.arg(id), sqlc.arg(name)
)
RETURNING *;

-- name: DeleteGame :exec
DELETE FROM games
WHERE id = sqlc.arg(id);

-- ########################################
-- # GLICKO RATING
-- ########################################

-- name: GetGlickoRating :one
SELECT * FROM glicko_rating
WHERE 
  user_id = sqlc.arg(userId) AND 
  game_id = sqlc.arg(gameId)
LIMIT 1;

-- name: UpsertCurrentGlickoRating :exec
INSERT INTO glicko_rating (
  user_id, game_id, current_rating
) VALUES (
  $1, $2, $3
)
ON CONFLICT DO 
UPDATE SET current_rating = $3;

-- name: UpsertFullGlickoRating :exec
INSERT INTO glicko_rating (
  user_id, game_id, current_rating, glicko_rating, glicko_deviation
) VALUES (
  $1, $2, $3, $4, $5
)
ON CONFLICT DO 
UPDATE SET current_rating = $3, glicko_rating = $4, glicko_deviation = $5;

-- ########################################
-- # ELO
-- ########################################

-- name: GetEloRating :one
SELECT rating FROM elo_rating
WHERE 
  user_id = sqlc.arg(userId) AND 
  game_id = sqlc.arg(gameId);

-- name: UpsertEloRating :exec
INSERT INTO elo_rating (
  user_id, game_id, rating
) VALUES (
  $1, $2, $3
)
ON CONFLICT DO 
UPDATE SET rating = $3;

-- ########################################
-- # MATCHES
-- ########################################

-- name: CreateMatch :one
INSERT INTO matches (
  id, game_id, happened_at
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: AddPlayerToMatch :one
INSERT INTO match_player (
  match_id, user_id, is_winner, score
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetEloMatchResult :many
SELECT B.user_id, B.is_winner, B.score, C.rating, A.happened_at
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

-- name: GetGlickoMatchesAfter :many
SELECT A.match_id as match_id, A.user_id, A.is_winner, A.score, B.current_rating, B.glicko_rating, B.glicko_deviation
FROM (
  SELECT match_id, user_id, is_winner, score 
  FROM match_player
  WHERE match_id IN 
  (
    SELECT matches.id
    FROM match_player INNER JOIN matches ON match_player.match_id = matches.id
    WHERE match_player.user_id = $1 AND matches.game_id = $2 AND matches.happened_at >= $3
  )
) AS A
INNER JOIN glicko_rating AS B ON A.user_id = B.user_id
GROUP BY A.match_id;