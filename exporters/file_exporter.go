package exporters

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/duckclick/wing/config"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FileExporter struct {
	Config config.FileExporter
}

func (fe *FileExporter) Export(trackEntry *TrackEntry, recordId string) error {
	htmlBytes, err := base64.StdEncoding.DecodeString(trackEntry.Markup)
	if err != nil {
		return errors.New("Invalid base64 payload")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(htmlBytes)))
	if err != nil {
		return errors.New("Failed to parse HTML")
	}

	scripts := doc.Find("script")
	scripts.Each(func(i int, s *goquery.Selection) { s.Remove() })

	content, err := doc.Html()
	if err != nil {
		return errors.New("Failed to generate secure HTML")
	}

	recordPath := filepath.Join(fe.Config.Folder, recordId)
	os.MkdirAll(recordPath, os.ModePerm)

	fileName := filepath.Join(recordPath, fmt.Sprintf("%d.html", trackEntry.CreatedAt))
	err = ioutil.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to save track entry to '%s'", fileName))
	}

	return nil
}
