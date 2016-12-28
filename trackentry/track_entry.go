package trackentry

import (
	"encoding/base64"
	"errors"
)

type TrackEntry struct {
	CreatedAt int    `json:"created_at"`
	Origin    string `json:"origin"`
	Markup    string `json:"markup"`
}

func (trackEntry *TrackEntry) MarkupBytes() ([]byte, error) {
	bytes, err := base64.StdEncoding.DecodeString(trackEntry.Markup)
	if err != nil {
		return nil, errors.New("Invalid base64 payload")
	}

	return bytes, nil
}