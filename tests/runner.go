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
	MatchService  *match.Service
	RatingService *rating.Service
	UserService   *user.UserService
	GameService   *game.Service
	Queries       *db.Queries
	DB            *sql.DB
}

type IntegrationTestSuite struct {
	suite.Suite

	TestDependencies *TestDependencies
	TestData         TestData
}

func (s *IntegrationTestSuite) SetupTest() {
	sqldb := SetupTestDB()
	queries := db.New(sqldb)
	ratingService := rating.NewRatingService(queries, sqldb)
	matchService := match.New(queries, sqldb, ratingService)
	userService := user.NewUserService(queries)
	gameService := game.New(queries)

	s.TestDependencies = &TestDependencies{
		MatchService:  matchService,
		RatingService: ratingService,
		UserService:   userService,
		GameService:   gameService,
		Queries:       queries,
		DB:            sqldb,
	}
}

func (s *IntegrationTestSuite) TearDownTest() {
	s.TestDependencies.DB.Close()
}
