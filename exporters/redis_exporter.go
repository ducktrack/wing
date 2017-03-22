package exporters

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/events"
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

// RedisExporter definition
type RedisExporter struct {
	Config config.RedisExporter
	Pool   *redis.Pool
}

// NewRedisExporter is the construtor of RedisExporter
func NewRedisExporter(config config.RedisExporter) *RedisExporter {
	return &RedisExporter{Config: config}
}

// Initialize establishes and verify the connection with the redis host
func (re *RedisExporter) Initialize() error {
	connString := fmt.Sprintf("%s:%d", re.Config.Host, re.Config.Port)
	log.Infof("Initializing RedisExporter (connection string '%s')", connString)

	if re.Pool == nil {
		re.Pool = createConnectionPool(connString)
	}

	conn := re.Pool.Get()
	defer conn.Close()

	reply, err := redis.String(conn.Do("PING"))
	if err != nil || reply != "PONG" {
		if err == nil {
			err = errors.Errorf("Wrong reply, expected: 'PONG', received: '%s'", reply)
		}
		return errors.Wrap(err, "Failed to test connection with redis")
	}

	return nil
}

// Stop closes the connection pool
func (re *RedisExporter) Stop() error {
	return re.Pool.Close()
}

// Export saves the entry using HSET with recordID as the key. The field name is the created at value
// of the TrackEntry.
// To list all fields use HGETALL <recordID>, example: hgetall "593a177d-e250-4fc2-a6a4-5b0ec33ed56a"
func (re *RedisExporter) Export(trackable events.Trackable, recordID string) error {
	event := trackable.GetEvent()
	json, err := trackable.ToJSON()
	if err != nil {
		return errors.Wrap(err, "Failed to encode json")
	}

	if re.Pool == nil {
		return errors.New("Not connected to Redis, must connect first")
	}
	conn := re.Pool.Get()
	defer conn.Close()

	createdAtStr := strconv.Itoa(event.CreatedAt)
	log.Infof("Storing redis entry at: %s, %s", recordID, createdAtStr)
	reply, err := conn.Do("HSET", recordID, createdAtStr, json)
	return errors.Wrapf(err, "Failed to store track entry in redis, error: %s, reply: %s", err, reply)
}

func createConnectionPool(connString string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", connString)
		},
	}
}
