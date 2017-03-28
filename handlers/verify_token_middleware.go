package handlers

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// VerifyTokenMiddleware definition
func VerifyTokenMiddleware(route Route) Route {
	return func(appContext *AppContext) httprouter.Handle {
		return func(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
			authHeader := request.Header.Get("Authorization")
			token, err := ParseAndVerifyAuthenticationHeader(authHeader, appContext.Config.SessionTokenSecret)

			if err != nil {
				log.WithError(err).Error("Failed to parse authentication header")
				response.WriteHeader(http.StatusForbidden)
				fmt.Fprint(response, `{"message": "Forbidden"}`)
				return
			}

			err = VerifyRecordToken(token, request.Header.Get("Origin"), appContext.JWEPrivateKey)

			if err != nil {
				log.WithError(err).Error("Failed to verify record token")
				response.WriteHeader(http.StatusForbidden)
				fmt.Fprint(response, `{"message": "Forbidden"}`)
				return
			}

			route(appContext)(response, request, params)
		}
	}
}
