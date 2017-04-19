package testing

import (
	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"
)

func CreateRedisConn() *redigomock.Conn {
	return redigomock.NewConn()
}

func CreateRedisPool(conn *redigomock.Conn) *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return conn, nil
		},
	}
}
