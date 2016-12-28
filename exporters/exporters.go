package exporters

import (
	"errors"
	"fmt"
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/trackentry"
)

type Exporter interface {
	Export(trackEntry *trackentry.TrackEntry, recordId string) error
	Stop() error
}

func Lookup(config *config.Config) (Exporter, error) {
	switch config.Exporter {
	case "file":
		return &FileExporter{Config: config.FileExporter}, nil
	case "redis":
		return NewRedisExporter(config.RedisExporter), nil
	}

	errorMessage := fmt.Sprintf("No exporter found for '%s'", config.Exporter)
	return nil, errors.New(errorMessage)
}
