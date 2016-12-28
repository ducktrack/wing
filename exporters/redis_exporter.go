package exporters

import (
	"github.com/duckclick/wing/config"
	"github.com/garyburd/redigo/redis"
	"fmt"
	"time"
	"encoding/base64"
	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
	"strings"
	"errors"
	"strconv"
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
	re.pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			connString := fmt.Sprintf("%s:%d", re.config.Host, re.config.Port)
			log.Infof("Redis connection string: %s", connString)

			c, err := redis.Dial("tcp", connString)
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
}

func (re *redisExporter) Stop() error {
	return re.pool.Close()
}

func (re *redisExporter) Export(trackEntry *TrackEntry, recordId string) error {
	htmlBytes, err := base64.StdEncoding.DecodeString(trackEntry.Markup)
	if err != nil {
		return errors.New("Invalid base64 payload")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(htmlBytes)))
	if err != nil {
		return errors.New("Failed to parse HTML")
	}

	scripts := doc.Find("script")
	scripts.Each(func(i int, s *goquery.Selection) { s.Remove() })

	content, err := doc.Html()
	if err != nil {
		return errors.New("Failed to generate secure HTML")
	}

	if re.pool == nil {
		return errors.New("Not connected, must connect first")
	}
	conn := re.pool.Get()
	defer conn.Close()

	createdAtStr := strconv.Itoa(trackEntry.CreatedAt)
	log.Infof("Storing redis entry at: %s, %s", recordId, createdAtStr)
	reply, err := conn.Do("HSET", recordId, createdAtStr, content)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to store track entry in redis, error: %s, reply: %s", err, reply))
	}

	return nil
}
