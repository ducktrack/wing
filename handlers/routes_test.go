package handlers_test

import (
	helpers "github.com/duckclick/wing/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RoutesTestSuite struct {
	suite.Suite
}

func (suite *RoutesTestSuite) TestDrawRoutesExecution() {
	router := helpers.CreateRouter()
	assert.NotNil(suite.T(), router)
	router.DrawRoutes()
}

func TestRoutes(t *testing.T) {
	suite.Run(t, new(RoutesTestSuite))
}
