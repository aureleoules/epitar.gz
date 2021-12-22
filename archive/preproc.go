package archive

import (
	"bytes"
	"log"
	"strings"

	"code.sajari.com/docconv/client"
	"github.com/aureleoules/epitar/util"
	"github.com/expectedsh/go-sonic/sonic"
)

var (
	docconvClient *client.Client
	ingester      sonic.Ingestable
)

func preprocessNews(data string) string {
	content := strings.Split(string(data), "\n\n")
	if len(content) < 2 {
		return ""
	}

	return util.NormalizeText(strings.Join(content[1:], " "))
}

func preprocessDocument(filename string, data []byte) string {
	res, err := docconvClient.Convert(bytes.NewReader(data), filename)
	if err != nil {
		log.Fatal(err)
	}

	return util.NormalizeText(res.Body)
}
