package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/nfk93/rating-service/generated/database"
)

type Repo struct {
	db *sql.DB
	q  *database.Queries
}

func NewRepo(db *sql.DB) *Repo {
	q := database.New(db)
	return &Repo{
		db: db,
		q:  q,
	}
}

func (r *Repo) CreateUser(ctx context.Context, id uuid.UUID, name string) error {
	_, err := r.q.CreateUser(ctx, database.CreateUserParams{
		ID:   id,
		Name: name,
	})
	return err
}

func (r *Repo) GetUsers(ctx context.Context) ([]database.User, error) {
	q := database.New(r.db)
	return q.ListUsers(ctx)
}

func (r *Repo) CreateMatch(ctx context.Context, gameID uuid.UUID, timestamp time.Time) (database.Match, error) {
	return r.q.CreateMatch(ctx, database.CreateMatchParams{
		GameID:     gameID,
		HappenedAt: timestamp,
	})
}

func (r *Repo) RegisterMatchResults(ctx context.Context, matchID uuid.UUID, playerIDs []uuid.UUID, winnerID uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	transaction := r.q.WithTx(tx)

	match, err := transaction.GetMatchForUpdate(ctx, matchID)
	if err != nil {
		return err
	}
	if match.RatingsUpdated || match.IsFinished {
		return errors.New("match already has results registered")
	}

	for _, playerID := range playerIDs {
		_, err := transaction.AddPlayerToMatch(ctx, database.AddPlayerToMatchParams{
			MatchID:  matchID,
			UserID:   playerID,
			IsWinner: playerID == winnerID,
		})
		if err != nil {
			return err
		}
	}

	err = transaction.SetMatchFinished(ctx, matchID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
