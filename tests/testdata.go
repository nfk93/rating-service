package tests

import (
	"github.com/nfk93/rating-service/sqlc/db"
	"math/rand"

	"github.com/google/uuid"
)

type TestData struct {
	users map[string]db.User
	game  db.Game

	matchID uuid.UUID
}

func defaultData() TestData {
	users := defaultUsers(4)
	game := defaultGame()

	return TestData{
		users: users,
		game:  game,

		matchID: uuid.New(),
	}
}

func defaultGame() db.Game {
	gameName := randStr(50)
	return db.Game{
		ID:           uuid.New(),
		Name:         gameName,
		RatingSystem: db.RatingSystemEnumElo,
	}
}

func defaultUsers(nUsers int) map[string]db.User {
	users := make(map[string]db.User)
	for i := 0; i < nUsers; i++ {
		name := randStr(50)
		users[name] = db.User{
			ID:   uuid.New(),
			Name: name,
		}
	}
	return users
}

func randStr(n int) string {
	str := make([]rune, n)
	for i := range str {
		str[i] = rune(65 + rand.Intn(26))
	}

	return string(str)
}
