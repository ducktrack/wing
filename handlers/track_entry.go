package handlers

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/duckclick/wing/config"
	"github.com/duckclick/wing/exporters"
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

const RECORD_ID_COOKIE_NAME = "record_id"
const RECORD_ID_EXPIRATION = 2 * time.Hour

type TrackEntryHandler struct {
	Config *config.Config
}

func (h *TrackEntryHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	origin := request.Header.Get("Origin")
	response.Header().Set("Access-Control-Allow-Origin", origin)
	response.Header().Set("Access-Control-Allow-Headers", "content-type")
	response.Header().Set("Access-Control-Allow-Credentials", "true")
	response.Header().Set("Content-Type", "application/json")

	if request.Method == "OPTIONS" {
		return
	}

	if request.Method != "POST" {
		response.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(response, `{"message": "Method Not Allowed"}`)
		return
	}

	if h.Config == nil {
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, `{"message": "App not configured"}`)
		return
	}

	recordCookie := createRecordCookie(response, request)
	recordId := recordCookie.Value

	decoder := json.NewDecoder(request.Body)
	var trackEntry exporters.TrackEntry
	err := decoder.Decode(&trackEntry)
	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(response, `{"message": "Invalid JSON payload"}`)
		return
	}

	trackEntry.Origin = origin
	exporter, err := exporters.Lookup(h.Config)
	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(response, fmt.Sprintf(`{"message": "%s"}`, err.Error()))
		return
	}

	err = exporter.Export(&trackEntry, recordId)
	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(response, fmt.Sprintf(`{"message": "%s"}`, err.Error()))
		return
	}

	response.WriteHeader(http.StatusCreated)
	log.Infof("Tracking dom, record_id: %s, created_at: %d, origin: %s", recordId, trackEntry.CreatedAt, origin)
	fmt.Fprintf(response, fmt.Sprintf(`{"recorded": true}`))
}

func createRecordCookie(response http.ResponseWriter, request *http.Request) *http.Cookie {
	cookie, err := request.Cookie(RECORD_ID_COOKIE_NAME)

	if err != nil || cookie.Value == "" {
		cookie = &http.Cookie{
			Name:     RECORD_ID_COOKIE_NAME,
			Value:    uuid.NewV4().String(),
			Expires:  time.Now().Add(RECORD_ID_EXPIRATION),
			Path:     "/",
			HttpOnly: true,
		}

		http.SetCookie(response, cookie)
	}

	return cookie
}
