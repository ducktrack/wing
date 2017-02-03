package trackentry

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeBase64(t *testing.T) {
	htmlSample := `<html><head></head><body></body></html>`
	entry := &TrackEntry{
		CreatedAt: 123456,
		Markup:    base64.StdEncoding.EncodeToString([]byte(htmlSample)),
	}

	bytes, err := entry.MarkupBytes()
	assert.Nil(t, err, "MarkupBytes() should succeed")
	assert.Equal(t, "<html><head></head><body></body></html>", string(bytes), "base64 should be decoded")
}
