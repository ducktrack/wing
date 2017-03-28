package handlers

import (
	"bytes"
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
	jose "gopkg.in/square/go-jose.v1"
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
func VerifyRecordToken(jwtToken *jwt.Token, originHost string, privateKey interface{}) error {
	host, err := extractOriginHost(originHost)
	if err != nil {
		return err
	}

	claims, ok := jwtToken.Claims.(*jwt.StandardClaims)
	recordToken, err := DecryptEncodedRecordToken(claims.Subject, privateKey)
	if err != nil {
		return errors.Wrap(err, "Failed to decrypt encoded record token")
	}

	if !ok || host != recordToken.Host {
		return errors.Wrapf(err, "Failed to match session token host ('%s') with Origin header host ('%s')", recordToken.Host, host)
	}

	return nil
}

// EncodeAndEncryptRecordToken definition
func EncodeAndEncryptRecordToken(recordToken RecordToken, publicKey interface{}) (string, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	err := encoder.Encode(recordToken)

	if err != nil {
		return "", errors.Wrap(err, "Failed to encode RecordToken to JSON")
	}

	encrypter, err := jose.NewEncrypter(jose.RSA_OAEP, jose.A128GCM, publicKey)
	if err != nil {
		return "", errors.Wrap(err, "Failed to create encrypter")
	}

	jweObject, err := encrypter.Encrypt(buffer.Bytes())
	if err != nil {
		return "", errors.Wrap(err, "Failed to encrypt JSON payload")
	}

	return jweObject.FullSerialize(), nil
}

// DecryptEncodedRecordToken definition
func DecryptEncodedRecordToken(encryptedPayload string, privateKey interface{}) (RecordToken, error) {
	jweObject, err := jose.ParseEncrypted(encryptedPayload)
	if err != nil {
		return RecordToken{}, errors.Wrap(err, "Failed to parse encrypted payload")
	}

	decryptedPayload, err := jweObject.Decrypt(privateKey)
	if err != nil {
		return RecordToken{}, errors.Wrap(err, "Failed to decrypt JWE Object")
	}

	var recordToken RecordToken
	err = json.Unmarshal([]byte(decryptedPayload), &recordToken)
	if err != nil {
		return RecordToken{}, errors.Wrap(err, "Failed to decode decrypted payload")
	}

	return recordToken, nil
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
