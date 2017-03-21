package handlers_test

import (
	"fmt"
	"github.com/duckclick/wing/handlers"
	helpers "github.com/duckclick/wing/testing"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWhenJsonPayloadIsInvalid(t *testing.T) {
	appContext := helpers.CreateFileExporterAppContext()
	handler := handlers.CollectHandler(appContext)
	params := httprouter.Params{}

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", strings.NewReader("{")) // invalid
	req.Header.Set("Content-Type", "application/json")

	handler(rr, req, params)
	assert.Equal(t, 422, rr.Code, "should respond with 422 to invalid payload")
	assert.Equal(t, `{"message": "Invalid JSON payload"}`, rr.Body.String(), "should respond with an error message")
}

func TestWhenBase64IsInvalid(t *testing.T) {
	appContext := helpers.CreateFileExporterAppContext()
	handler := handlers.CollectHandler(appContext)
	params := httprouter.Params{}

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/",
		strings.NewReader(`[{"type": "TrackDOM", "created_at": 12345, "payload": {"markup": "123"}}]`),
	)

	req.Header.Set("Content-Type", "application/json")
	handler(rr, req, params)

	assert.Equal(t, 422, rr.Code, "should respond with 422 to invalid base64 payload")
	assert.Equal(t, `{"message": "Invalid JSON payload"}`, rr.Body.String(), "should respond with an error message")
}

func TestWhenItSavesTheRequest(t *testing.T) {
	appContext := helpers.CreateFileExporterAppContext()
	handler := handlers.CollectHandler(appContext)
	params := httprouter.Params{}

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/",
		strings.NewReader(`[{"type": "TrackDOM", "created_at": 1480979268, "payload": {"markup": "PGh0bWw+PC9odG1sPg=="}}]`),
	)

	req.Header.Set("Content-Type", "application/json")
	handler(rr, req, params)

	assert.Equal(t, 201, rr.Code, "should respond with 201 to to valid request")
	assert.Equal(t, `{"recorded": true}`, rr.Body.String(), "should respond with valid json")

	request := http.Request{Header: http.Header{"Cookie": rr.HeaderMap["Set-Cookie"]}}
	_, err := request.Cookie(handlers.RecordIDCookieName)

	assert.Nil(t, err, fmt.Sprintf("expected h to create '%s' cookie", handlers.RecordIDCookieName))
}
