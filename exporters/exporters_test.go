package exporters

import (
	"github.com/ducktrack/wing/config"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestLookupWhenExporterIsFile(t *testing.T) {
	appConfig := config.Config{
		Exporter: "file",
		FileExporter: config.FileExporter{
			Folder: "/tmp/test/track_entries",
		},
	}

	exporter, err := Lookup(&appConfig)
	assert.Nil(t, err)

	interfaceType := reflect.TypeOf(exporter).String()
	assert.Equal(t, "*exporters.FileExporter", interfaceType, "should use interface type from the configuration")
}

func TestLookupWhenExporterIsInvalid(t *testing.T) {
	appConfig := config.Config{
		Exporter: "invalid",
	}

	_, err := Lookup(&appConfig)
	assert.NotNil(t, err)

	assert.Equal(t, "No exporter found for 'invalid'", err.Error(), "should fail with not found exporter")
}
