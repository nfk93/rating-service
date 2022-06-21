package db

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
	"github.com/nfk93/rating-service/generated/database"
)

type Repo struct {
	db      *sql.DB
	queries *database.Queries
}

func NewRepo() (*Repo, error) {
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		return nil, err
	}

	queries := database.New(db)

	return &Repo{
		db:      db,
		queries: queries,
	}, nil
}

func (r *Repo) CreateUser(ctx context.Context, id uuid.UUID, name string) error {
	_, err := r.queries.CreateUser(ctx, database.CreateUserParams{
		ID:   id.String(),
		Name: name,
	})
	return err
}

func (r *Repo) GetUsers(ctx context.Context) ([]database.User, error) {
	return r.queries.ListUsers(ctx)
}
