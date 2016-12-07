package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	Exporter     string
	FileExporter FileExporter `yaml:"file_exporter,omitempty"`
}

type FileExporter struct {
	Folder string
}

func ReadConfigFile(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		errorMessage := fmt.Sprintf("The configuration file is missing, expected file \"%s\"", path)
		return nil, errors.New(errorMessage)
	}

	ymlBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal([]byte(ymlBytes), &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
