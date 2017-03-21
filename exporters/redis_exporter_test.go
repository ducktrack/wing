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
	"os"
	"testing"
)

var event events.Trackable
var eventJSON string
var recordID string
var exporter *exporters.RedisExporter
var mockedConnection *redigomock.Conn

func TestMain(m *testing.M) {
	recordID = uuid.NewV4().String()

	exporterConfig := config.RedisExporter{
		Host: "foo",
		Port: 1234,
	}
	exporter = &exporters.RedisExporter{Config: exporterConfig}
	mockedConnection = redigomock.NewConn()
	exporter.Pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return mockedConnection, nil
		},
	}
	defer exporter.Stop()

	os.Exit(m.Run())
}

func TestRedisExport(t *testing.T) {
	trackDOM, err := events.TrackDOMFromJSON(events.Event{
		CreatedAt:  1487696788863,
		URL:        "http://example.org/some/path",
		RawPayload: helpers.CreateRawMessage(`{"markup": %s}`, helpers.Base64BlankMarkup),
	})

	assert.Nil(t, err, "events.TrackDOMFromJSON() should succeed")

	eventJSON, err = trackDOM.ToJSON()
	assert.Nil(t, err, "trackDOM.ToJSON() should succeed")

	mockedConnection.
		Command("HSET", recordID, "1487696788863", eventJSON).
		Expect(nil)

	err = exporter.Export(trackDOM, recordID)
	assert.Nil(t, err, "RedisExporter#Export should succeed")
}

func TestExportReturnsErrorOnRedisError(t *testing.T) {
	trackDOM, err := events.TrackDOMFromJSON(events.Event{
		CreatedAt:  1487696788863,
		URL:        "http://example.org/some/path",
		RawPayload: helpers.CreateRawMessage(`{"markup": %s}`, helpers.Base64BlankMarkup),
	})

	assert.Nil(t, err, "events.TrackDOMFromJSON() should succeed")

	eventJSON, err = trackDOM.ToJSON()
	assert.Nil(t, err, "trackDOM.ToJSON() should succeed")

	mockedConnection.
		Command("HSET", recordID, "1487696788863", eventJSON).
		ExpectError(errors.New("Redis error"))

	err = exporter.Export(trackDOM, recordID)
	assert.NotNil(t, err, "RedisExporter#Export should fail with an error")
}
