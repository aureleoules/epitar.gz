package archive

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"unicode"

	"code.sajari.com/docconv/client"
	"github.com/aureleoules/epitar/config"
	"github.com/aureleoules/epitar/db"
	"github.com/aureleoules/epitar/models"
	"github.com/expectedsh/go-sonic/sonic"
	"github.com/fatih/color"
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

func preprocessText(filename string, data []byte) string {
	res, err := docconvClient.Convert(bytes.NewReader(data), filename)
	if err != nil {
		log.Fatal(err)
	}

	result, _, _ := transform.String(transformChain, res.Body)
	return preprocRgx.ReplaceAllString(result, " ")
}

func checkExtension(accept []string, path string) bool {
	ext := filepath.Ext(path)
	if ext == "" {
		return false
	}
	ext = ext[1:]
	for _, a := range accept {
		if a == ext {
			return true
		}
	}
	return false
}

func buildKey(keyType string, hash []byte) []byte {
	return append([]byte(keyType), hash...)
}

func indexPath(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	h := sha1.Sum(data)
	key := hex.EncodeToString(h[:])

	_, err = models.GetFileMeta(key)
	if err == nil {
		color.Yellow("Skipped %s because it already exists", filepath.Base(path))
		return nil
	}

	name := filepath.Base(path)
	keywords := preprocessText(name, data)
	if keywords == "" {
		color.Yellow("Skipped %s because it has no keywords", filepath.Base(path))
		return nil
	}

	meta := models.FileMeta{
		ID:      hex.EncodeToString(h[:]),
		Name:    name,
		Size:    int64(len(data)),
		Summary: keywords[:100],
	}

	if err := meta.Save(); err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(config.Cfg.Index.Store, "files", key), data, 0644)
	if err != nil {
		return err
	}

	err = ingester.Push("files", "default", "key:"+key, keywords, "")
	if err != nil {
		return err
	}

	return nil
}

func listIndexableFiles(rootDir string, accept []string) []string {
	var files []string
	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if !checkExtension(accept, path) {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files
}

func (m *Module) Index() error {
	p := filepath.Join(config.Cfg.Output, m.Name)
	files := listIndexableFiles(p, m.IndexOptions.Files)

	l := len(files)
	i := 0

	for _, f := range files {
		fmt.Printf("\rIndexing %s (%d/%d)\n", f, i, l)
		err := indexPath(f)
		if err != nil {
			return err
		}

		i++
	}

	return nil
}

func IndexModules() error {
	var err error
	color.Green("Indexing modules")

	docconvClient = client.New()

	db.Init()
	defer db.Close()

	ingester, err = sonic.NewIngester("localhost", 1491, "password")
	if err != nil {
		panic(err)
	}

	for _, m := range modules {
		color.Magenta("Indexing module %s", m.Name)
		if err := m.Index(); err != nil {
			return err
		}
	}

	return nil
}
