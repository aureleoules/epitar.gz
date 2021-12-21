package archive

import (
	"bytes"
	"log"
	"regexp"
	"strings"
	"unicode"

	"code.sajari.com/docconv/client"
	"github.com/expectedsh/go-sonic/sonic"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	transformChain = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	preprocRgx     = regexp.MustCompile(`[^\w]`)

	docconvClient *client.Client
	ingester      sonic.Ingestable
)

func normalize(txt string) string {
	result, _, _ := transform.String(transformChain, txt)
	return preprocRgx.ReplaceAllString(result, " ")
}

func preprocessNews(data string) string {
	content := strings.Split(string(data), "\n\n")
	if len(content) < 2 {
		return ""
	}

	return normalize(strings.Join(content[1:], " "))
}

func preprocessDocument(filename string, data []byte) string {
	res, err := docconvClient.Convert(bytes.NewReader(data), filename)
	if err != nil {
		log.Fatal(err)
	}

	return normalize(res.Body)
}
