package exporters

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/trackentry"
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

// RedisExporter definition
type RedisExporter struct {
	config config.RedisExporter
	pool   *redis.Pool
}

// NewRedisExporter is the construtor of RedisExporter
func NewRedisExporter(config config.RedisExporter) *RedisExporter {
	exporter := &RedisExporter{config: config}
	exporter.Connect()
	return exporter
}

// Connect establishes the connection with the redis host
func (re *RedisExporter) Connect() {
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

// Stop closes the connection pool
func (re *RedisExporter) Stop() error {
	return re.pool.Close()
}

// Export saves the entry using HSET with recordID as the key. The field name is the created at value of the TrackEntry.
// To list all fields use HGETALL <recordID>, example: hgetall "593a177d-e250-4fc2-a6a4-5b0ec33ed56a"
func (re *RedisExporter) Export(trackEntry *trackentry.TrackEntry, recordID string) error {
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
	log.Infof("Storing redis entry at: %s, %s", recordID, createdAtStr)
	reply, err := conn.Do("HSET", recordID, createdAtStr, markup)
	return errors.Wrapf(err, "Failed to store track entry in redis, error: %s, reply: %s", err, reply)
}