package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TrackEntry struct {
	CreatedAt int    `json:"created_at"`
	Markup    string `json:"markup"`
}

func handler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	responseWriter.Header().Set("Access-Control-Allow-Headers", "content-type")

	if request.Method == "OPTIONS" {
		return
	}

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

	fileName := fmt.Sprintf("./track-entries/%d.html", trackEntry.CreatedAt)
	err = ioutil.WriteFile(fileName, htmlBytes, 0644)
	if err != nil {
		fmt.Fprintf(responseWriter, "quack! error writing HTML")
		panic(err)
	}

	fmt.Printf("%+v\n", trackEntry)
	fmt.Fprintf(responseWriter, "working")
}

func main() {
	fmt.Printf("Starging Wing at port 7273\n")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":7273", nil)
}
