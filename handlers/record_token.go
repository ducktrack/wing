package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/duckclick/wing/config"
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
	jose "gopkg.in/square/go-jose.v1"
	"net/url"
)

// RecordTokenRedisNamespace definition
const RecordTokenRedisNamespace = "record_tokens"

// RecordToken definition
type RecordToken struct {
	ID   string `json:"id"`
	Host string `json:"host"`
}

// FindRecordTokenByHost find the record token in redis under the namespace defined by RecordTokenRedisNamespace (e.g: `record_tokens/<host>`)
// e.g:
//   GET "record_tokens/github.com"
//   # "a38103c9-a19e-409d-94de-0e3c86085a5a"
func FindRecordTokenByHost(redisPool *redis.Pool, originHost string) (RecordToken, error) {
	host, err := extractOriginHost(originHost)
	if err != nil {
		return RecordToken{}, err
	}

	conn := redisPool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s/%s", RecordTokenRedisNamespace, host)
	id, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return RecordToken{}, errors.Wrapf(err, "Record token for host '%s' not found", host)
	}

	recordToken := RecordToken{ID: id, Host: host}
	return recordToken, nil
}

// VerifyRecordToken verifies if the token was issued for the given host
func VerifyRecordToken(jwtToken *jwt.Token, originHost string, privateKey config.PrivateKey) error {
	host, err := extractOriginHost(originHost)
	if err != nil {
		return err
	}

	claims, ok := jwtToken.Claims.(*jwt.StandardClaims)
	recordToken, err := decryptEncodedRecordToken(claims.Subject, privateKey)
	if err != nil {
		return errors.Wrap(err, "Failed to decrypt encoded record token")
	}

	if !ok || host != recordToken.Host {
		return errors.Errorf("Failed to match session token host ('%s') with Origin header host ('%s')", recordToken.Host, host)
	}

	return nil
}

// EncodeAndEncryptRecordToken definition
func EncodeAndEncryptRecordToken(recordToken RecordToken, publicKey config.PublicKey) (string, error) {
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

func decryptEncodedRecordToken(encryptedPayload string, privateKey config.PrivateKey) (RecordToken, error) {
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
