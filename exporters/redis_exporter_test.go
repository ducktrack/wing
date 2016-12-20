package exporters

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"encoding/base64"
	"github.com/duckclick/wing/config"
	"github.com/satori/go.uuid"
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

	exporter := redisExporter{Config: exporterConfig}
	err := exporter.Export(trackEntry, recordId)
	assert.Nil(t, err, "export should succeed")
}
