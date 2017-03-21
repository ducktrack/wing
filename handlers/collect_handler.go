package handlers

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/duckclick/wing/events"
	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"io"
	"net/http"
	"time"
)

// cookie name
const RecordIDCookieName = "record_id"

// expiration time
const RecordIDExpiration = 2 * time.Hour

// CollectHandler definition
func CollectHandler(appContext *AppContext) httprouter.Handle {
	return func(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		recordCookie := recordCookie(response, request)
		recordID := recordCookie.Value
		trackableEvents, err := events.DecodeJSON(streamToByte(request.Body))

		if err != nil {
			response.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprint(response, `{"message": "Invalid JSON payload"}`)
			return
		}

		log.Infof("Tracking %d entries, record_id: %s", len(trackableEvents), recordID)

		for i := 0; i < len(trackableEvents); i++ {
			trackable := trackableEvents[i]
			event := trackable.GetEvent()
			log.Infof("%s, record_id: %s, created_at: %d, URL: %s", event.Type, recordID, event.CreatedAt, event.URL)
			err = appContext.Exporter.Export(trackable, recordID)

			if err != nil {
				log.WithError(err).Errorf("Failed to export event: %+v", event)
				response.WriteHeader(http.StatusUnprocessableEntity)
				fmt.Fprint(response, `{"message": "Failed to export event"}`)
				return
			}
		}

		response.WriteHeader(http.StatusCreated)
		fmt.Fprint(response, `{"recorded": true}`)
	}
}

func recordCookie(response http.ResponseWriter, request *http.Request) *http.Cookie {
	cookie, err := request.Cookie(RecordIDCookieName)

	if err != nil || cookie.Value == "" {
		cookie = &http.Cookie{
			Name:     RecordIDCookieName,
			Value:    uuid.NewV4().String(),
			Expires:  time.Now().Add(RecordIDExpiration),
			Path:     "/",
			HttpOnly: true,
		}

		http.SetCookie(response, cookie)
	}

	return cookie
}

func streamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}
