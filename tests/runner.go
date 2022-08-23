package tests

import (
	"database/sql"
	"github.com/nfk93/rating-service/internal/logic/game"
	"github.com/nfk93/rating-service/sqlc/db"

	"github.com/nfk93/rating-service/internal/logic/match"
	"github.com/nfk93/rating-service/internal/logic/rating"
	"github.com/nfk93/rating-service/internal/logic/user"
	"github.com/stretchr/testify/suite"
)

type TestDependencies struct {
	MatchService  *match.MatchService
	RatingService *rating.RatingService
	UserService   *user.UserService
	GameService   *game.Service
	Queries       *db.Queries
	DB            *sql.DB

	TestData TestData
}

type IntegrationTestSuite struct {
	suite.Suite

	TestDependencies *TestDependencies
}

func (s *IntegrationTestSuite) SetupTest() {
	sqldb := SetupTestDB()
	queries := db.New(sqldb)
	matchService := match.New(queries, sqldb)
	ratingService := rating.NewRatingService(queries, sqldb)
	userService := user.NewUserService(queries)
	gameService := game.New(queries)

	data := defaultData()

	s.TestDependencies = &TestDependencies{
		MatchService:  matchService,
		RatingService: ratingService,
		UserService:   userService,
		GameService:   gameService,
		Queries:       queries,
		DB:            sqldb,

		TestData: data,
	}
}

func (s *IntegrationTestSuite) TearDownTest() {
	s.TestDependencies.DB.Close()
}
