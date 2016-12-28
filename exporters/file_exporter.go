package exporters

import (
	"fmt"
	"github.com/duckclick/wing/config"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/duckclick/wing/trackentry"
	"github.com/pkg/errors"
)

type FileExporter struct {
	Config config.FileExporter
}

func (fe *FileExporter) Export(trackEntry *trackentry.TrackEntry, recordId string) error {
	markup, err := trackEntry.Rinse()
	if err != nil {
		return errors.Wrap(err, "Failed to rinse the markup")
	}

	recordPath := filepath.Join(fe.Config.Folder, recordId)
	os.MkdirAll(recordPath, os.ModePerm)

	fileName := filepath.Join(recordPath, fmt.Sprintf("%d.html", trackEntry.CreatedAt))
	err = ioutil.WriteFile(fileName, []byte(markup), 0644)
	return errors.Wrapf(err, "Failed to save track entry to '%s'", fileName)
}

func (fe *FileExporter) Stop() error {
	return nil
}