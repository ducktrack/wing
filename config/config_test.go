package config

import (
	"testing"
)

func TestWhenFileIsMissing(t *testing.T) {
	_, err := ReadConfigFile("/tmp/missing-wing-application.yml")
	expected := "The configuration file is missing, expected file \"/tmp/missing-wing-application.yml\""
	if err == nil || err.Error() != expected {
		t.Errorf(
			"ReadConfigFile returned unexpected error\ngot:\n%v\n\nwant:\n%v",
			err.Error(),
			expected,
		)
	}
}
