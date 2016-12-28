package exporters

import (
	"errors"
	"fmt"
	"github.com/duckclick/wing/config"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/duckclick/wing/trackentry"
)

type FileExporter struct {
	Config config.FileExporter
}

func (fe *FileExporter) Export(trackEntry *trackentry.TrackEntry, recordId string) error {
	markup, err := trackEntry.Rinse()
	if err != nil {
		return errors.New("Failed to rinse the markup")
	}

	recordPath := filepath.Join(fe.Config.Folder, recordId)
	os.MkdirAll(recordPath, os.ModePerm)

	fileName := filepath.Join(recordPath, fmt.Sprintf("%d.html", trackEntry.CreatedAt))
	err = ioutil.WriteFile(fileName, []byte(markup), 0644)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to save track entry to '%s'", fileName))
	}

	return nil
}

func (fe *FileExporter) Stop() error {
	return nil
}