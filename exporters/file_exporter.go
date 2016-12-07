package exporters

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/ducktrack/wing/config"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileExporter struct {
	Config config.FileExporter
}

func (fe *FileExporter) Export(trackEntry *TrackEntry, recordId string) error {
	htmlBytes, err := base64.StdEncoding.DecodeString(trackEntry.Markup)
	if err != nil {
		return errors.New("Invalid base64 payload")
	}

	recordPath := filepath.Join(fe.Config.Folder, recordId)
	os.MkdirAll(recordPath, os.ModePerm)

	fileName := filepath.Join(recordPath, fmt.Sprintf("%d.html", trackEntry.CreatedAt))
	err = ioutil.WriteFile(fileName, htmlBytes, 0644)
	if err != nil {
		return errors.New("Fail to save track entry")
	}

	return nil
}
