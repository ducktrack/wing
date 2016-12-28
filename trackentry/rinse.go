package trackentry

import (
	"github.com/PuerkitoBio/goquery"
	"bytes"
	"errors"
)

func (trackEntry *TrackEntry) Rinse() (secureMarkup string, error error) {
	htmlBytes, err := trackEntry.MarkupBytes()
	if err != nil {
		return "", errors.New("Failed to get markup")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBytes))
	if err != nil {
		return "", errors.New("Failed to parse HTML")
	}

	scripts := doc.Find("script")
	scripts.Each(func(i int, s *goquery.Selection) { s.Remove() })

	content, err := doc.Html()
	if err != nil {
		return "", errors.New("Failed to generate secure HTML")
	}

	return content, nil
}

