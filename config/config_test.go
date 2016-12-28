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
	content := `exporter: file

file_exporter:
  folder: /tmp/track_entries`
	tempFile, err := writeTempFile(content)
	assert.Nil(t, err)
	defer os.Remove(tempFile)

	c, err := ReadConfigFile(tempFile)
	assert.Nil(t, err)

	assert.Equal(t, "file", c.Exporter, "should have correct exporter")
	assert.NotNil(t, c.FileExporter, "should have a file exporter")
	assert.Equal(t, "/tmp/track_entries", c.FileExporter.Folder, "should have correct folder")
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
