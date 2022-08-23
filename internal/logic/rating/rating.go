package rating

import (
	"context"
	"database/sql"
	"errors"
	"github.com/nfk93/rating-service/sqlc/db"

	"github.com/google/uuid"
	"github.com/nfk93/rating-service/internal/logic/rating/rating_system"
)

type Service struct {
	queries *db.Queries
	db      *sql.DB
}

func NewRatingService(q *db.Queries, db *sql.DB) *Service {
	return &Service{
		queries: q,
		db:      db,
	}
}

// UpdateRatings updates players' rating based on the results of the given match
// TODO: support different rating systems
func (s *Service) UpdateRatings(ctx context.Context, matchID uuid.UUID) error {
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
	if !match.IsFinished {
		return errors.New("match is not finished, can't apply rating updates")
	}
	if match.RatingsUpdated {
		// match rating changes have already been applied so we exit
		return nil
	}

	matchPlayers, err := transaction.GetMatchResult(ctx, matchID)
	if err != nil {
		return err
	}

	ratings := make([]int, len(matchPlayers))
	winnerIndex := -1
	for i, v := range matchPlayers {
		ratings[i] = int(v.CurrentRating)

		if v.IsWinner {
			if winnerIndex != -1 {
				return errors.New("match has two registered winners")
			}

			winnerIndex = i
		}

	}

	// TODO: add support for multiple rating systems
	ratingSystem := rating_system.NewEloRatingSystem()
	ratingDiffs, err := ratingSystem.RatingsDiffs(ratings, winnerIndex)
	if err != nil {
		return err
	}

	for i := 0; i < len(matchPlayers); i++ {
		args := db.ApplyRatingDiffParams{
			UserID:     matchPlayers[i].UserID,
			GameID:     match.GameID,
			Ratingdiff: int32(ratingDiffs[i]),
		}
		err := transaction.ApplyRatingDiff(ctx, args)
		if err != nil {
			return err
		}
	}

	err = transaction.SetMatchRatingsUpdated(ctx, matchID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
