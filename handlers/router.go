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
	appContext *AppContext
}

// Route definition
type Route func(*AppContext) httprouter.Handle

// NewRouter creates a new router with the app context
func NewRouter(wingConfig *config.Config, exporter exporters.Exporter) (*Router, error) {
	privateKey, err := readPrivateKey(wingConfig.JWEPrivateKeyFile)
	if err != nil {
		return nil, err
	}

	publicKey, err := readPublicKey(wingConfig.JWEPublicKeyFile)
	if err != nil {
		return nil, err
	}

	return &Router{
		Router: httprouter.New(),
		appContext: &AppContext{
			Config:        wingConfig,
			Exporter:      exporter,
			JWEPrivateKey: privateKey,
			JWEPublicKey:  publicKey,
		},
	}, nil
}

// AppContext returns the app context
func (r *Router) AppContext() *AppContext {
	if r.appContext.Redis == nil {
		r.appContext.Redis = createRedisConnectionPool("localhost:6379")
	}

	return r.appContext
}

// SetAppContextRedis definition
func (r *Router) SetAppContextRedis(redis *redis.Pool) {
	r.appContext.Redis = redis
}

// GET draw a get handler
func (r *Router) GET(path string, route Route) {
	r.Handle("GET", path, route(r.AppContext()))
}

// POST draw a post handler
func (r *Router) POST(path string, route Route) {
	r.Handle("POST", path, route(r.AppContext()))
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

func readPrivateKey(path string) (interface{}, error) {
	privateKeyData, err := readFile(path)
	if err != nil {
		return nil, err
	}

	return jose.LoadPrivateKey([]byte(privateKeyData))
}

func readPublicKey(path string) (interface{}, error) {
	publicKeyData, err := readFile(path)
	if err != nil {
		return nil, err
	}

	return jose.LoadPublicKey([]byte(publicKeyData))
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
