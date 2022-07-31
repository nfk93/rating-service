package match

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/nfk93/rating-service/generated/database"
)

type MatchService struct {
	queries *database.Queries
	db      *sql.DB
}

func (s *MatchService) RegisterMatch(ctx context.Context, gameId uuid.UUID) (database.Match, error) {
	return s.queries.CreateMatch(ctx, database.CreateMatchParams{
		GameID:     gameId,
		HappenedAt: time.Now(),
	})
}

func (s *MatchService) RegisterMatchResult(ctx context.Context, matchID uuid.UUID, playerIDs []uuid.UUID, winnerID uuid.UUID) error {
	if !contains(winnerID, playerIDs) {
		return errors.New("winner id is not in the list of players")
	}

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
	if match.RatingsUpdated || match.IsFinished {
		return errors.New("match already has results registered")
	}

	for _, playerID := range playerIDs {
		_, err := transaction.AddPlayerToMatch(ctx, database.AddPlayerToMatchParams{
			MatchID:  matchID,
			UserID:   playerID,
			IsWinner: playerID == winnerID,
			// TODO: Set current rating
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

func contains(v uuid.UUID, l []uuid.UUID) bool {
	for _, element := range l {
		if v == element {
			return true
		}
	}

	return false
}
