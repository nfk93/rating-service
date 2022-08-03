package user

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/nfk93/rating-service/generated/database"
)

type UserService struct {
	q *database.Queries
}

func NewUserService(q *database.Queries) *UserService {
	return &UserService{
		q: q,
	}
}

func (s *UserService) CreateUser(ctx context.Context, name string) (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Printf("error: %s", err.Error())
		return "", err
	}

	_, err = s.q.CreateUser(ctx, database.CreateUserParams{
		ID:   id,
		Name: name,
	})
	if err != nil {
		log.Printf("error: %s", err.Error())
		return "", err
	}

	return id.String(), nil
}

func (s *UserService) GetUsers(ctx context.Context) ([]database.User, error) {
	users, err := s.q.ListUsers(ctx)
	if err != nil {
		log.Printf("error: %s", err.Error())
		return nil, err
	}

	return users, nil
}
