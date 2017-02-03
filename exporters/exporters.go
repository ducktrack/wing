package exporters

import (
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/trackentry"
	"github.com/pkg/errors"
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
	default:
		return nil, errors.Errorf("No exporter found for '%s'", config.Exporter)
	}
}
