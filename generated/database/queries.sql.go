// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: queries.sql

package database

import (
	"context"
)

const createGame = `-- name: CreateGame :one

INSERT INTO games (
  id, name
) VALUES (
  $1, $2
)
RETURNING id, name
`

type CreateGameParams struct {
	ID   string
	Name string
}

// ############
// # Ratings
// ############
func (q *Queries) CreateGame(ctx context.Context, arg CreateGameParams) (Game, error) {
	row := q.db.QueryRowContext(ctx, createGame, arg.ID, arg.Name)
	var i Game
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  id, name
) VALUES (
  $1, $2
)
RETURNING id, name
`

type CreateUserParams struct {
	ID   string
	Name string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.ID, arg.Name)
	var i User
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const deleteGame = `-- name: DeleteGame :exec
DELETE FROM games
WHERE id = $1
`

func (q *Queries) DeleteGame(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteGame, id)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getRating = `-- name: GetRating :one

SELECT user_id, game_id, rating FROM ratings
WHERE 
  user_id = $1 AND 
  game_id = $2
LIMIT 1
`

type GetRatingParams struct {
	Userid string
	Gameid string
}

// ############
// # Games
// ############
func (q *Queries) GetRating(ctx context.Context, arg GetRatingParams) (Rating, error) {
	row := q.db.QueryRowContext(ctx, getRating, arg.Userid, arg.Gameid)
	var i Rating
	err := row.Scan(&i.UserID, &i.GameID, &i.Rating)
	return i, err
}

const getUser = `-- name: GetUser :one

SELECT id, name FROM users
WHERE id = $1 LIMIT 1
`

// ############
// # Users
// ############
func (q *Queries) GetUser(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, name FROM users
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const upsertRating = `-- name: UpsertRating :exec
INSERT INTO ratings (
  user_id, game_id, rating
) VALUES (
  $1, $2, $3
)
ON CONFLICT ON CONSTRAINT ratings_pkey DO 
UPDATE SET rating = $3
`

type UpsertRatingParams struct {
	UserID string
	GameID string
	Rating int32
}

func (q *Queries) UpsertRating(ctx context.Context, arg UpsertRatingParams) error {
	_, err := q.db.ExecContext(ctx, upsertRating, arg.UserID, arg.GameID, arg.Rating)
	return err
}
