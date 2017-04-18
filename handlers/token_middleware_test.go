package handlers_test

import (
	"fmt"
	"github.com/duckclick/wing/handlers"
	helpers "github.com/duckclick/wing/testing"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TokenMiddlewareTestHandler(appContext *handlers.AppContext) httprouter.Handle {
	return func(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		response.WriteHeader(http.StatusOK)
		fmt.Fprint(response, `{"status": "ok"}`)
	}
}

type TokenMiddlewareTestSuite struct {
	suite.Suite
	handler          httprouter.Handle
	params           httprouter.Params
	appContext       *handlers.AppContext
	originHost       string
	recordKey        string
	encryptedPayload string
	duration         time.Duration
}

func (suite *TokenMiddlewareTestSuite) SetupTest() {
	suite.appContext = helpers.CreateFileExporterAppContext()
	suite.params = httprouter.Params{}
	suite.handler = handlers.TokenMiddleware(TokenMiddlewareTestHandler)(suite.appContext)

	suite.originHost = "duckclick.com"
	suite.recordKey = "a38103c9-a19e-409d-94de-0e3c86085a5a"

	encryptedPayload, _ := handlers.EncodeAndEncryptRecordToken(handlers.RecordToken{
		ID:   suite.recordKey,
		Host: suite.originHost,
	}, suite.appContext.JWEPublicKey)

	suite.encryptedPayload = encryptedPayload
	suite.duration = time.Duration(10) * time.Second
}

func (suite *TokenMiddlewareTestSuite) TestReturns403WhenTokenIsInvalid() {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid")

	suite.handler(rr, req, suite.params)
	assert.Equal(suite.T(), 403, rr.Code, "should respond with 403")
	assert.Equal(suite.T(), `{"message": "Forbidden"}`, rr.Body.String(), "should respond with valid json")
}

func (suite *TokenMiddlewareTestSuite) TestReturns403WhenRecordTokenIsInvalid() {
	encryptedPayload, _ := handlers.EncodeAndEncryptRecordToken(handlers.RecordToken{
		ID:   "wrong-key",
		Host: suite.originHost,
	}, suite.appContext.JWEPublicKey)

	token := handlers.CreateToken(encryptedPayload, suite.duration)
	tokenString, _ := handlers.SignToken(token, suite.appContext.Config.SessionTokenSecret)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	suite.handler(rr, req, suite.params)
	assert.Equal(suite.T(), 403, rr.Code, "should respond with 403")
	assert.Equal(suite.T(), `{"message": "Forbidden"}`, rr.Body.String(), "should respond with valid json")
}

func (suite *TokenMiddlewareTestSuite) TestCallsTheOriginalHandlerWhenTokenAndRecordTokenAreCorrect() {
	token := handlers.CreateToken(suite.encryptedPayload, suite.duration)
	tokenString, _ := handlers.SignToken(token, suite.appContext.Config.SessionTokenSecret)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	req.Header.Set("Origin", fmt.Sprintf("https://%s", suite.originHost))

	suite.handler(rr, req, suite.params)
	assert.Equal(suite.T(), 200, rr.Code, "should respond with 200")
	assert.Equal(suite.T(), `{"status": "ok"}`, rr.Body.String(), "should respond with valid json")
}

func TestTokenMiddleware(t *testing.T) {
	suite.Run(t, new(TokenMiddlewareTestSuite))
}
