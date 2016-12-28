package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/handlers"
	"net/http"
	"os"
	"github.com/duckclick/wing/exporters"
)

const DEFAULT_CONFIG_FILE = "application.yml"

func main() {
	port := getPort()
	configFilePath := getConfigFilePath()
	log.Infof("Config file: %s", configFilePath)

	wingConfig, err := config.ReadConfigFile(configFilePath)
	if err != nil {
		log.WithError(err).Fatal("Failed to read config file")
		os.Exit(1)
	}

	exporter, err := exporters.Lookup(wingConfig)
	if err != nil {
		log.WithError(err).Fatal("Failed to instantiate an exporter")
		os.Exit(1)
	}
	defer exporter.Stop()

	log.Infof("Using exporter: %s", wingConfig.Exporter)
	log.Infof("Starting Wing at port %s", port)

	http.Handle("/", &handlers.TrackEntryHandler{Config: *wingConfig, Exporter: exporter})
	host := fmt.Sprintf(":%s", port)

	if wingConfig.TLSCertFile != "" && wingConfig.TLSKeyFile != "" {
		log.Infof("Using TLS")
		http.ListenAndServeTLS(host, wingConfig.TLSCertFile, wingConfig.TLSKeyFile, nil)

	} else {
		http.ListenAndServe(host, nil)
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "7273"
	}

	return port
}

func getConfigFilePath() string {
	configFilePath := os.Getenv("CONFIG")
	if len(configFilePath) == 0 {
		configFilePath = DEFAULT_CONFIG_FILE
	}

	return configFilePath
}
