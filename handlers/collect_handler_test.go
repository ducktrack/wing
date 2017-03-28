package handlers_test

import (
	"fmt"
	"github.com/duckclick/wing/events"
	"github.com/duckclick/wing/handlers"
	helpers "github.com/duckclick/wing/testing"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MyMockedExporter struct {
	mock.Mock
}

func (m *MyMockedExporter) Initialize() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MyMockedExporter) Export(trackable events.Trackable, recordID string) error {
	args := m.Called(trackable, recordID)
	return args.Error(0)
}

func (m *MyMockedExporter) Stop() error {
	args := m.Called()
	return args.Error(0)
}

type CollectHandlerTestSuite struct {
	suite.Suite
	handler    httprouter.Handle
	params     httprouter.Params
	appContext *handlers.AppContext
}

func (suite *CollectHandlerTestSuite) SetupTest() {
	suite.appContext = helpers.CreateFileExporterAppContext()
	suite.handler = handlers.CollectHandler(suite.appContext)
	suite.params = httprouter.Params{}
}

func (suite *CollectHandlerTestSuite) TestWhenJsonPayloadIsInvalid() {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", strings.NewReader("{")) // invalid
	req.Header.Set("Content-Type", "application/json")

	suite.handler(rr, req, suite.params)
	assert.Equal(suite.T(), 422, rr.Code, "should respond with 422 to invalid payload")
	assert.Equal(suite.T(), `{"message": "Invalid JSON payload"}`, rr.Body.String(), "should respond with an error message")
}

func (suite *CollectHandlerTestSuite) TestWhenBase64IsInvalid() {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/",
		strings.NewReader(`[{"type": "TrackDOM", "created_at": 12345, "payload": {"markup": "123"}}]`),
	)

	req.Header.Set("Content-Type", "application/json")
	suite.handler(rr, req, suite.params)

	assert.Equal(suite.T(), 422, rr.Code, "should respond with 422 to invalid base64 payload")
	assert.Equal(suite.T(), `{"message": "Invalid JSON payload"}`, rr.Body.String(), "should respond with an error message")
}

func (suite *CollectHandlerTestSuite) TestWhenExporterFail() {
	mockedExporter := new(MyMockedExporter)
	mockedExporter.
		On(
			"Export",
			mock.AnythingOfType("*events.TrackDOM"),
			mock.AnythingOfType("string"),
		).
		Return(errors.New("Failed"))

	suite.appContext.Exporter = mockedExporter
	suite.handler = handlers.CollectHandler(suite.appContext)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/",
		strings.NewReader(`[{"type": "TrackDOM", "created_at": 1480979268, "payload": {"markup": "PGh0bWw+PC9odG1sPg=="}}]`),
	)

	req.Header.Set("Content-Type", "application/json")
	suite.handler(rr, req, suite.params)

	assert.Equal(suite.T(), 422, rr.Code, "should respond with 422")
	assert.Equal(suite.T(), `{"message": "Failed to export event"}`, rr.Body.String(), "should respond with an error message")
}

func (suite *CollectHandlerTestSuite) TestWhenItSavesTheRequest() {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/",
		strings.NewReader(`[{"type": "TrackDOM", "created_at": 1480979268, "payload": {"markup": "PGh0bWw+PC9odG1sPg=="}}]`),
	)

	req.Header.Set("Content-Type", "application/json")
	suite.handler(rr, req, suite.params)

	assert.Equal(suite.T(), 201, rr.Code, "should respond with 201 to to valid request")
	assert.Equal(suite.T(), `{"recorded": true}`, rr.Body.String(), "should respond with valid json")
}

func (suite *CollectHandlerTestSuite) TestCookieConfiguration() {
	suite.appContext.Config.RecordIdCookieName = "custom_cookie_name"
	suite.handler = handlers.CollectHandler(suite.appContext)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/",
		strings.NewReader(`[{"type": "TrackDOM", "created_at": 1480979268, "payload": {"markup": "PGh0bWw+PC9odG1sPg=="}}]`),
	)

	req.Header.Set("Content-Type", "application/json")
	suite.handler(rr, req, suite.params)

	request := http.Request{Header: http.Header{"Cookie": rr.HeaderMap["Set-Cookie"]}}
	_, err := request.Cookie("custom_cookie_name")

	assert.Nil(suite.T(), err, fmt.Sprintf("expected h to create '%s' cookie", suite.appContext.Config.RecordIdCookieName))
}

func TestCollectHandler(t *testing.T) {
	suite.Run(t, new(CollectHandlerTestSuite))
}
