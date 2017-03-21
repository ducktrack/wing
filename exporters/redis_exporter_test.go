package exporters_test

import (
	"errors"
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/events"
	"github.com/duckclick/wing/exporters"
	helpers "github.com/duckclick/wing/testing"
	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RedisExporterTestSuite struct {
	suite.Suite
	recordID         string
	exporter         *exporters.RedisExporter
	mockedConnection *redigomock.Conn
}

func (suite *RedisExporterTestSuite) SetupTest() {
	suite.recordID = uuid.NewV4().String()
	suite.exporter = &exporters.RedisExporter{
		Config: config.RedisExporter{
			Host: "foo",
			Port: 1234,
		},
	}

	suite.mockedConnection = redigomock.NewConn()
	suite.exporter.Pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return suite.mockedConnection, nil
		},
	}
}

func (suite *RedisExporterTestSuite) TestRedisExport() {
	defer suite.exporter.Stop()
	trackDOM, err := events.TrackDOMFromJSON(events.Event{
		CreatedAt:  1487696788863,
		URL:        "http://example.org/some/path",
		RawPayload: helpers.CreateRawMessage(`{"markup": %s}`, helpers.Base64BlankMarkup),
	})

	assert.Nil(suite.T(), err, "events.TrackDOMFromJSON() should succeed")
	eventJSON, err := trackDOM.ToJSON()
	assert.Nil(suite.T(), err, "trackDOM.ToJSON() should succeed")

	suite.mockedConnection.
		Command("HSET", suite.recordID, "1487696788863", eventJSON).
		Expect(nil)

	err = suite.exporter.Export(trackDOM, suite.recordID)
	assert.Nil(suite.T(), err, "RedisExporter#Export should succeed")
}

func (suite *RedisExporterTestSuite) TestExportReturnsErrorOnRedisError() {
	defer suite.exporter.Stop()
	trackDOM, err := events.TrackDOMFromJSON(events.Event{
		CreatedAt:  1487696788863,
		URL:        "http://example.org/some/path",
		RawPayload: helpers.CreateRawMessage(`{"markup": %s}`, helpers.Base64BlankMarkup),
	})

	assert.Nil(suite.T(), err, "events.TrackDOMFromJSON() should succeed")
	eventJSON, err := trackDOM.ToJSON()
	assert.Nil(suite.T(), err, "trackDOM.ToJSON() should succeed")

	suite.mockedConnection.
		Command("HSET", suite.recordID, "1487696788863", eventJSON).
		ExpectError(errors.New("Redis error"))

	err = suite.exporter.Export(trackDOM, suite.recordID)
	assert.NotNil(suite.T(), err, "RedisExporter#Export should fail with an error")
}

func TestRedisExporter(t *testing.T) {
	suite.Run(t, new(RedisExporterTestSuite))
}
