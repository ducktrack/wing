package testing

import (
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/exporters"
	"github.com/duckclick/wing/handlers"
)

// CreateFileExporterAppContext creates a new *handlers.AppContext configured with FileExporter
func CreateFileExporterAppContext() *handlers.AppContext {
	appConfig := config.Config{
		RecordIdCookieName: "record_id",
		RecordIdCookieExpiration: 120,
		Exporter: "file",
		FileExporter: config.FileExporter{
			Folder: "/tmp/test/track_entries",
		},
	}
	fileExporter, _ := exporters.Lookup(&appConfig)
	return &handlers.AppContext{
		Config:   &appConfig,
		Exporter: fileExporter,
	}
}

func CreateRouter() *handlers.Router {
	appContext := CreateFileExporterAppContext()
	return handlers.NewRouter(appContext.Config, appContext.Exporter)
}
