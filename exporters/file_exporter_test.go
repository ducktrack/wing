package exporters

import (
	"fmt"
	"github.com/ducktrack/wing/config"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"testing"
)

func TestExport(t *testing.T) {
	exporterConfig := config.FileExporter{
		Folder: "/tmp/test/track_entries",
	}

	trackEntry := &TrackEntry{
		CreatedAt: 123456,
		Markup:    "PGh0bWw+PC9odG1sPg==",
	}

	recordId := uuid.NewV4().String()

	exporter := FileExporter{Config: exporterConfig}
	err := exporter.Export(trackEntry, recordId)

	if err != nil {
		t.Errorf("FileExporter#Export failed with:\n%v\n", err.Error())
	}

	recordPath := fmt.Sprintf("/tmp/test/track_entries/%s/%d.html", recordId, trackEntry.CreatedAt)
	if _, err := os.Stat(recordPath); os.IsNotExist(err) {
		t.Errorf("FileExporter#Export failed to save track entry, expected:\n%v\n to exist", recordPath)
		return
	}

	htmlBytes, _ := ioutil.ReadFile(recordPath)
	expected := "<html></html>"
	if string(htmlBytes) != expected {
		t.Errorf(
			"FileExporter#Export saved the wrong content:\ngot:\n%v\nwant:\n%v\n",
			string(htmlBytes),
			expected,
		)
	}

	_ = os.Remove(recordPath)
}
