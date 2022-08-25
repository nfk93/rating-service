package tests

import (
	"github.com/nfk93/rating-service/sqlc/db"
	"math/rand"

	"github.com/google/uuid"
)

type TestData struct {
	users []db.User
	game  db.Game

	winnerIdx int
	match     *db.Match
}

func defaultData(nUsers int) TestData {
	users := defaultUsers(nUsers)
	game := defaultGame()

	return TestData{
		users: users,
		game:  game,
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

func defaultUsers(nUsers int) []db.User {
	users := make([]db.User, nUsers)
	for i := 0; i < nUsers; i++ {
		defaultUser := defaultUser()
		users[i] = defaultUser
	}
	return users
}

func defaultUser() db.User {
	name := randStr(50)
	defaultUser := db.User{
		ID:   uuid.New(),
		Name: name,
	}
	return defaultUser
}

func randStr(n int) string {
	str := make([]rune, n)
	for i := range str {
		str[i] = rune(65 + rand.Intn(26))
	}

	return string(str)
}
