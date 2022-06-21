package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/nfk93/rating-service/db"
	"github.com/nfk93/rating-service/generated/database"
)

type Service struct {
	repo *db.Repo
}

func (s *Service) CreateUser(ctx context.Context, name string) (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	err = s.repo.CreateUser(ctx, id, name)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (s *Service) GetUsers(ctx context.Context) ([]database.User, error) {
	return s.repo.GetUsers(ctx)
}
