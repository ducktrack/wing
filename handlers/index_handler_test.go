package handlers_test

import (
	"github.com/duckclick/wing/handlers"
	helpers "github.com/duckclick/wing/testing"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type IndexHandlerTestSuite struct {
	suite.Suite
	handler httprouter.Handle
	params  httprouter.Params
}

func (suite *IndexHandlerTestSuite) SetupTest() {
	suite.handler = handlers.IndexHandler(helpers.CreateFileExporterAppContext())
	suite.params = httprouter.Params{}
}

func (suite *IndexHandlerTestSuite) TestReturns200() {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	suite.handler(rr, req, suite.params)
	assert.Equal(suite.T(), 200, rr.Code, "should respond with 200 to valid request")
	assert.Equal(suite.T(), `{"name": "wing"}`, rr.Body.String(), "should respond with valid json")
}

func TestIndexHandler(t *testing.T) {
	suite.Run(t, new(IndexHandlerTestSuite))
}
