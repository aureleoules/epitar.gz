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
	"strings"
	"unicode"

	"code.sajari.com/docconv/client"
	"github.com/aureleoules/epitar/config"
	"github.com/aureleoules/epitar/db"
	"github.com/aureleoules/epitar/models"
	"github.com/expectedsh/go-sonic/sonic"
	"github.com/fatih/color"
	"go.uber.org/zap"
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

func isNewOrigin(origins []models.FileOrigin, name string) bool {
	for _, o := range origins {
		if o.Module == name {
			return false
		}
	}
	return true
}

func (m *Module) indexPath(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	originURL, err := ioutil.ReadFile(path + ".url")
	if err != nil {
		return err
	}

	url := strings.TrimSpace(string(originURL))

	h := sha1.Sum(data)
	key := hex.EncodeToString(h[:])

	fileMeta, err := models.GetFileMeta(key)
	if err == nil {
		fmt.Println(fileMeta.Origins)
		if isNewOrigin(fileMeta.Origins, m.Meta.Slug) {
			zap.S().Debugf("Adding origin %s to file %s", m.Meta.Slug, key)
			err = fileMeta.AddOrigin(m.Meta.Slug, url)
			if err != nil {
				return err
			}

			return nil
		} else {
			color.Yellow("Skipped %s because it already exists", filepath.Base(path))
			return nil
		}
	}

	name := filepath.Base(path)
	zap.S().Debugf("Pre processing %s", name)
	keywords := preprocessText(name, data)

	maxLength := len(keywords)
	if maxLength > 100 {
		maxLength = 100
	}

	meta := models.FileMeta{
		ID:      hex.EncodeToString(h[:]),
		Name:    name,
		Size:    int64(len(data)),
		Summary: keywords[:maxLength],
	}

	if err := meta.Save(); err != nil {
		return err
	}

	if err := meta.AddOrigin(m.Meta.Slug, url); err != nil {
		return err
	}

	zap.S().Infof("Writing file... %s", name)
	err = os.WriteFile(filepath.Join(config.Cfg.Index.Store, "files", key), data, 0644)
	if err != nil {
		return err
	}

	keywords = name + " " + keywords
	zap.S().Infof("Pushing keyworks... %s", name)
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
	p := filepath.Join(config.Cfg.Output, m.Meta.Slug)
	files := listIndexableFiles(p, m.IndexOptions.Files)

	l := len(files)
	i := 0

	for _, f := range files {
		fmt.Printf("\rIndexing %s (%d/%d)\n", f, i, l)
		err := m.indexPath(f)
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

	ingester, err = sonic.NewIngester(config.Cfg.Index.Sonic.Host, config.Cfg.Index.Sonic.Port, config.Cfg.Index.Sonic.Password)
	if err != nil {
		panic(err)
	}

	for _, m := range modules {
		if m.IndexOptions.Enable {
			color.Magenta("Indexing module %s", m.Meta.Name)
			if err := m.Index(); err != nil {
				return err
			}
		}
	}

	return nil
}
