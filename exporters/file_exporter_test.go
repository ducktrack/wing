package exporters_test

import (
	"fmt"
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/events"
	"github.com/duckclick/wing/exporters"
	helpers "github.com/duckclick/wing/testing"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type FileExplorerTestSuite struct {
	suite.Suite
	recordID string
	exporter exporters.FileExporter
}

func (suite *FileExplorerTestSuite) SetupTest() {
	suite.recordID = uuid.NewV4().String()
	suite.exporter = exporters.FileExporter{
		Config: config.FileExporter{
			Folder: "/tmp/test/track_entries",
		},
	}
}

func (suite *FileExplorerTestSuite) TestExport() {
	trackDOM, err := events.TrackDOMFromJSON(events.Event{
		CreatedAt:  1487696788863,
		URL:        "http://example.org/some/path",
		RawPayload: helpers.CreateRawMessage(`{"markup": %s}`, helpers.Base64BlankMarkup),
	})

	assert.Nil(suite.T(), err, "events.TrackDOMFromJSON() should succeed")
	err = suite.exporter.Export(trackDOM, suite.recordID)
	assert.Nil(suite.T(), err, "FileExporter#Export should succeed")

	recordPath := fmt.Sprintf("/tmp/test/track_entries/%s/%d.json", suite.recordID, trackDOM.CreatedAt)
	if _, err := os.Stat(recordPath); os.IsNotExist(err) {
		fmt.Printf("FileExporter#Export failed to save track entry, expected:\n%v\n to exist", recordPath)
		suite.T().FailNow()
	}
	defer os.Remove(recordPath)

	htmlBytes, _ := ioutil.ReadFile(recordPath)
	assert.Equal(
		suite.T(),
		`{"created_at":1487696788863,"url":"http://example.org/some/path","type":"TrackDOM","payload":{"markup":"<html><head></head><body></body></html>"}}`,
		strings.TrimSpace(string(htmlBytes)),
		"FileExporter#Export should save the expected content",
	)
}

func TestFileExplorer(t *testing.T) {
	suite.Run(t, new(FileExplorerTestSuite))
}
