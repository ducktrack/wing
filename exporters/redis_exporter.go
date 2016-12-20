package exporters

import (
	"github.com/duckclick/wing/config"
	"github.com/garyburd/redigo/redis"
	"fmt"
)

type redisExporter struct {
	Config config.RedisExporter
	Pool   *redis.Pool
}

func NewRedisExporter(config config.RedisExporter) *redisExporter {
	exporter := &redisExporter{Config: config}
	exporter.Connect()
	return exporter
}

func (re *redisExporter) Connect() {
	re.Pool = redis.NewPool(
		func() (redis.Conn, error) {
			host := re.Config.Host
			port := re.Config.Host
			c, err := redis.Dial("tcp", fmt.Sprintf("%v:%v", host, port))
			if err != nil {
				return nil, err
			}
			return c, err
		},
		re.Config.PoolSize,
	)
}

func (re *redisExporter) Stop() error {
	return re.Pool.Close()
}

func (re *redisExporter) Export(trackEntry *TrackEntry, recordId string) error {
	return nil
}
