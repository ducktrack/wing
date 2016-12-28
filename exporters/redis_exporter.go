package exporters

import (
	"github.com/duckclick/wing/config"
	"github.com/garyburd/redigo/redis"
	"fmt"
	"time"
	log "github.com/Sirupsen/logrus"
	"strconv"
	"github.com/duckclick/wing/trackentry"
	"github.com/pkg/errors"
)

type redisExporter struct {
	config config.RedisExporter
	pool   *redis.Pool
}

func NewRedisExporter(config config.RedisExporter) *redisExporter {
	exporter := &redisExporter{config: config}
	exporter.Connect()
	return exporter
}

func (re *redisExporter) Connect() {
	connString := fmt.Sprintf("%s:%d", re.config.Host, re.config.Port)
	log.Infof("Redis connection string: %s", connString)

	re.pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", connString)
		},
	}
}

func (re *redisExporter) Stop() error {
	return re.pool.Close()
}

func (re *redisExporter) Export(trackEntry *trackentry.TrackEntry, recordId string) error {
	markup, err := trackEntry.Rinse()
	if err != nil {
		return errors.Wrap(err, "Failed to rinse the markup")
	}

	if re.pool == nil {
		return errors.New("Not connected to Redis, must connect first")
	}
	conn := re.pool.Get()
	defer conn.Close()

	createdAtStr := strconv.Itoa(trackEntry.CreatedAt)
	log.Infof("Storing redis entry at: %s, %s", recordId, createdAtStr)
	reply, err := conn.Do("HSET", recordId, createdAtStr, markup)
	return errors.Wrapf(err, "Failed to store track entry in redis, error: %s, reply: %s", err, reply)
}
