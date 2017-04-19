package handlers_test

import (
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/exporters"
	"github.com/duckclick/wing/handlers"
	helpers "github.com/duckclick/wing/testing"
	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RouterTestSuite struct {
	suite.Suite
	conn  *redigomock.Conn
	redis *redis.Pool
}

func (suite *RouterTestSuite) SetupTest() {
	suite.conn = helpers.CreateRedisConn()
	suite.redis = helpers.CreateRedisPool(suite.conn)
}

func (suite *RouterTestSuite) TestNewRouter() {
	appConfig := helpers.CreateBasicConfig()
	appConfig.Exporter = "file"
	appConfig.FileExporter = config.FileExporter{
		Folder: "/tmp/test/track_entries",
	}

	fileExporter, _ := exporters.Lookup(&appConfig)
	router, err := handlers.NewRouter(&appConfig, fileExporter)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), router.AppContext().Redis)
}

func (suite *RouterTestSuite) TestNewRouterWhenCertificatesAreNotConfigured() {
	appConfig := helpers.CreateBasicConfig()
	appConfig.JWEPrivateKeyFile = "wrong"

	appConfig.Exporter = "file"
	appConfig.FileExporter = config.FileExporter{
		Folder: "/tmp/test/track_entries",
	}

	fileExporter, _ := exporters.Lookup(&appConfig)
	_, err := handlers.NewRouter(&appConfig, fileExporter)
	assert.NotNil(suite.T(), err)
}

func TestRouter(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}
