package events_test

import (
	"github.com/duckclick/wing/events"
	helpers "github.com/duckclick/wing/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type TrackDOMTestSuite struct {
	suite.Suite
}

func (suite *TrackDOMTestSuite) Test_ToJSON() {
	trackDOM, err := events.TrackDOMFromJSON(events.Event{
		CreatedAt:  1487696788863,
		URL:        "http://example.org/some/path",
		RawPayload: helpers.CreateRawMessage(`{"markup": %s}`, helpers.Base64BlankMarkup),
	})
	assert.Nil(suite.T(), err, "TrackDOMFromJSON() should succeed")

	json, err := trackDOM.ToJSON()
	assert.Nil(suite.T(), err, "ToJSON() should succeed")
	assert.Equal(
		suite.T(),
		`{"created_at":1487696788863,"url":"http://example.org/some/path","type":"TrackDOM","payload":{"markup":"<html><head></head><body></body></html>"}}`,
		strings.TrimSpace(string(json)),
		"generates a valid JSON",
	)
}

func (suite *TrackDOMTestSuite) Test_TrackDOMFromJSON_StripScriptTags() {
	htmlSample := `<html><head><script src="evil"></script></head><body><script src="g-evil"></script></body></html>`
	trackDOM, err := events.TrackDOMFromJSON(events.Event{
		CreatedAt:  123456,
		RawPayload: helpers.CreateRawMessage(`{"markup": %s}`, helpers.ToBase64(htmlSample)),
	})
	assert.Nil(suite.T(), err, "TrackDOMFromJSON() should succeed")
	assert.Equal(suite.T(), "<html><head></head><body></body></html>", string(trackDOM.Payload.Markup), "script tags should be removed")
}

func (suite *TrackDOMTestSuite) Test_TrackDOMFromJSON_DecodesBase64Markup() {
	trackDOM, err := events.TrackDOMFromJSON(events.Event{
		CreatedAt:  123456,
		RawPayload: helpers.CreateRawMessage(`{"markup": %s}`, helpers.Base64BlankMarkup),
	})
	assert.Nil(suite.T(), err, "TrackDOMFromJSON() should succeed")
	assert.Equal(suite.T(), "<html><head></head><body></body></html>", string(trackDOM.Payload.Markup), "base64 should be decoded")
}

func TestTrackDOM(t *testing.T) {
	suite.Run(t, new(TrackDOMTestSuite))
}
