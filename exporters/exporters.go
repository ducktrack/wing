package exporters

import (
	"errors"
	"fmt"
	"github.com/duckclick/wing/config"
)

type TrackEntry struct {
	CreatedAt int    `json:"created_at"`
	Origin    string `json:"origin"`
	Markup    string `json:"markup"`
}

type Exporter interface {
	Export(trackEntry *TrackEntry, recordId string) error
}

func Lookup(config *config.Config) (Exporter, error) {
	switch config.Exporter {
	case "file":
		return &FileExporter{Config: config.FileExporter}, nil
	}

	errorMessage := fmt.Sprintf("No exporter found for '%s'", config.Exporter)
	return nil, errors.New(errorMessage)
}
