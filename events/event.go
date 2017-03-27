package events

import (
	"encoding/json"
	"github.com/pkg/errors"
)

// Event definition
type Event struct {
	CreatedAt  int             `json:"created_at"`
	URL        string          `json:"url"`
	Type       string          `json:"type"`
	RawPayload json.RawMessage `json:"payload"`
}

// Trackable interface
type Trackable interface {
	ToJSON() (string, error)
	GetEvent() Event
}

// DecodeJSON decodes json to a list of trackable types
func DecodeJSON(bytes []byte) ([]Trackable, error) {
	var entries []Event
	var result []Trackable
	err := json.Unmarshal(bytes, &entries)

	if err != nil {
		return result, errors.Wrap(err, "Failed to decode JSON")
	}

	result = make([]Trackable, len(entries))

	for i, entry := range entries {
		switch entry.Type {
		case "TrackDOM":
			result[i], err = TrackDOMFromJSON(entry)

		default:
			return result, errors.Errorf("Invalid event type: %s", entry.Type)
		}

		if err != nil {
			return result, errors.Wrapf(err, "Failed to decode entry of type %s", entry.Type)
		}
	}

	return result, err
}
