package main

import (
	"fmt"
	"github.com/ducktrack/wing/config"
	"github.com/ducktrack/wing/handlers"
	"net/http"
	"os"
)

const DEFAULT_CONFIG_FILE = "application.yml"

func main() {
	port := getPort()
	configFilePath := getConfigFilePath()

	config, err := config.ReadConfigFile(configFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Config file: %s\n", configFilePath)
	fmt.Printf("Using exporter: %s\n", config.Exporter)
	fmt.Printf("Starting Wing at port %s\n", port)

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
