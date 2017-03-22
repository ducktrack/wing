package exporters

import (
	"fmt"
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/events"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// FileExporter definition
type FileExporter struct {
	Config config.FileExporter
}

// Initialize checks if the application has write permissions
func (fe *FileExporter) Initialize() error {
	return nil
}

// Export writes a file (<createdAt>.html) with the markup
func (fe *FileExporter) Export(trackable events.Trackable, recordID string) error {
	event := trackable.GetEvent()
	json, err := trackable.ToJSON()
	if err != nil {
		return errors.Wrap(err, "Failed to encode JSON")
	}

	recordPath := filepath.Join(fe.Config.Folder, recordID)
	os.MkdirAll(recordPath, os.ModePerm)

	fileName := filepath.Join(recordPath, fmt.Sprintf("%d.json", event.CreatedAt))
	err = ioutil.WriteFile(fileName, []byte(json), 0644)
	return errors.Wrapf(err, "Failed to save track entry to '%s'", fileName)
}

// Stop doesn't do anything for FileExporter
func (fe *FileExporter) Stop() error {
	return nil
}
