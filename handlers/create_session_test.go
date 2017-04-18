package handlers_test

import (
	"encoding/json"
	"fmt"
	"github.com/duckclick/wing/handlers"
	helpers "github.com/duckclick/wing/testing"
	"github.com/garyburd/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type CreateSessionHandlerTestSuite struct {
	suite.Suite
	handler          httprouter.Handle
	params           httprouter.Params
	mockedConnection *redigomock.Conn
	redisPool        *redis.Pool
}

func (suite *CreateSessionHandlerTestSuite) SetupTest() {
	suite.params = httprouter.Params{}
	suite.mockedConnection = redigomock.NewConn()
	suite.redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return suite.mockedConnection, nil
		},
	}

	appContext := helpers.CreateFileExporterAppContext()
	appContext.Redis = suite.redisPool
	suite.handler = handlers.CreateSessionHandler(appContext)
}

func (suite *CreateSessionHandlerTestSuite) TestWhenCalledWithAnInvalidOriginItReturns403() {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", nil)
	req.Header.Set("Origin", "https://foo.example.com")

	key := fmt.Sprintf("%s/%s", handlers.RecordTokenRedisNamespace, "foo.example.com")
	suite.mockedConnection.Command("GET", key).Expect(nil)

	suite.handler(rr, req, suite.params)
	assert.Equal(suite.T(), 403, rr.Code, "should respond with 403")
	assert.Equal(suite.T(), `{"message": "Failed to create session"}`, rr.Body.String(), "should respond with valid json")
}

func (suite *CreateSessionHandlerTestSuite) TestReturns201WhenOriginIsValid() {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", nil)
	req.Header.Set("Origin", "https://valid.example.com")

	key := fmt.Sprintf("%s/%s", handlers.RecordTokenRedisNamespace, "valid.example.com")
	suite.mockedConnection.Command("GET", key).Expect("a38103c9-a19e-409d-94de-0e3c86085a5a")

	type SessionJSON struct {
		AccessToken string `json:"access_token"`
	}

	suite.handler(rr, req, suite.params)
	assert.Equal(suite.T(), 201, rr.Code, "should respond with 201")

	var jsonPayload SessionJSON
	err := json.Unmarshal([]byte(rr.Body.String()), &jsonPayload)
	assert.Nil(suite.T(), err, "should be a valid JSON")
	assert.NotNil(suite.T(), jsonPayload.AccessToken, "should be present")
}

func TestCreateSessionHandler(t *testing.T) {
	suite.Run(t, new(CreateSessionHandlerTestSuite))
}
