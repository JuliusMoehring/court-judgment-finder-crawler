package bgh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PathFromURL(t *testing.T) {
	t.Run("Creates the correct file path from URL", func(t *testing.T) {
		expected := "judgements/bgh/2021/117424_3571_2950.pdf"
		actual, err := PathFromURL("https://juris.bundesgerichtshof.de/cgi-bin/rechtsprechung/document.py?Gericht=bgh&Art=en&Datum=2021&Seite=98&nr=117424&anz=3571&pos=2950&Blank=1.pdf")

		assert.NoError(t, err, "Should not return an error")
		assert.Equal(t, expected, actual, "Should return the correct file path")
	})

	t.Run("Returns error if the URL search params does not include 'Gericht'", func(t *testing.T) {
		_, err := PathFromURL("https://juris.bundesgerichtshof.de/cgi-bin/rechtsprechung/document.py?Art=en&Datum=2021&Seite=98&nr=117424&anz=3571&pos=2950&Blank=1.pdf")

		assert.ErrorIs(t, err, InvalidURLError, "Should return an `InvalidURLError` error")
	})

	t.Run("Returns error if the URL search params does not include 'Datum'", func(t *testing.T) {
		_, err := PathFromURL("https://juris.bundesgerichtshof.de/cgi-bin/rechtsprechung/document.py?Gericht=bgh&Art=en&Seite=98&nr=117424&anz=3571&pos=2950&Blank=1.pdf")

		assert.ErrorIs(t, err, InvalidURLError, "Should return an `InvalidURLError` error")
	})

	t.Run("Returns error if the URL search params does not include 'nr'", func(t *testing.T) {
		_, err := PathFromURL("https://juris.bundesgerichtshof.de/cgi-bin/rechtsprechung/document.py?Gericht=bgh&Art=en&Datum=2021&Seite=98&anz=3571&pos=2950&Blank=1.pdf")

		assert.ErrorIs(t, err, InvalidURLError, "Should return an `InvalidURLError` error")
	})

	t.Run("Returns error if the URL search params does not include 'anz'", func(t *testing.T) {
		_, err := PathFromURL("https://juris.bundesgerichtshof.de/cgi-bin/rechtsprechung/document.py?Gericht=bgh&Art=en&Datum=2021&Seite=98&nr=117424&pos=2950&Blank=1.pdf")

		assert.ErrorIs(t, err, InvalidURLError, "Should return an `InvalidURLError` error")
	})

	t.Run("Returns error if the URL search params does not include 'pos'", func(t *testing.T) {
		_, err := PathFromURL("https://juris.bundesgerichtshof.de/cgi-bin/rechtsprechung/document.py?Gericht=bgh&Art=en&Datum=2021&Seite=98&nr=117424&anz=3571&Blank=1.pdf")

		assert.ErrorIs(t, err, InvalidURLError, "Should return an `InvalidURLError` error")
	})
}
