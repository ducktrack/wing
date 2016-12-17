package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/handlers"
	"net/http"
	"os"
)

const DEFAULT_CONFIG_FILE = "application.yml"

func main() {
	port := getPort()
	configFilePath := getConfigFilePath()
	log.Infof("Config file: %s", configFilePath)

	config, err := config.ReadConfigFile(configFilePath)
	if err != nil {
		log.WithError(err).Fatal("Failed to read config file")
		os.Exit(1)
	}

	log.Infof("Using exporter: %s", config.Exporter)
	log.Infof("Starting Wing at port %s", port)

	http.Handle("/", &handlers.TrackEntryHandler{Config: config})

	host := fmt.Sprintf(":%s", port)
	http.ListenAndServe(host, nil)
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
