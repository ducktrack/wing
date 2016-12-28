package trackentry

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"encoding/base64"
)

func TestStripScriptTags(t *testing.T) {
	htmlSample := `<html><head><script src="evil"></script></head><body><script src="g-evil"></script></body></html>`
	entry := &TrackEntry{
		CreatedAt: 123456,
		Markup:    base64.StdEncoding.EncodeToString([]byte(htmlSample)),
	}

	rinsedMarkup, err := entry.Rinse()
	assert.Nil(t, err, "rinse should succeed")
	assert.Equal(t, "<html><head></head><body></body></html>", string(rinsedMarkup), "script tags should be removed")
}

