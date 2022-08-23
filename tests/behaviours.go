package tests

import (
	"context"
	"github.com/nfk93/rating-service/sqlc/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

func SetupUser(t *testing.T, deps *TestDependencies) {
	t.Run("SetupUsers", func(t *testing.T) {
		for name, _ := range deps.TestData.users {
			userID, err := deps.UserService.CreateUser(context.TODO(), name)
			assert.NoError(t, err, "error inserting user")

			deps.TestData.users[name] = db.User{
				ID:   userID,
				Name: name,
			}
		}
	})

}

func SetupGame(t *testing.T, deps *TestDependencies) {
	t.Run("SetupGame", func(t *testing.T) {
		gameID, err := deps.GameService.CreateGame(context.TODO(), deps.TestData.game.Name, deps.TestData.game.RatingSystem)
		assert.NoError(t, err, "error creating game")

		deps.TestData.game.ID = gameID
	})
}
