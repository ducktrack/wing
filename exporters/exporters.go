package exporters

import (
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/events"
	"github.com/pkg/errors"
)

// Exporter interface
type Exporter interface {
	Initialize() error
	Export(trackable events.Trackable, recordID string) error
	Stop() error
}

// Lookup returns an exporter
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
