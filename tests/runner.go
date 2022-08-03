package tests

import (
	"database/sql"

	"github.com/nfk93/rating-service/generated/database"
	"github.com/nfk93/rating-service/internal/logic/match"
	"github.com/nfk93/rating-service/internal/logic/rating"
	"github.com/nfk93/rating-service/internal/logic/user"
	"github.com/stretchr/testify/suite"
)

type TestDependencies struct {
	MatchService  *match.MatchService
	RatingService *rating.RatingService
	UserService   *user.UserService
	Queries       *database.Queries
	DB            *sql.DB

	TestData TestData
}

type IntegrationTestSuite struct {
	suite.Suite

	TestDependencies *TestDependencies
}

func (s *IntegrationTestSuite) SetupTest() {
	sqldb := SetupTestDB()
	queries := database.New(sqldb)
	matchService := match.New(queries, sqldb)
	ratingService := rating.NewRatingService(queries, sqldb)
	userService := user.NewUserService(queries)

	data := defaultData()

	s.TestDependencies = &TestDependencies{
		MatchService:  matchService,
		RatingService: ratingService,
		UserService:   userService,
		Queries:       queries,
		DB:            sqldb,

		TestData: data,
	}
}

func (s *IntegrationTestSuite) TearDownTest() {
	s.TestDependencies.DB.Close()
}
