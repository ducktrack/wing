package handlers

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// IndexHandler definition
func IndexHandler(appContext *AppContext) httprouter.Handle {
	return func(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		response.WriteHeader(http.StatusOK)
		fmt.Fprint(response, `{"name": "wing"}`)
	}
}
