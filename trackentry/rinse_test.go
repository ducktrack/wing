package trackentry

import (
	helpers "github.com/duckclick/wing/testing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStripScriptTags(t *testing.T) {
	htmlSample := `<html><head><script src="evil"></script></head><body><script src="g-evil"></script></body></html>`
	entry := &TrackEntry{
		CreatedAt: 123456,
		Markup:    helpers.ToBase64(htmlSample),
	}

	rinsedMarkup, err := entry.Rinse()
	assert.Nil(t, err, "Rinse() should succeed")
	assert.Equal(t, "<html><head></head><body></body></html>", string(rinsedMarkup), "script tags should be removed")
}
