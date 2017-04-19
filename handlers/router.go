package handlers

import (
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/exporters"
	"github.com/garyburd/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"time"
)

// AppContext definition
type AppContext struct {
	Config        *config.Config
	Exporter      exporters.Exporter
	Redis         *redis.Pool
	JWEPrivateKey config.PrivateKey
	JWEPublicKey  config.PublicKey
}

// Router definition
type Router struct {
	*httprouter.Router
	appContext *AppContext
}

// Route definition
type Route func(*AppContext) httprouter.Handle

// NewAppContext creates a new app context
func NewAppContext(wingConfig *config.Config, exporter exporters.Exporter) (*AppContext, error) {
	privateKey, publicKey, err := config.LoadJWEKeys(wingConfig)
	if err != nil {
		return nil, err
	}

	return &AppContext{
		Config:        wingConfig,
		Exporter:      exporter,
		JWEPrivateKey: privateKey,
		JWEPublicKey:  publicKey,
	}, nil
}

// NewRouter creates a new router with the app context
func NewRouter(wingConfig *config.Config, exporter exporters.Exporter) (*Router, error) {
	appContext, err := NewAppContext(wingConfig, exporter)
	if err != nil {
		return nil, err
	}

	return &Router{
		Router:     httprouter.New(),
		appContext: appContext,
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
