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
	handler httprouter.Handle
	params  httprouter.Params
}

func (suite *CollectHandlerTestSuite) SetupTest() {
	suite.handler = handlers.CollectHandler(helpers.CreateFileExporterAppContext())
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
	appContext := helpers.CreateFileExporterAppContext()

	mockedExporter := new(MyMockedExporter)
	mockedExporter.
		On(
			"Export",
			mock.AnythingOfType("*events.TrackDOM"),
			mock.AnythingOfType("string"),
		).
		Return(errors.New("Failed"))

	appContext.Exporter = mockedExporter
	suite.handler = handlers.CollectHandler(appContext)

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

	request := http.Request{Header: http.Header{"Cookie": rr.HeaderMap["Set-Cookie"]}}
	_, err := request.Cookie(handlers.RecordIDCookieName)

	assert.Nil(suite.T(), err, fmt.Sprintf("expected h to create '%s' cookie", handlers.RecordIDCookieName))
}

func TestCollectHandler(t *testing.T) {
	suite.Run(t, new(CollectHandlerTestSuite))
}
