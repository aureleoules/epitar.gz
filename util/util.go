package util

import (
	"regexp"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	transformChain = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	preprocRgx     = regexp.MustCompile(`[^\w]`)
)

func NormalizeText(txt string) string {
	result, _, _ := transform.String(transformChain, txt)
	return preprocRgx.ReplaceAllString(result, " ")
}
