package bgh

import (
	"fmt"
	"net/url"
)

var InvalidURLError = fmt.Errorf("URL does not contain all required query parameters")

func PathFromURL(u string) (string, error) {
	url, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	query := url.Query()

	court := query.Get("Gericht")
	date := query.Get("Datum")
	nr := query.Get("nr")
	anz := query.Get("anz")
	pos := query.Get("pos")

	if court == "" || date == "" || nr == "" || anz == "" || pos == "" {
		return "", InvalidURLError
	}

	return fmt.Sprintf("judgements/%s/%s/%s_%s_%s.pdf", court, date, nr, anz, pos), nil
}
