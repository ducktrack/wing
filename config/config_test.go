package config_test

import (
	"github.com/duckclick/wing/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"testing"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (suite *ConfigTestSuite) TestWhenFileIsMissing() {
	_, err := config.ReadConfigFile("/tmp/missing-wing-application.yml")
	assert.NotNil(suite.T(), err)

	expected := "The configuration file is missing, expected file '/tmp/missing-wing-application.yml'"
	assert.Equal(suite.T(), expected, err.Error(), "should complain about missing configuration file")
}

func (suite *ConfigTestSuite) TestWhenFileIsValid() {
	content := `exporter: file

file_exporter:
  folder: /tmp/track_entries`
	tempFile, err := writeTempFile(content)
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile)

	c, err := config.ReadConfigFile(tempFile)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), "file", c.Exporter, "should have correct exporter")
	assert.NotNil(suite.T(), c.FileExporter, "should have a file exporter")
	assert.Equal(suite.T(), "/tmp/track_entries", c.FileExporter.Folder, "should have correct folder")
}

func TestConfig(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func writeTempFile(content string) (fName string, err error) {
	byteContent := []byte(content)
	tempFile, err := ioutil.TempFile("", "example")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	fName = tempFile.Name()
	err = ioutil.WriteFile(fName, byteContent, 0644)
	if err != nil {
		return "", err
	}

	return fName, nil
}
