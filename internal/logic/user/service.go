package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/nfk93/rating-service/sqlc/db"
	"log"
)

type UserService struct {
	q *db.Queries
}

func NewUserService(q *db.Queries) *UserService {
	return &UserService{
		q: q,
	}
}

func (s *UserService) CreateUser(ctx context.Context, name string) (uuid.UUID, error) {
	user, err := s.q.CreateUser(ctx, name)
	if err != nil {
		log.Printf("error: %s", err.Error())
		return uuid.Nil, err
	}

	return user.ID, nil
}

func (s *UserService) GetUsers(ctx context.Context) ([]db.User, error) {
	users, err := s.q.ListUsers(ctx)
	if err != nil {
		log.Printf("error: %s", err.Error())
		return nil, err
	}

	return users, nil
}
