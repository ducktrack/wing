package exporters

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"encoding/base64"
	"github.com/duckclick/wing/config"
	"github.com/satori/go.uuid"
	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"errors"
	"github.com/duckclick/wing/trackentry"
	"os"
)

var htmlSample string
var trackEntry *trackentry.TrackEntry
var recordId string
var exporter redisExporter
var mockedConnection *redigomock.Conn

func TestMain(m *testing.M) {
	htmlSample = "<html><head></head><body></body></html>"
	trackEntry = &trackentry.TrackEntry{
		CreatedAt: 123456,
		Markup:    base64.StdEncoding.EncodeToString([]byte(htmlSample)),
	}

	recordId = uuid.NewV4().String()

	exporterConfig := config.RedisExporter{
		Host: "foo",
		Port: 1234,
	}
	exporter = redisExporter{config: exporterConfig}
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
	mockedConnection.Command("HSET", recordId, "123456", htmlSample).Expect(nil)

	err := exporter.Export(trackEntry, recordId)
	assert.Nil(t, err, "export should succeed")
}

func TestExportReturnsErrorOnRedisError(t *testing.T) {
	mockedConnection.Command("HSET", recordId, "123456", htmlSample).ExpectError(errors.New("Redis error"))

	err := exporter.Export(trackEntry, recordId)
	assert.NotNil(t, err, "export should fail with an error")
}
