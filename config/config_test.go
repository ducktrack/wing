package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWhenFileIsMissing(t *testing.T) {
	_, err := ReadConfigFile("/tmp/missing-wing-application.yml")
	assert.NotNil(t, err)

	expected := "The configuration file is missing, expected file '/tmp/missing-wing-application.yml'"
	assert.Equal(t, expected, err.Error(), "should complain about missing configuration file")
}
