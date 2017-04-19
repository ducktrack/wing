package handlers_test

import (
	"fmt"
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/handlers"
	helpers "github.com/duckclick/wing/testing"
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
	"time"
)

type RecordTokenTestSuite struct {
	suite.Suite
	conn             *redigomock.Conn
	redis            *redis.Pool
	config           config.Config
	privateKey       config.PrivateKey
	publicKey        config.PublicKey
	originHost       string
	recordKey        string
	redisKey         string
	encryptedPayload string
	duration         time.Duration
}

func (suite *RecordTokenTestSuite) SetupTest() {
	suite.conn = helpers.CreateRedisConn()
	suite.redis = helpers.CreateRedisPool(suite.conn)
	suite.config = helpers.CreateBasicConfig()

	privateKey, publicKey, _ := config.LoadJWEKeys(&suite.config)
	suite.privateKey = privateKey
	suite.publicKey = publicKey

	suite.originHost = "duckclick.com"
	suite.recordKey = "a38103c9-a19e-409d-94de-0e3c86085a5a"
	suite.redisKey = fmt.Sprintf("%s/%s", handlers.RecordTokenRedisNamespace, suite.originHost)

	encryptedPayload, _ := handlers.EncodeAndEncryptRecordToken(handlers.RecordToken{
		ID:   suite.recordKey,
		Host: suite.originHost,
	}, suite.publicKey)

	suite.encryptedPayload = encryptedPayload
	suite.duration = time.Duration(10) * time.Second
}

func (suite *RecordTokenTestSuite) TestFindRecordTokenByHost() {
	cmd := suite.conn.Command("GET", suite.redisKey).Expect(suite.recordKey)
	recordToken, err := handlers.FindRecordTokenByHost(suite.redis, fmt.Sprintf("https://%s", suite.originHost))
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), recordToken.ID, suite.recordKey)
	assert.Equal(suite.T(), recordToken.Host, suite.originHost)
	assert.True(suite.T(), suite.conn.Stats(cmd) == 1, "Command should have been called")
}

func (suite *RecordTokenTestSuite) TestFindRecordTokenByHostWhenOriginHostIsBlank() {
	_, err := handlers.FindRecordTokenByHost(suite.redis, "")
	assert.NotNil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(err.Error(), "Origin header is blank"), err.Error())
}

func (suite *RecordTokenTestSuite) TestFindRecordTokenByHostWhenOriginHostIsNotAValidURL() {
	_, err := handlers.FindRecordTokenByHost(suite.redis, "http://[::1]%23")
	assert.NotNil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(err.Error(), "Failed to parse origin host"), err.Error())
}

func (suite *RecordTokenTestSuite) TestFindRecordTokenByHostWhenHostDoesNotExist() {
	cmd := suite.conn.Command("GET", suite.redisKey).Expect(nil)
	_, err := handlers.FindRecordTokenByHost(suite.redis, fmt.Sprintf("https://%s", suite.originHost))
	assert.NotNil(suite.T(), err)
	assert.True(suite.T(), suite.conn.Stats(cmd) == 1, "Command should have been called")
	assert.True(suite.T(), strings.Contains(err.Error(), "Record token for host 'duckclick.com' not found"), err.Error())
}

func (suite *RecordTokenTestSuite) TestFindRecordTokenByHostWhenRedisReturnsAnError() {
	cmd := suite.conn.Command("GET", suite.redisKey).ExpectError(errors.Errorf("Boom!"))
	_, err := handlers.FindRecordTokenByHost(suite.redis, fmt.Sprintf("https://%s", suite.originHost))
	assert.NotNil(suite.T(), err)
	assert.True(suite.T(), suite.conn.Stats(cmd) == 1, "Command should have been called")
	assert.True(suite.T(), strings.Contains(err.Error(), "Boom!"), err.Error())
}

func (suite *RecordTokenTestSuite) TestVerifyRecordToken() {
	token := handlers.CreateToken(suite.encryptedPayload, suite.duration)
	err := handlers.VerifyRecordToken(token, fmt.Sprintf("https://%s", suite.originHost), suite.privateKey)
	assert.Nil(suite.T(), err)
}

func (suite *RecordTokenTestSuite) TestVerifyRecordTokenWhenOriginHostIsBlank() {
	token := handlers.CreateToken(suite.encryptedPayload, suite.duration)
	err := handlers.VerifyRecordToken(token, "", suite.privateKey)
	assert.NotNil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(err.Error(), "Origin header is blank"), err.Error())
}

func (suite *RecordTokenTestSuite) TestVerifyRecordTokenWhenOriginHostIsNotAValidURL() {
	token := handlers.CreateToken(suite.encryptedPayload, suite.duration)
	err := handlers.VerifyRecordToken(token, "http://[::1]%23", suite.privateKey)
	assert.NotNil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(err.Error(), "Failed to parse origin host"), err.Error())
}

func (suite *RecordTokenTestSuite) TestVerifyRecordTokenWhenTheEncryptedPayloadIsInvalid() {
	encryptedPayload := `{"protected":"eyJhb","encrypted_key":"PpxfsRXoV2eKxWBV","iv":"1M3I7P_Oe51bVQIq","ciphertext":"SwGtN_iPRvKB7","tag":"GXJdsGN"}`
	token := handlers.CreateToken(encryptedPayload, suite.duration)
	err := handlers.VerifyRecordToken(token, fmt.Sprintf("https://%s", suite.originHost), suite.privateKey)
	assert.NotNil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(err.Error(), "Failed to decrypt encoded record token"), err.Error())
}

func (suite *RecordTokenTestSuite) TestVerifyRecordTokenWhenTokenDoesNotMatchOriginHost() {
	token := handlers.CreateToken(suite.encryptedPayload, suite.duration)
	err := handlers.VerifyRecordToken(token, "https://github.com", suite.privateKey)
	assert.NotNil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(err.Error(), "Failed to match session token host"), err.Error())
}

func TestRecordToken(t *testing.T) {
	suite.Run(t, new(RecordTokenTestSuite))
}
