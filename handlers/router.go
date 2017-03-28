package handlers

import (
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/exporters"
	"github.com/garyburd/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	jose "gopkg.in/square/go-jose.v1"
	"io/ioutil"
	"os"
	"time"
)

// AppContext definition
type AppContext struct {
	Config        *config.Config
	Exporter      exporters.Exporter
	Redis         *redis.Pool
	JWEPublicKey  interface{}
	JWEPrivateKey interface{}
}

// Router definition
type Router struct {
	*httprouter.Router
	*AppContext
}

// Route definition
type Route func(*AppContext) httprouter.Handle

// NewRouter creates a new router with the app context
func NewRouter(wingConfig *config.Config, exporter exporters.Exporter) (*Router, error) {
	privateKeyData, err := readFile(wingConfig.JWEPrivateKeyFile)
	if err != nil {
		return nil, err
	}

	publicKeyData, err := readFile(wingConfig.JWEPublicKeyFile)
	if err != nil {
		return nil, err
	}

	privateKey, err := jose.LoadPrivateKey([]byte(privateKeyData))
	if err != nil {
		return nil, err
	}

	publicKey, err := jose.LoadPublicKey([]byte(publicKeyData))
	if err != nil {
		return nil, err
	}

	return &Router{
		httprouter.New(),
		&AppContext{
			Config:        wingConfig,
			Exporter:      exporter,
			Redis:         createRedisConnectionPool("localhost:6379"),
			JWEPrivateKey: privateKey,
			JWEPublicKey:  publicKey,
		},
	}, nil
}

// GET draw a get handler
func (r *Router) GET(path string, route Route) {
	r.Handle("GET", path, route(r.AppContext))
}

// POST draw a post handler
func (r *Router) POST(path string, route Route) {
	r.Handle("POST", path, route(r.AppContext))
}

func createRedisConnectionPool(connString string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", connString)
		},
	}
}

func readFile(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", errors.Wrapf(err, "File '%s' is missing, expected file '%s'", path)
	}

	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to read file '%s'", path)
	}

	return string(fileBytes), nil
}
