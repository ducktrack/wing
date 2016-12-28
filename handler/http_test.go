package handler

import (
	"fmt"
	"github.com/duckclick/wing/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"github.com/duckclick/wing/exporters"
)

var appConfig config.Config
var fileExporter exporters.Exporter

func TestMain(m *testing.M) {
	appConfig = config.Config{
		Exporter: "file",
		FileExporter: config.FileExporter{
			Folder: "/tmp/test/track_entries",
		},
	}
	fileExporter, _ = exporters.Lookup(&appConfig)

	os.Exit(m.Run())
}

func TestWhenRequestMethodOptions(t *testing.T) {
	rr := httptest.NewRecorder()
	h := &TrackEntryHandler{Config: appConfig, Exporter: fileExporter}

	req, _ := http.NewRequest("OPTIONS", "/", nil)
	h.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code, "should respond with 200 to OPTIONS request")
	assert.Equal(t, "", rr.Body.String(), "should respond with an empty body")
}

func TestWhenRequestMethodDifferentThanPost(t *testing.T) {
	for _, method := range []string{"GET", "PUT", "DELETE", "PATCH"} {
		rr := httptest.NewRecorder()
		h := &TrackEntryHandler{Config: appConfig, Exporter: fileExporter}
		req, _ := http.NewRequest(method, "/", nil)
		h.ServeHTTP(rr, req)

		assert.Equal(t, 405, rr.Code, "should respond with 405 to unhandled method")
		assert.Equal(t, `{"message": "Method Not Allowed"}`, rr.Body.String(), "should respond with an error message")
	}
}

func TestWhenJsonPayloadIsInvalid(t *testing.T) {
	rr := httptest.NewRecorder()
	h := &TrackEntryHandler{Config: appConfig, Exporter: fileExporter}

	req, _ := http.NewRequest("POST", "/", strings.NewReader("{")) // invalid
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(rr, req)

	assert.Equal(t, 422, rr.Code, "should respond with 422 to invalid payload")
	assert.Equal(t, `{"message": "Invalid JSON payload"}`, rr.Body.String(), "should respond with an error message")
}

func TestWhenBase64IsInvalid(t *testing.T) {
	rr := httptest.NewRecorder()
	h := &TrackEntryHandler{Config: appConfig, Exporter: fileExporter}

	req, _ := http.NewRequest(
		"POST",
		"/",
		strings.NewReader(`{"created_at": 12345, "markup": "123"}`),
	)

	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(rr, req)

	assert.Equal(t, 422, rr.Code, "should respond with 422 to invalid base64 payload")
	assert.Equal(t, `{"message": "Failed to export track entry"}`, rr.Body.String(), "should respond with an error message")
}

func TestWhenItSavesTheRequest(t *testing.T) {
	rr := httptest.NewRecorder()
	h := &TrackEntryHandler{Config: appConfig, Exporter: fileExporter}

	req, _ := http.NewRequest(
		"POST",
		"/",
		strings.NewReader(`{"created_at": 1480979268, "markup": "PGh0bWw+PC9odG1sPg=="}`),
	)

	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(rr, req)

	assert.Equal(t, 201, rr.Code, "should respond with 201 to to valid request")
	assert.Equal(t, "", rr.Body.String(), "should respond with an empty body")

	request := http.Request{Header: http.Header{"Cookie": rr.HeaderMap["Set-Cookie"]}}
	_, err := request.Cookie(RECORD_ID_COOKIE_NAME)

	assert.Nil(t, err, fmt.Sprintf("expected h to create '%s' cookie", RECORD_ID_COOKIE_NAME))
}
