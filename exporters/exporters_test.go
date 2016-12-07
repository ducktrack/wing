package exporters

import (
	"github.com/ducktrack/wing/config"
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
	if err != nil {
		t.Errorf("Lookup failed with:\n%v\n", err.Error())
	}

	interfaceType := reflect.TypeOf(exporter).String()
	if interfaceType != "*exporters.FileExporter" {
		t.Errorf("Lookup returned a different type than expected:\n%v\n", interfaceType)
	}
}

func TestLookupWhenExporterIsInvalid(t *testing.T) {
	appConfig := config.Config{
		Exporter: "invalid",
	}

	expected := "No exporter found for \"invalid\""
	_, err := Lookup(&appConfig)
	if err == nil || err.Error() != expected {
		t.Errorf(
			"Lookup failed with a different error:\ngot:\n%v\nwant:\n%v\n",
			err.Error(),
			expected,
		)
	}
}
