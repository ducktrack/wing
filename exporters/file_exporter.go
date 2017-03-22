package exporters

import (
	"fmt"
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/events"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
)

const filePermission = 0644

// FileExporter definition
type FileExporter struct {
	Config config.FileExporter
	Fs     afero.Fs
}

// NewFileExporter is the construtor of FileExporter
func NewFileExporter(config config.FileExporter) *FileExporter {
	return &FileExporter{
		Config: config,
		Fs:     afero.NewOsFs(),
	}
}

// Initialize checks if the application has write permissions
func (fe *FileExporter) Initialize() error {
	rootFolder := fe.Config.Folder
	if len(rootFolder) == 0 || rootFolder == "/" {
		return errors.Errorf("Invalid folder '%s'", rootFolder)
	}

	err := fe.Fs.MkdirAll(rootFolder, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "Failed to prepare root folder")
	}

	checkFile := filepath.Join(rootFolder, ".check")
	err = afero.WriteFile(fe.Fs, checkFile, []byte(`check`), filePermission)
	if err != nil {
		return errors.Wrapf(err, "Failed to write to '%s'", rootFolder)
	}
	defer fe.Fs.Remove(checkFile)

	return nil
}

// Export writes a file (<rootFolder>/<recordID>/<createdAt>.json) with the event
func (fe *FileExporter) Export(trackable events.Trackable, recordID string) error {
	event := trackable.GetEvent()
	json, err := trackable.ToJSON()
	if err != nil {
		return errors.Wrap(err, "Failed to encode JSON")
	}

	recordPath := filepath.Join(fe.Config.Folder, recordID)
	fe.Fs.MkdirAll(recordPath, os.ModePerm)

	fileName := filepath.Join(recordPath, fmt.Sprintf("%d.json", event.CreatedAt))
	err = afero.WriteFile(fe.Fs, fileName, []byte(json), filePermission)
	return errors.Wrapf(err, "Failed to export event to '%s'", fileName)
}

// Stop doesn't do anything for FileExporter
func (fe *FileExporter) Stop() error {
	return nil
}
