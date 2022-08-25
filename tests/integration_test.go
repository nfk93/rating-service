package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestShit() {
	_, err := s.TestDependencies.UserService.CreateUser(context.TODO(), "peter")
	if err != nil {
		s.T().Fatalf("error: %s", err.Error())
	}
}

func (s *IntegrationTestSuite) TestEloMatch() {
	s.TestData = defaultData(2)

	s.SetupUsers()
	s.SetupGame()
	s.CreateMatch()
	s.JoinMatch()
	s.FinishMatch(1)
	s.UpdateRatings()
	s.AssertRatings([]int{980, 1020})
	s.CreateMatch()
	s.JoinMatch()
	s.FinishMatch(-1)
	s.UpdateRatings()
	s.AssertRatings([]int{982, 1018})
	s.CreateMatch()
	s.JoinMatch()
	s.FinishMatch(0)
	s.UpdateRatings()
	s.AssertRatings([]int{1004, 996})
}
