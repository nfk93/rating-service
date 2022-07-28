package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/nfk93/rating-service/generated/database"
)

type Repo struct {
	queries *database.Queries
}

func NewRepo(queries *database.Queries) *Repo {
	return &Repo{
		queries: queries,
	}
}

func (r *Repo) CreateUser(ctx context.Context, id uuid.UUID, name string) error {
	_, err := r.queries.CreateUser(ctx, database.CreateUserParams{
		ID:   id,
		Name: name,
	})
	return err
}

func (r *Repo) GetUsers(ctx context.Context) ([]database.User, error) {
	return r.queries.ListUsers(ctx)
}
