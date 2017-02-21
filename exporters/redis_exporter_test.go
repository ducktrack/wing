package exporters

import (
	"errors"
	"github.com/duckclick/wing/config"
	helpers "github.com/duckclick/wing/testing"
	"github.com/duckclick/wing/trackentry"
	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var trackEntry *trackentry.TrackEntry
var trackEntryJSON string
var recordID string
var exporter *RedisExporter
var mockedConnection *redigomock.Conn

func TestMain(m *testing.M) {
	recordID = uuid.NewV4().String()
	trackEntry = &trackentry.TrackEntry{
		CreatedAt: 1487696788863,
		URL:       "http://example.org/some/path",
		Markup:    helpers.Base64BlankMarkup,
	}

	trackEntryJSON, _ = trackEntry.ToJSON()
	exporterConfig := config.RedisExporter{
		Host: "foo",
		Port: 1234,
	}
	exporter = &RedisExporter{config: exporterConfig}
	mockedConnection = redigomock.NewConn()
	exporter.pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return mockedConnection, nil
		},
	}
	defer exporter.Stop()

	os.Exit(m.Run())
}

func TestRedisExport(t *testing.T) {
	mockedConnection.
		Command("HSET", recordID, "1487696788863", trackEntryJSON).
		Expect(nil)

	err := exporter.Export(trackEntry, recordID)
	assert.Nil(t, err, "RedisExporter#Export should succeed")
}

func TestExportReturnsErrorOnRedisError(t *testing.T) {
	mockedConnection.
		Command("HSET", recordID, "1487696788863", trackEntryJSON).
		ExpectError(errors.New("Redis error"))

	err := exporter.Export(trackEntry, recordID)
	assert.NotNil(t, err, "RedisExporter#Export should fail with an error")
}
