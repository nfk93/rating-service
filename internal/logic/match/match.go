package match

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/nfk93/rating-service/sqlc/db"
)

type Service struct {
	queries *db.Queries
	db      *sql.DB
}

func New(q *db.Queries, db *sql.DB) *Service {
	return &Service{
		q,
		db,
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

	_, err = s.queries.AddPlayerToMatch(ctx, db.AddPlayerToMatchParams{
		MatchID: matchID,
		UserID:  userID,
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}

// FinishMatch marks the match as finished and updates all participants elo based on their when the match finished
func (s *Service) FinishMatch(ctx context.Context, matchID uuid.UUID, winnerID uuid.UUID) error {
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

	ratings := make(map[uuid.UUID]int)
	for _, player := range players {
		// TODO: default rating if not found
		rating, err := transaction.GetEloRatingForUpdate(ctx, db.GetEloRatingForUpdateParams{
			Userid: player.UserID,
			Gameid: match.GameID,
		})
		if err != nil {
			return err
		}

		score := sql.NullInt32{
			Int32: 0,
			Valid: true,
		}
		if player.UserID == winnerID {
			score.Int32 = 1
		}

		ratings[player.UserID] = int(rating)
		err = transaction.UpdateMatchPlayer(ctx, db.UpdateMatchPlayerParams{
			MatchID: matchID,
			UserID:  player.UserID,
			Rating: sql.NullInt32{
				Int32: rating,
				Valid: true,
			},
			Score: score,
		})
		if err != nil {
			return err
		}
	}

	// TODO: calc new ratings and update

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
