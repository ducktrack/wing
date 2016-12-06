package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const RECORD_ID_COOKIE_NAME = "record_id"
const RECORD_ID_EXPIRATION = 2 * time.Hour

type TrackEntryHandler struct {
}

type TrackEntry struct {
	CreatedAt int    `json:"created_at"`
	Markup    string `json:"markup"`
}

func (h *TrackEntryHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
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

	recordCookie := createRecordCookie(response, request)
	recordId := recordCookie.Value

	decoder := json.NewDecoder(request.Body)
	var trackEntry TrackEntry
	err := decoder.Decode(&trackEntry)
	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(response, `{"message": "Invalid JSON payload"}`)
		return
	}

	htmlBytes, err := base64.StdEncoding.DecodeString(trackEntry.Markup)
	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(response, `{"message": "Invalid base64 payload"}`)
		return
	}

	trackEntriesPath := filepath.Join("/tmp", "track_entries", recordId)
	os.MkdirAll(trackEntriesPath, os.ModePerm)
	fileName := filepath.Join(trackEntriesPath, fmt.Sprintf("%d.html", trackEntry.CreatedAt))
	err = ioutil.WriteFile(fileName, htmlBytes, 0644)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, `{"message": "Fail to save request"}`)
		return
	}

	response.WriteHeader(http.StatusCreated)
	fmt.Printf("Tracking dom, record_id: %s, created_at: %d\n", recordId, trackEntry.CreatedAt)
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
