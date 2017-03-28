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
	Config   *config.Config
	Exporter exporters.Exporter
	Redis    *redis.Pool
}

// Router definition
type Router struct {
	*httprouter.Router
	*AppContext
}

// Route definition
type Route func(*AppContext) httprouter.Handle

// NewRouter creates a new router with the app context
func NewRouter(wingConfig *config.Config, exporter exporters.Exporter) *Router {
	return &Router{
		httprouter.New(),
		&AppContext{
			Config:   wingConfig,
			Exporter: exporter,
			Redis:    createRedisConnectionPool("localhost:6379"),
		},
	}
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
