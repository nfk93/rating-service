package endpoints

import (
	"github.com/nfk93/rating-service/generated/api"
	"github.com/nfk93/rating-service/generated/database"
)

func mapUsers(users []database.User) []api.User {
	result := make([]api.User, len(users))
	for i, v := range users {
		result[i] = mapUser(v)
	}
	return result
}

func mapUser(user database.User) api.User {
	return api.User{
		Id:   user.ID,
		Name: user.Name,
	}
}
