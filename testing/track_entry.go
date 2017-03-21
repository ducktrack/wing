package testing

import (
	"encoding/json"
	"fmt"
)

// BlankMarkup
var BlankMarkup = "<html><head></head><body></body></html>"

// Base64BlankMarkup
var Base64BlankMarkup = ToBase64(BlankMarkup)

// CreateRawMessage a new json.RawMessage
func CreateRawMessage(rawJSON string, args ...interface{}) json.RawMessage {
	if len(args) > 0 {
		return json.RawMessage(fmt.Sprintf(rawJSON, args))
	}

	return json.RawMessage(rawJSON)
}
