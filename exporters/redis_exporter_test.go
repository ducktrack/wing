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

var htmlSample string
var trackEntry *trackentry.TrackEntry
var recordID string
var exporter *RedisExporter
var mockedConnection *redigomock.Conn

func TestMain(m *testing.M) {
	htmlSample = "<html><head></head><body></body></html>"
	trackEntry = &trackentry.TrackEntry{
		CreatedAt: 123456,
		Markup:    helpers.ToBase64(htmlSample),
	}

	recordID = uuid.NewV4().String()

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
	mockedConnection.Command("HSET", recordID, "123456", htmlSample).Expect(nil)

	err := exporter.Export(trackEntry, recordID)
	assert.Nil(t, err, "export should succeed")
}

func TestExportReturnsErrorOnRedisError(t *testing.T) {
	mockedConnection.Command("HSET", recordID, "123456", htmlSample).ExpectError(errors.New("Redis error"))

	err := exporter.Export(trackEntry, recordID)
	assert.NotNil(t, err, "export should fail with an error")
}
