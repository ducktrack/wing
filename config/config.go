package config

import (
	"github.com/pkg/errors"
	jose "gopkg.in/square/go-jose.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// Config definition
type Config struct {
	TLSCertFile        string        `yaml:"tls_cert_file"`
	TLSKeyFile         string        `yaml:"tls_key_file"`
	JWEPrivateKeyFile  string        `yaml:"jwe_private_key_file"`
	JWEPublicKeyFile   string        `yaml:"jwe_public_key_file"`
	SessionTokenSecret string        `yaml:"session_token_secret"`
	Exporter           string        `yaml:"exporter"`
	FileExporter       FileExporter  `yaml:"file_exporter,omitempty"`
	RedisExporter      RedisExporter `yaml:"redis_exporter,omitempty"`
}

// FileExporter definition
type FileExporter struct {
	Folder string `yaml:"folder"`
}

// RedisExporter definition
type RedisExporter struct {
	Host string `yaml:"host"`
	Port int16  `yaml:"port"`
}

// PrivateKey used for go-jose.v1 private keys
type PrivateKey interface{}

// PublicKey used for go-jose.v1 public keys
type PublicKey interface{}

// ReadConfigFile reads and parses the configuration file (application.yml)
func ReadConfigFile(path string) (*Config, error) {
	filePayload, err := readFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "The configuration file is missing, expected file '%s'", path)
	}

	c := &Config{}
	err = yaml.Unmarshal([]byte(filePayload), &c)
	return c, err
}

// LoadJWEKeys returns the private and the public keys parsed by go-jose.v1
func LoadJWEKeys(wingConfig *Config) (PrivateKey, PublicKey, error) {
	privateKeyData, err := readFile(wingConfig.JWEPrivateKeyFile)
	if err != nil {
		return nil, nil, err
	}

	privateKey, err := jose.LoadPrivateKey([]byte(privateKeyData))

	publicKeyData, err := readFile(wingConfig.JWEPublicKeyFile)
	if err != nil {
		return nil, nil, err
	}

	publicKey, err := jose.LoadPublicKey([]byte(publicKeyData))
	return privateKey.(PrivateKey), publicKey.(PublicKey), nil
}

func readFile(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", errors.Wrapf(err, "File '%s' is missing", path)
	}

	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to read file '%s'", path)
	}

	return string(fileBytes), nil
}
