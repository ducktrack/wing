package main

import (
	"fmt"
	"github.com/ducktrack/wing/handlers"
	"net/http"
)

func main() {
	fmt.Printf("Starging Wing at port 7273\n")
	http.Handle("/", &handlers.TrackEntryHandler{})
	http.ListenAndServe(":7273", nil)
}
