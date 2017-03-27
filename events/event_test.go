package events_test

import (
	"fmt"
	"github.com/duckclick/wing/events"
	helpers "github.com/duckclick/wing/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type EventTestSuite struct {
	suite.Suite
}

func (suite *EventTestSuite) Test_DecodeJSON_TrackDOM() {
	jsonPayload := fmt.Sprintf(
		`[{"created_at":1487696788863,"url":"http://example.org/some/path","type":"TrackDOM","payload":{"markup":"%s"}}]`,
		helpers.Base64BlankMarkup,
	)

	trackables, err := events.DecodeJSON([]byte(jsonPayload))
	assert.Nil(suite.T(), err, "events#DecodeJSON() should succeed")
	assert.Equal(suite.T(), len(trackables), 1, "events#DecodeJSON() should decode 1 Trackable")

	trackDOM := trackables[0]
	assert.Equal(suite.T(), fmt.Sprintf("%T", trackDOM), "*events.TrackDOM", "trackable should be a TrackDOM")
}

func (suite *EventTestSuite) Test_DecodeJSON_TrackDOM_WhenItFailsToDecode() {
	jsonPayload := `[{"created_at":1487696788863,"url":"http://example.org/some/path","type":"TrackDOM","payload":{"markup":"invalid"}}]`
	_, err := events.DecodeJSON([]byte(jsonPayload))
	assert.NotNil(suite.T(), err, "events#DecodeJSON() should fail because the event is not valid")
}

func (suite *EventTestSuite) Test_DecodeJSON_UnknownEvent() {
	jsonPayload := `[{"created_at":1487696788863,"url":"http://example.org/some/path"}]`
	_, err := events.DecodeJSON([]byte(jsonPayload))
	assert.NotNil(suite.T(), err, "events#DecodeJSON() should fail because the event type is unknown")
}

func (suite *EventTestSuite) Test_DecodeJSON_WithInvalidJSON() {
	jsonPayload := `invalid-json`
	_, err := events.DecodeJSON([]byte(jsonPayload))
	assert.NotNil(suite.T(), err, "events#DecodeJSON() should fail because the JSON payload is invalid")
}

func TestEvent(t *testing.T) {
	suite.Run(t, new(EventTestSuite))
}
