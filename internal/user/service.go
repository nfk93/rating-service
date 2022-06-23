package user

import (
	"context"
	"log"

	"github.com/gofrs/uuid"
	"github.com/nfk93/rating-service/db"
	"github.com/nfk93/rating-service/generated/database"
)

type UserService struct {
	repo *db.Repo
}

func NewUserService(repo *db.Repo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, name string) (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		log.Printf("error: %s", err.Error())
		return "", err
	}

	err = s.repo.CreateUser(ctx, id, name)
	if err != nil {
		log.Printf("error: %s", err.Error())
		return "", err
	}

	return id.String(), nil
}

func (s *UserService) GetUsers(ctx context.Context) ([]database.User, error) {
	users, err := s.repo.GetUsers(ctx)
	if err != nil {
		log.Printf("error: %s", err.Error())
		return nil, err
	}

	return users, nil
}
