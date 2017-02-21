package trackentry

import (
	helpers "github.com/duckclick/wing/testing"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestMarkupBytes(t *testing.T) {
	entry := &TrackEntry{
		Markup: helpers.ToBase64("<html><head></head><body></body></html>"),
	}

	bytes, err := entry.MarkupBytes()
	assert.Nil(t, err, "MarkupBytes() should succeed")
	assert.Equal(t, "<html><head></head><body></body></html>", string(bytes), "base64 should be decoded")
}

func TestToJSON(t *testing.T) {
	entry := &TrackEntry{
		CreatedAt: 1487696788863,
		URL:       "http://example.org/some/path",
		Markup:    helpers.ToBase64("<html><head></head><body></body></html>"),
	}

	json, err := entry.ToJSON()
	assert.Nil(t, err, "ToJSON() should succeed")
	assert.Equal(
		t,
		`{"created_at":1487696788863,"url":"http://example.org/some/path","markup":"<html><head></head><body></body></html>"}`,
		strings.TrimSpace(string(json)),
		"generates a valid JSON",
	)
}
