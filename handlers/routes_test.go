package handlers_test

import (
	helpers "github.com/duckclick/wing/testing"
	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RoutesTestSuite struct {
	suite.Suite
	conn  *redigomock.Conn
	redis *redis.Pool
}

func (suite *RoutesTestSuite) SetupTest() {
	suite.conn = helpers.CreateRedisConn()
	suite.redis = helpers.CreateRedisPool(suite.conn)
}

func (suite *RoutesTestSuite) TestDrawRoutesExecution() {
	router, err := helpers.CreateRouter()
	assert.Nil(suite.T(), err, "helpers.CreateRouter should succeed")

	if assert.NotNil(suite.T(), router) {
		router.SetAppContextRedis(suite.redis)
		assert.NotNil(suite.T(), router)
		router.DrawRoutes()
	}
}

func TestRoutes(t *testing.T) {
	suite.Run(t, new(RoutesTestSuite))
}
