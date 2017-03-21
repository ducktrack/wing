package handlers_test

import (
	"github.com/duckclick/wing/handlers"
	helpers "github.com/duckclick/wing/testing"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReturns200(t *testing.T) {
	appContext := helpers.CreateFileExporterAppContext()
	handler := handlers.IndexHandler(appContext)
	params := httprouter.Params{}

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	handler(rr, req, params)
	assert.Equal(t, 200, rr.Code, "should respond with 200 to to valid request")
	assert.Equal(t, `{"name": "wing"}`, rr.Body.String(), "should respond with valid json")
}
