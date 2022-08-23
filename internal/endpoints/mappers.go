package endpoints

import (
	"github.com/nfk93/rating-service/generated/api"
	"github.com/nfk93/rating-service/sqlc/db"
)

func mapUsers(users []db.User) []api.User {
	result := make([]api.User, len(users))
	for i, v := range users {
		result[i] = mapUser(v)
	}
	return result
}

func mapUser(user db.User) api.User {
	return api.User{
		Id:   user.ID.String(),
		Name: user.Name,
	}
}
