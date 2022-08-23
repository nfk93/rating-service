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
	s.Run(
		SetupUser,
		SetupGame,
	)
}

type Behaviour func(t *testing.T, deps *TestDependencies)

func (s *IntegrationTestSuite) Run(behaviours ...Behaviour) {
	for _, b := range behaviours {
		b(s.T(), s.TestDependencies)
	}
}
