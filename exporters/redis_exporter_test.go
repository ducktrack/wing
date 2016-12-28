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
)

func TestRedisExport(t *testing.T) {
	exporterConfig := config.RedisExporter{
		Host: "foo",
		Port: 1234,
	}

	htmlSample := "<html><head></head><body></body></html>"
	htmlAsBase64 := base64.StdEncoding.EncodeToString([]byte(htmlSample))

	trackEntry := &TrackEntry{
		CreatedAt: 123456,
		Markup:    htmlAsBase64,
	}

	recordId := uuid.NewV4().String()

	exporter := redisExporter{config: exporterConfig}
	mockedConnection := redigomock.NewConn()
	exporter.pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return mockedConnection, nil
		},
	}
	defer exporter.Stop()

	mockedConnection.Command("HSET", recordId, "123456", htmlSample).Expect(nil)

	err := exporter.Export(trackEntry, recordId)
	assert.Nil(t, err, "export should succeed")
}

func TestExportReturnsErrorOnRedisError(t *testing.T) {
	exporterConfig := config.RedisExporter{
		Host: "foo",
		Port: 1234,
	}

	htmlSample := "<html><head></head><body></body></html>"
	htmlAsBase64 := base64.StdEncoding.EncodeToString([]byte(htmlSample))

	trackEntry := &TrackEntry{
		CreatedAt: 123456,
		Markup:    htmlAsBase64,
	}

	recordId := uuid.NewV4().String()

	exporter := redisExporter{config: exporterConfig}
	mockedConnection := redigomock.NewConn()
	exporter.pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return mockedConnection, nil
		},
	}
	defer exporter.Stop()

	mockedConnection.Command("HSET", recordId, "123456", htmlSample).ExpectError(errors.New("Redis error"))

	err := exporter.Export(trackEntry, recordId)
	assert.NotNil(t, err, "export should fail with an error")
}
