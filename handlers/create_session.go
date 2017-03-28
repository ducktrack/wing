package handlers

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

// CreateSessionHandler definition
func CreateSessionHandler(appContext *AppContext) httprouter.Handle {
	return func(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
		recordToken, err := FindRecordTokenByHost(appContext.Redis, request.Header.Get("Origin"))
		if err != nil {
			log.WithError(err).Error("Failed to find record token")
			response.WriteHeader(http.StatusForbidden)
			fmt.Fprint(response, `{"message": "Failed to create session"}`)
			return
		}

		token, err := CreateToken(recordToken, time.Duration(30)*time.Minute)
		if err != nil {
			log.WithError(err).Errorf("Failed to create token for RecordToken (%v)", recordToken)
			response.WriteHeader(http.StatusForbidden)
			fmt.Fprint(response, `{"message": "Failed to create session"}`)
			return
		}

		tokenString, err := SignToken(token, appContext.Config.SessionTokenSecret)
		if err != nil {
			log.WithError(err).Error("Failed to sign JWT")
			response.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(response, `{"message": "Failed to create session"}`)
			return
		}

		response.WriteHeader(http.StatusCreated)
		fmt.Fprint(response, fmt.Sprintf(`{"session_token": "%s"}`, tokenString))
	}
}
