package exporters

import (
	"fmt"
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/events"
	helpers "github.com/duckclick/wing/testing"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestExport(t *testing.T) {
	exporterConfig := config.FileExporter{
		Folder: "/tmp/test/track_entries",
	}

	trackDOM, err := events.TrackDOMFromJSON(events.Event{
		CreatedAt:  1487696788863,
		URL:        "http://example.org/some/path",
		RawPayload: helpers.CreateRawMessage(`{"markup": %s}`, helpers.Base64BlankMarkup),
	})

	assert.Nil(t, err, "events.TrackDOMFromJSON() should succeed")
	recordID := uuid.NewV4().String()

	exporter := FileExporter{Config: exporterConfig}
	err = exporter.Export(trackDOM, recordID)
	assert.Nil(t, err, "FileExporter#Export should succeed")

	recordPath := fmt.Sprintf("/tmp/test/track_entries/%s/%d.json", recordID, trackDOM.CreatedAt)
	if _, err := os.Stat(recordPath); os.IsNotExist(err) {
		fmt.Printf("FileExporter#Export failed to save track entry, expected:\n%v\n to exist", recordPath)
		t.FailNow()
	}
	defer os.Remove(recordPath)

	htmlBytes, _ := ioutil.ReadFile(recordPath)
	assert.Equal(
		t,
		`{"created_at":1487696788863,"url":"http://example.org/some/path","type":"TrackDOM","payload":{"markup":"<html><head></head><body></body></html>"}}`,
		strings.TrimSpace(string(htmlBytes)),
		"FileExporter#Export should save the expected content",
	)
}
