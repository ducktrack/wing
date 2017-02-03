package trackentry

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

// Rinse generates a secure markup (no script tags)
func (trackEntry *TrackEntry) Rinse() (string, error) {
	htmlBytes, err := trackEntry.MarkupBytes()
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
