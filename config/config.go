package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	TLSCertFile   string        `yaml:"tls_cert_file"`
	TLSKeyFile    string        `yaml:"tls_key_file"`
	Exporter      string
	FileExporter  FileExporter  `yaml:"file_exporter,omitempty"`
	RedisExporter RedisExporter `yaml:"redis_exporter,omitempty"`
}

type FileExporter struct {
	Folder string `yaml:"folder"`
}

type RedisExporter struct {
	Host     string `yaml:"host"`
	Port     int16  `yaml:"port"`
}

func ReadConfigFile(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		errorMessage := fmt.Sprintf("The configuration file is missing, expected file '%s'", path)
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
