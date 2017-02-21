package trackentry

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/pkg/errors"
)

// TrackEntry definition
type TrackEntry struct {
	CreatedAt int    `json:"created_at"`
	URL       string `json:"url"`
	Markup    string `json:"markup"`
}

// MarkupBytes decodes the base64 markup
func (trackEntry *TrackEntry) MarkupBytes() ([]byte, error) {
	bytes, err := base64.StdEncoding.DecodeString(trackEntry.Markup)
	return bytes, errors.Wrap(err, "Invalid base64 payload")
}

// ToJSON serialize to a JSON string
func (trackEntry *TrackEntry) ToJSON() (string, error) {
	markup, err := trackEntry.Rinse()

	if err != nil {
		return "", errors.Wrap(err, "Failed to generate a secure markup")
	}

	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)

	err = encoder.Encode(TrackEntry{
		CreatedAt: trackEntry.CreatedAt,
		URL:       trackEntry.URL,
		Markup:    markup,
	})

	if err != nil {
		return "", errors.Wrap(err, "Failed to encode JSON")
	}

	return string(buffer.Bytes()), err
}
