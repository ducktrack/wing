package testing

import (
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/exporters"
	"github.com/duckclick/wing/handlers"
)

// CreateBasicConfig definition
func CreateBasicConfig() config.Config {
	return config.Config{
		JWEPrivateKeyFile: "../jwe_certs/privatekey.pem",
		JWEPublicKeyFile:  "../jwe_certs/publickey.pem",
	}
}

// CreateFileExporterAppContext creates a new *handlers.AppContext configured with FileExporter
func CreateFileExporterAppContext() *handlers.AppContext {
	appConfig := CreateBasicConfig()
	appConfig.Exporter = "file"
	appConfig.FileExporter = config.FileExporter{
		Folder: "/tmp/test/track_entries",
	}

	fileExporter, _ := exporters.Lookup(&appConfig)
	appContext, _ := handlers.NewAppContext(&appConfig, fileExporter)
	return appContext
}

// CreateRouter definition
func CreateRouter() (*handlers.Router, error) {
	appContext := CreateFileExporterAppContext()
	return handlers.NewRouter(appContext.Config, appContext.Exporter)
}
