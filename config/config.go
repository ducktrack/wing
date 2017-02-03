package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"github.com/pkg/errors"
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
		return nil, errors.Errorf("The configuration file is missing, expected file '%s'", path)
	}

	ymlBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal([]byte(ymlBytes), &c)
	return c, err
}
