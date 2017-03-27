package exporters_test

import (
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/exporters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type ExportersTestSuite struct {
	suite.Suite
}

func (suite *ExportersTestSuite) TestLookupWhenExporterIsFile() {
	appConfig := config.Config{
		Exporter: "file",
		FileExporter: config.FileExporter{
			Folder: "/tmp/test/track_entries",
		},
	}

	exporter, err := exporters.Lookup(&appConfig)
	assert.Nil(suite.T(), err)

	interfaceType := reflect.TypeOf(exporter).String()
	assert.Equal(suite.T(), "*exporters.FileExporter", interfaceType, "should use interface type from the configuration")
}

func (suite *ExportersTestSuite) TestLookupWhenExporterIsRedis() {
	appConfig := config.Config{
		Exporter:      "redis",
		RedisExporter: config.RedisExporter{},
	}

	exporter, err := exporters.Lookup(&appConfig)
	assert.Nil(suite.T(), err)

	interfaceType := reflect.TypeOf(exporter).String()
	assert.Equal(suite.T(), "*exporters.RedisExporter", interfaceType, "should use interface type from the configuration")
}

func (suite *ExportersTestSuite) TestLookupWhenExporterIsInvalid() {
	appConfig := config.Config{
		Exporter: "invalid",
	}

	_, err := exporters.Lookup(&appConfig)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "No exporter found for 'invalid'", err.Error(), "should fail with not found exporter")
}

func TestExporters(t *testing.T) {
	suite.Run(t, new(ExportersTestSuite))
}
