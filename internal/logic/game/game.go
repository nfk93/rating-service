package game

import (
	"context"
	"github.com/google/uuid"

	"github.com/nfk93/rating-service/sqlc/db"
)

type Service struct {
	queries *db.Queries
}

func New(q *db.Queries) *Service {
	return &Service{
		q,
	}
}

func (s *Service) CreateGame(ctx context.Context, name string, ratingSystem db.RatingSystemEnum) (uuid.UUID, error) {
	game, err := s.queries.CreateGame(ctx, db.CreateGameParams{
		Name:         name,
		RatingSystem: ratingSystem,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return game.ID, err
}
