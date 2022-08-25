package rating

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nfk93/rating-service/internal/logic/rating/rating_system"
	"github.com/nfk93/rating-service/sqlc/db"
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

func (s *Service) UpdateRatingsForGame(ctx context.Context, matchID uuid.UUID) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	transaction := s.queries.WithTx(tx)

	players, game, err := getMatchData(ctx, matchID, transaction)
	if err != nil {
		return err
	}

	switch game.RatingSystem {
	case db.RatingSystemEnumElo:
		err := UpdateEloRating(players)
		if err != nil {
			return err
		}
	default:
		return errors.New("unrecognized rating system")
	}

	for _, player := range players {
		if !player.RatingBefore.Valid || !player.RatingChange.Valid {
			return fmt.Errorf("player (id=%s) rating not set", player.UserID)
		}

		err := transaction.UpdateMatchPlayer(ctx, db.UpdateMatchPlayerParams{
			MatchID:      matchID,
			UserID:       player.UserID,
			Score:        player.Score,
			RatingBefore: player.RatingBefore,
			RatingChange: player.RatingChange,
		})
		if err != nil {
			return err
		}

		err = transaction.UpsertEloRating(ctx, db.UpsertEloRatingParams{
			UserID: player.UserID,
			GameID: game.ID,
			Rating: player.RatingBefore.Int32 + player.RatingChange.Int32,
		})
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

// Get game and match players with ratings (if a rating exists for this player)
// It gets the match and ratings for update
// TODO: Get all the match data in one query
func getMatchData(ctx context.Context, matchID uuid.UUID, tx *db.Queries) ([]db.MatchPlayer, db.Game, error) {
	match, err := tx.GetMatchForUpdate(ctx, matchID)
	if err != nil {
		return nil, db.Game{}, err
	}

	players, err := tx.GetMatchPlayers(ctx, match.ID)
	if err != nil {
		return nil, db.Game{}, err
	}

	playerMap := make(map[uuid.UUID]*db.MatchPlayer)
	playerUUIDs := make([]uuid.UUID, len(players))
	for i := range players {
		player := players[i]
		playerUUIDs[i] = player.UserID
		playerMap[player.UserID] = &player
	}

	playerRatings, err := tx.GetRatingsForUpdate(ctx, db.GetRatingsForUpdateParams{
		GameID:  uuid.UUID{},
		Column2: playerUUIDs,
	})
	for _, rating := range playerRatings {
		playerMap[rating.UserID].Score = sql.NullInt32{
			Int32: rating.Rating,
			Valid: true,
		}
	}

	game, err := tx.GetGame(ctx, match.GameID)
	if err != nil {
		return nil, db.Game{}, err
	}
	return players, game, nil
}

func UpdateEloRating(players []db.MatchPlayer) error {
	if len(players) != 2 {
		return errors.New("number of players must be 2 for elo matches")
	}

	winnerScore := int32(0)
	winnerIdx := -1
	for i, player := range players {
		if player.Score.Valid && player.Score.Int32 > winnerScore {
			winnerIdx = i
			winnerScore = player.Score.Int32
		}
	}

	eloRatingSystem := rating_system.NewEloRatingSystem()
	ratingA := getRatingFallback(players[0], eloRatingSystem.DefaultRating())
	ratingB := getRatingFallback(players[1], eloRatingSystem.DefaultRating())

	result := winnerIdxToResult(winnerIdx)
	newRatingA, newRatingB := eloRatingSystem.CalculateUpdatedEloRatings(ratingA, ratingB, result)
	players[0].RatingBefore = sql.NullInt32{
		Int32: int32(ratingA),
		Valid: true,
	}
	players[0].RatingChange = sql.NullInt32{
		Int32: int32(newRatingA - ratingA),
		Valid: true,
	}

	players[1].RatingBefore = sql.NullInt32{
		Int32: int32(ratingB),
		Valid: true,
	}
	players[1].RatingChange = sql.NullInt32{
		Int32: int32(newRatingB - ratingB),
		Valid: true,
	}

	return nil
}

func winnerIdxToResult(idx int) float64 {
	if idx == 0 {
		return 1
	}
	if idx == 1 {
		return 0
	}

	return 0.5
}

func getRatingFallback(player db.MatchPlayer, fallback int) int {
	if player.RatingBefore.Valid {
		return int(player.RatingBefore.Int32)
	}

	return fallback
}
