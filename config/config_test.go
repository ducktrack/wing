package config_test

import (
	"github.com/duckclick/wing/config"
	helpers "github.com/duckclick/wing/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (suite *ConfigTestSuite) TestReadConfigFile() {
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

func (suite *ConfigTestSuite) TestReadConfigFileWhenFileIsMissing() {
	_, err := config.ReadConfigFile("/tmp/missing-wing-application.yml")
	assert.NotNil(suite.T(), err)

	expected := "The configuration file is missing, expected file '/tmp/missing-wing-application.yml'"
	assert.True(suite.T(), strings.Contains(err.Error(), expected), "should complain about missing configuration file")
}

func (suite *ConfigTestSuite) TestLoadJWEKeys() {
	wingConfig := helpers.CreateBasicConfig()
	privateKey, publicKey, err := config.LoadJWEKeys(&wingConfig)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), privateKey)
	assert.NotNil(suite.T(), publicKey)
}

func (suite *ConfigTestSuite) TestLoadJWEKeysWhenPrivateKeyDoesNotExist() {
	wingConfig := helpers.CreateBasicConfig()
	wingConfig.JWEPrivateKeyFile = "wrong"
	privateKey, publicKey, err := config.LoadJWEKeys(&wingConfig)
	assert.NotNil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(err.Error(), "File 'wrong' is missing"))
	assert.Nil(suite.T(), privateKey)
	assert.Nil(suite.T(), publicKey)
}

func (suite *ConfigTestSuite) TestLoadJWEKeysWhenPublicKeyDoesNotExist() {
	wingConfig := helpers.CreateBasicConfig()
	wingConfig.JWEPublicKeyFile = "wrong"
	privateKey, publicKey, err := config.LoadJWEKeys(&wingConfig)
	assert.NotNil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(err.Error(), "File 'wrong' is missing"))
	assert.Nil(suite.T(), privateKey)
	assert.Nil(suite.T(), publicKey)
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
