package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/exporters"
	"github.com/duckclick/wing/handlers"
	"github.com/rs/cors"
	"net/http"
	"os"
)

const defaultConfigFile = "application.yml"

func main() {
	port := getPort()
	wingConfig := readConfig()
	exporter := lookupExporter(wingConfig)

	log.Infof("Using exporter: %s", wingConfig.Exporter)

	err := exporter.Initialize()
	if err != nil {
		log.WithError(err).Fatal("Failed to Initialize exporter")
		os.Exit(1)
	}
	defer exporter.Stop()

	router, err := handlers.NewRouter(wingConfig, exporter)
	if err != nil {
		log.WithError(err).Fatal("Failed to create router")
		os.Exit(1)
	}

	router.DrawRoutes()
	mux := corsMiddleware(router)
	host := fmt.Sprintf(":%s", port)

	log.Infof("Starting Wing at port %s", port)

	if wingConfig.TLSCertFile != "" && wingConfig.TLSKeyFile != "" {
		log.Infof("Using TLS")
		http.ListenAndServeTLS(host, wingConfig.TLSCertFile, wingConfig.TLSKeyFile, mux)

	} else {
		http.ListenAndServe(host, mux)
	}
}

func corsMiddleware(router *handlers.Router) http.Handler {
	middleware := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"content-type"},
	})

	return middleware.Handler(router)
}

func lookupExporter(wingConfig *config.Config) exporters.Exporter {
	exporter, err := exporters.Lookup(wingConfig)
	if err != nil {
		log.WithError(err).Fatal("Failed to instantiate an exporter")
		os.Exit(1)
	}

	return exporter
}

func readConfig() *config.Config {
	configFilePath := getConfigFilePath()
	log.Infof("Config file: %s", configFilePath)

	wingConfig, err := config.ReadConfigFile(configFilePath)
	if err != nil {
		log.WithError(err).Fatal("Failed to read config file")
		os.Exit(1)
	}

	return wingConfig
}

func getConfigFilePath() string {
	configFilePath := os.Getenv("CONFIG")
	if len(configFilePath) == 0 {
		configFilePath = defaultConfigFile
	}

	return configFilePath
}

func getPort() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "7273"
	}

	return port
}
