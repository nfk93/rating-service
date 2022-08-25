package match

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/nfk93/rating-service/internal/logic/rating"
	"time"

	"github.com/google/uuid"
	"github.com/nfk93/rating-service/sqlc/db"
)

type Service struct {
	queries *db.Queries
	db      *sql.DB

	ratingService *rating.Service
}

func New(q *db.Queries, db *sql.DB, ratingService *rating.Service) *Service {
	return &Service{
		q,
		db,
		ratingService,
	}
}

func (s *Service) RegisterMatch(ctx context.Context, gameId uuid.UUID) (db.Match, error) {
	return s.queries.CreateMatch(ctx, db.CreateMatchParams{
		GameID:     gameId,
		HappenedAt: time.Now(),
	})
}

// AddPlayerToMatch fails if the match is already finished
func (s *Service) AddPlayerToMatch(ctx context.Context, userID, matchID uuid.UUID) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	transaction := s.queries.WithTx(tx)

	match, err := transaction.GetMatchForUpdate(ctx, matchID)
	if err != nil {
		return err
	}
	if match.Finished {
		return errors.New("match is already finished")
	}

	_, err = transaction.AddPlayerToMatch(ctx, db.AddPlayerToMatchParams{
		MatchID: matchID,
		UserID:  userID,
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}

// FinishMatch marks the match as finished and updates all participants elo based on their when the match finished
func (s *Service) FinishMatch(ctx context.Context, matchID uuid.UUID, scores map[uuid.UUID]int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	transaction := s.queries.WithTx(tx)

	match, err := transaction.GetMatchForUpdate(ctx, matchID)
	if err != nil {
		return err
	}
	if match.Finished {
		return errors.New("match is already finished")
	}

	players, err := transaction.GetMatchPlayers(ctx, matchID)
	if err != nil {
		return err
	}

	for _, player := range players {
		score, ok := scores[player.UserID]
		if !ok {
			return fmt.Errorf("missing score for player=%s", player.UserID)
		}

		err = transaction.UpdateMatchPlayer(ctx, db.UpdateMatchPlayerParams{
			MatchID: matchID,
			UserID:  player.UserID,
			Score: sql.NullInt32{
				Int32: int32(score),
				Valid: true,
			},
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
