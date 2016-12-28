package trackentry

import (
	"github.com/PuerkitoBio/goquery"
	"bytes"
	"github.com/pkg/errors"
)

func (trackEntry *TrackEntry) Rinse() (secureMarkup string, error error) {
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

