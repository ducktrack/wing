package events

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

// TrackDOM definition
type TrackDOM struct {
	Event
	Payload TrackDOMPayload `json:"payload"`
}

// TrackDOMPayload definition
type TrackDOMPayload struct {
	Markup string `json:"markup"`
}

// NewTrackDOM creates a new TrackDOM event
func NewTrackDOM(event Event, payload TrackDOMPayload) *TrackDOM {
	return &TrackDOM{
		Event{
			CreatedAt: event.CreatedAt,
			URL:       event.URL,
			Type:      "TrackDOM",
		},
		payload,
	}
}

// TrackDOMFromJSON creates a new TrackDOM based on event.RawPayload. It decodes
// the base64 and generates a secure markup
func TrackDOMFromJSON(event Event) (*TrackDOM, error) {
	var payload TrackDOMPayload
	err := json.Unmarshal([]byte(event.RawPayload), &payload)

	if err != nil {
		return &TrackDOM{}, errors.Wrap(err, "Failed to decode payload")
	}

	secureMarkup, err := rinse(payload.Markup)

	if err != nil {
		return &TrackDOM{}, errors.Wrap(err, "Failed to generate a secure markup")
	}

	payload.Markup = secureMarkup
	return NewTrackDOM(event, payload), err
}

// GetEvent returns the base event
func (trackDOM *TrackDOM) GetEvent() Event {
	return trackDOM.Event
}

// ToJSON serialize to a JSON string
func (trackDOM *TrackDOM) ToJSON() (string, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(trackDOM)

	if err != nil {
		return "", errors.Wrap(err, "Failed to encode TrackDOM to JSON")
	}

	return string(buffer.Bytes()), err
}

// Rinse generates a secure markup (no script tags)
func rinse(base64Markup string) (string, error) {
	htmlBytes, err := base64MarkupToBytes(base64Markup)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get the markup")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBytes))
	if err != nil {
		return "", errors.Wrap(err, "Failed to parse HTML")
	}

	scripts := doc.Find("script")
	scripts.Each(func(i int, s *goquery.Selection) { s.Remove() })

	content, err := doc.Html()
	return content, errors.Wrap(err, "Failed to generate secure HTML")
}

func base64MarkupToBytes(markup string) ([]byte, error) {
	bytes, err := base64.StdEncoding.DecodeString(markup)
	return bytes, errors.Wrap(err, "Invalid base64 payload")
}
