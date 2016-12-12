package config

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestWhenFileIsMissing(t *testing.T) {
	_, err := ReadConfigFile("/tmp/missing-wing-application.yml")
	assert.NotNil(t, err)

	expected := "The configuration file is missing, expected file '/tmp/missing-wing-application.yml'"
	assert.Equal(t, expected, err.Error(), "should complain about missing configuration file")
}

func TestWhenFileIsValid(t *testing.T) {
	configContent := `exporter: file

file_exporter:
  folder: /tmp/track_entries`
	tempFile, err := WriteTempFile(configContent)
	assert.Nil(t, err)
	defer os.Remove(tempFile)

	config, err := ReadConfigFile(tempFile)
	assert.Nil(t, err)

	assert.Equal(t, "file", config.Exporter, "should have correct exporter")
	assert.NotNil(t, config.FileExporter, "should have a file exporter")
	assert.Equal(t, "/tmp/track_entries", config.FileExporter.Folder, "should have correct folder")
}

func WriteTempFile(content string) (fName string, err error) {
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
