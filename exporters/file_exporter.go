package exporters

import (
	"fmt"
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/trackentry"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// FileExporter definition
type FileExporter struct {
	Config config.FileExporter
}

// Export writes a file (<createdAt>.html) with the markup
func (fe *FileExporter) Export(trackEntry *trackentry.TrackEntry, recordID string) error {
	markup, err := trackEntry.Rinse()
	if err != nil {
		return errors.Wrap(err, "Failed to rinse the markup")
	}

	recordPath := filepath.Join(fe.Config.Folder, recordID)
	os.MkdirAll(recordPath, os.ModePerm)

	fileName := filepath.Join(recordPath, fmt.Sprintf("%d.html", trackEntry.CreatedAt))
	err = ioutil.WriteFile(fileName, []byte(markup), 0644)
	return errors.Wrapf(err, "Failed to save track entry to '%s'", fileName)
}

// Stop doesn't do anything for FileExporter
func (fe *FileExporter) Stop() error {
	return nil
}
