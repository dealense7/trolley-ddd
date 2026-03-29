package embedder

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func NormalizeName(raw string) string {
	// 1. Unicode normalization (NFD → strip diacritics → NFC)
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	s, _, _ := transform.String(t, raw)

	// 2. Lowercase
	s = strings.ToLower(s)

	// 3. Remove punctuation except hyphens and spaces
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '-' {
			b.WriteRune(r)
		} else {
			b.WriteRune(' ')
		}
	}
	s = b.String()

	// 4. Collapse multiple spaces
	s = strings.Join(strings.Fields(s), " ")

	return strings.TrimSpace(s)
}
