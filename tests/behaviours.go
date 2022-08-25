package tests

import (
	"context"
	"github.com/google/uuid"
	"github.com/nfk93/rating-service/sqlc/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

func (s *IntegrationTestSuite) SetupUsers() {
	s.T().Run("SetupUsers", func(t *testing.T) {
		for i := range s.TestData.users {
			userID, err := s.TestDependencies.UserService.CreateUser(context.TODO(), s.TestData.users[i].Name)
			assert.NoError(t, err, "error inserting user")

			s.TestData.users[i].ID = userID
		}
	})

}

func (s *IntegrationTestSuite) SetupGame() {
	s.T().Run("SetupGame", func(t *testing.T) {
		gameID, err := s.TestDependencies.GameService.CreateGame(
			context.TODO(),
			s.TestData.game.Name,
			s.TestData.game.RatingSystem,
		)
		assert.NoError(t, err, "error creating game")

		s.TestData.game.ID = gameID
	})
}

func (s *IntegrationTestSuite) CreateMatch() {
	s.T().Run("CreateMatch", func(t *testing.T) {
		match, err := s.TestDependencies.MatchService.RegisterMatch(
			context.TODO(),
			s.TestData.game.ID,
		)
		assert.NoError(t, err, "error creating game")

		s.TestData.match = &match
	})
}

func (s *IntegrationTestSuite) JoinMatch() {
	s.T().Run("JoinMatch", func(t *testing.T) {
		for _, player := range s.TestData.users {
			err := s.TestDependencies.MatchService.AddPlayerToMatch(
				context.TODO(),
				player.ID,
				s.TestData.match.ID,
			)
			assert.NoError(t, err, "error joining game")
		}
	})
}

func (s *IntegrationTestSuite) FinishMatch(winnerIdx int) {
	s.T().Run("FinishMatch", func(t *testing.T) {
		s.TestData.winnerIdx = winnerIdx

		scores := make(map[uuid.UUID]int)
		for i, player := range s.TestData.users {
			score := 0
			if i == winnerIdx {
				score = 1
			}

			scores[player.ID] = score
		}

		err := s.TestDependencies.MatchService.FinishMatch(
			context.TODO(),
			s.TestData.match.ID,
			scores,
		)
		assert.NoError(t, err, "error joining game")
	})
}

func (s *IntegrationTestSuite) UpdateRatings() {
	s.T().Run("UpdateRatings", func(t *testing.T) {
		err := s.TestDependencies.RatingService.UpdateRatingsForGame(context.TODO(), s.TestData.match.ID)
		assert.NoError(t, err, "error updating ratings")
	})
}

func (s *IntegrationTestSuite) AssertRatings(ratings []int) {
	s.T().Run("AssertRatings", func(t *testing.T) {
		for idx, rating := range ratings {
			player := s.TestData.users[idx]

			actualRating, err := s.TestDependencies.Queries.GetRating(context.TODO(), db.GetRatingParams{
				Userid: player.ID,
				Gameid: s.TestData.game.ID,
			})
			assert.NoError(t, err, "error getting rating")

			assert.Equal(t, rating, int(actualRating))
		}
	})
}
