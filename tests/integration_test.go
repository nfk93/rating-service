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
