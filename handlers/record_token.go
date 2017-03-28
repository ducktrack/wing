package handlers

import (
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
	"net/url"
)

// RecordToken definition
type RecordToken struct {
	ID   string `json:"id"`
	Host string `json:"host"`
}

// FindRecordTokenByHost definition
func FindRecordTokenByHost(redisPool *redis.Pool, originHost string) (RecordToken, error) {
	host, err := extractOriginHost(originHost)
	if err != nil {
		return RecordToken{}, err
	}

	conn := redisPool.Get()
	defer conn.Close()

	id, err := redis.String(conn.Do("GET", "record_tokens/"+host))
	if err != nil {
		return RecordToken{}, errors.Wrapf(err, "Record token for host '%s' not found", host)
	}

	recordToken := RecordToken{ID: id, Host: host}
	return recordToken, nil
}

// VerifyRecordToken definition
func VerifyRecordToken(jwtToken *jwt.Token, originHost string) error {
	host, err := extractOriginHost(originHost)
	if err != nil {
		return err
	}

	claims, ok := jwtToken.Claims.(*jwt.StandardClaims)
	var recordToken RecordToken
	err = json.Unmarshal([]byte(claims.Subject), &recordToken)

	if err != nil {
		return errors.Wrap(err, "Failed to decode JWT Subject claim")
	}

	if !ok || host != recordToken.Host {
		return errors.Wrapf(err, "Failed to match session token host ('%s') with Origin header host ('%s')", recordToken.Host, host)
	}

	return nil
}

func extractOriginHost(originHeader string) (string, error) {
	if len(originHeader) == 0 {
		return "", errors.Errorf("Origin header is blank")
	}

	originURL, err := url.Parse(originHeader)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to parse origin host ('%s') ", originHeader)
	}

	return originURL.Host, nil
}
