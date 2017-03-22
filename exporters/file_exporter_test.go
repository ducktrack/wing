package exporters_test

import (
	"fmt"
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/events"
	"github.com/duckclick/wing/exporters"
	helpers "github.com/duckclick/wing/testing"
	"github.com/satori/go.uuid"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"strings"
	"testing"
)

type FileExporterTestSuite struct {
	suite.Suite
	Fs       afero.Fs
	recordID string
	exporter exporters.FileExporter
}

func (suite *FileExporterTestSuite) SetupTest() {
	suite.Fs = afero.NewMemMapFs()
	suite.recordID = uuid.NewV4().String()
	suite.exporter = exporters.FileExporter{
		Config: config.FileExporter{
			Folder: "/tmp/test/track_entries",
		},
		Fs: suite.Fs,
	}
}

func (suite *FileExporterTestSuite) TestInitialize() {
	err := suite.exporter.Initialize()
	rootFolder := suite.exporter.Config.Folder
	assert.Nil(suite.T(), err, "FileExporter#Initialize should succeed")

	rootDirExists, err := afero.DirExists(suite.Fs, rootFolder)
	assert.Nil(suite.T(), err, "afero.DirExists should succeed")
	assert.Equal(suite.T(), rootDirExists, true, "Root folder should be created")
}

func (suite *FileExporterTestSuite) TestInitializeWhenFSDoesNotHaveWriteAccess() {
	fs := afero.NewReadOnlyFs(suite.Fs)
	suite.exporter.Fs = fs
	err := suite.exporter.Initialize()
	assert.NotNil(suite.T(), err, "FileExporter#Initialize should fail with an error")
}

func (suite *FileExporterTestSuite) TestInitializeWhenRootFolderIsNotDefined() {
	suite.exporter.Config = config.FileExporter{}
	err := suite.exporter.Initialize()
	assert.NotNil(suite.T(), err, "FileExporter#Initialize should fail with an error")
}

func (suite *FileExporterTestSuite) TestExport() {
	trackDOM, err := events.TrackDOMFromJSON(events.Event{
		CreatedAt:  1487696788863,
		URL:        "http://example.org/some/path",
		RawPayload: helpers.CreateRawMessage(`{"markup": %s}`, helpers.Base64BlankMarkup),
	})

	assert.Nil(suite.T(), err, "events.TrackDOMFromJSON() should succeed")
	err = suite.exporter.Export(trackDOM, suite.recordID)
	assert.Nil(suite.T(), err, "FileExporter#Export should succeed")

	recordPath := fmt.Sprintf("/tmp/test/track_entries/%s/%d.json", suite.recordID, trackDOM.CreatedAt)
	if _, err := suite.Fs.Stat(recordPath); os.IsNotExist(err) {
		fmt.Printf("FileExporter#Export failed to save track entry, expected:\n%v\n to exist", recordPath)
		suite.T().FailNow()
	}
	defer suite.Fs.Remove(recordPath)

	htmlBytes, _ := afero.ReadFile(suite.Fs, recordPath)
	assert.Equal(
		suite.T(),
		`{"created_at":1487696788863,"url":"http://example.org/some/path","type":"TrackDOM","payload":{"markup":"<html><head></head><body></body></html>"}}`,
		strings.TrimSpace(string(htmlBytes)),
		"FileExporter#Export should save the expected content",
	)
}

func TestFileExporter(t *testing.T) {
	suite.Run(t, new(FileExporterTestSuite))
}
