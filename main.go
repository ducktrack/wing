package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type TrackEntry struct {
	CreatedAt int    `json:"created_at"`
	Markup    string `json:"markup"`
}

func handler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
	responseWriter.Header().Set("Access-Control-Allow-Headers", "content-type")
	responseWriter.Header().Set("Access-Control-Allow-Credentials", "true")

	if request.Method != "POST" {
		fmt.Fprintf(responseWriter, "nothing to do, but I'm working")
		return
	}

	recordCookie, noCookieErr := request.Cookie("record_id")
	if noCookieErr != nil || recordCookie.Value == "" {
		uuid, err := newUUID()
		if err != nil {
			fmt.Fprintf(responseWriter, "quack! error generating uuid")
			panic(err)
		}
		expiration := time.Now().Add(2 * time.Hour)
		recordCookie = &http.Cookie{
			Name:     "record_id",
			Value:    string(uuid),
			Expires:  expiration,
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(responseWriter, recordCookie)
	}
	recordId := recordCookie.Value

	decoder := json.NewDecoder(request.Body)
	var trackEntry TrackEntry
	err := decoder.Decode(&trackEntry)
	if err != nil {
		fmt.Fprintf(responseWriter, "quack! error decoding JSON")
		panic(err)
	}

	htmlBytes, err := base64.StdEncoding.DecodeString(trackEntry.Markup)
	if err != nil {
		fmt.Fprintf(responseWriter, "quack! error decoding base64")
		panic(err)
	}

	trackEntriesPath := filepath.Join("/tmp", "track_entries", recordId)
	os.MkdirAll(trackEntriesPath, os.ModePerm)
	fileName := filepath.Join(trackEntriesPath, fmt.Sprintf("%d.html", trackEntry.CreatedAt))
	err = ioutil.WriteFile(fileName, htmlBytes, 0644)
	if err != nil {
		fmt.Fprintf(responseWriter, "quack! error writing HTML")
		panic(err)
	}

	fmt.Printf("Tracking dom, record_id: %s, created_at: %d\n", recordId, trackEntry.CreatedAt)
	fmt.Fprintf(responseWriter, "record id: %s", recordId)
}

func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func main() {
	fmt.Printf("Starging Wing at port 7273\n")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":7273", nil)
}
