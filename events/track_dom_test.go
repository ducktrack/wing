package events

import (
	helpers "github.com/duckclick/wing/testing"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_ToJSON(t *testing.T) {
	trackDOM, err := TrackDOMFromJSON(Event{
		CreatedAt:  1487696788863,
		URL:        "http://example.org/some/path",
		RawPayload: helpers.CreateRawMessage(`{"markup": %s}`, helpers.Base64BlankMarkup),
	})
	assert.Nil(t, err, "TrackDOMFromJSON() should succeed")

	json, err := trackDOM.ToJSON()
	assert.Nil(t, err, "ToJSON() should succeed")
	assert.Equal(
		t,
		`{"created_at":1487696788863,"url":"http://example.org/some/path","type":"TrackDOM","payload":{"markup":"<html><head></head><body></body></html>"}}`,
		strings.TrimSpace(string(json)),
		"generates a valid JSON",
	)
}

func Test_TrackDOMFromJSON_StripScriptTags(t *testing.T) {
	htmlSample := `<html><head><script src="evil"></script></head><body><script src="g-evil"></script></body></html>`
	trackDOM, err := TrackDOMFromJSON(Event{
		CreatedAt:  123456,
		RawPayload: helpers.CreateRawMessage(`{"markup": %s}`, helpers.ToBase64(htmlSample)),
	})
	assert.Nil(t, err, "TrackDOMFromJSON() should succeed")
	assert.Equal(t, "<html><head></head><body></body></html>", string(trackDOM.Payload.Markup), "script tags should be removed")
}

func Test_TrackDOMFromJSON_DecodesBase64Markup(t *testing.T) {
	trackDOM, err := TrackDOMFromJSON(Event{
		CreatedAt:  123456,
		RawPayload: helpers.CreateRawMessage(`{"markup": %s}`, helpers.Base64BlankMarkup),
	})
	assert.Nil(t, err, "TrackDOMFromJSON() should succeed")
	assert.Equal(t, "<html><head></head><body></body></html>", string(trackDOM.Payload.Markup), "base64 should be decoded")
}
