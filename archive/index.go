package archive

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

	"code.sajari.com/docconv/client"
	"github.com/araddon/dateparse"
	"github.com/aureleoules/epitar/config"
	"github.com/aureleoules/epitar/db"
	"github.com/aureleoules/epitar/models"
	"github.com/aureleoules/epitar/util"
	"github.com/expectedsh/go-sonic/sonic"
	"github.com/fatih/color"
	"go.uber.org/zap"
)

func isNewOrigin(origins []models.FileOrigin, name string) bool {
	for _, o := range origins {
		if o.Module == name {
			return false
		}
	}
	return true
}

func (m *Module) indexDocument(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	key := getFileID(data)

	urlBytes, err := ioutil.ReadFile(path + ".url")
	if err != nil {
		return err
	}

	url := strings.TrimSpace(string(urlBytes))

	fileMeta, err := models.GetFileMeta(key)
	if err == nil {
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
	keywords := preprocessDocument(name, data)

	maxLength := len(keywords)
	if maxLength > 100 {
		maxLength = 100
	}

	meta := models.FileMeta{
		ID:      key,
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

	keywords = util.NormalizeText(name) + " " + keywords
	zap.S().Infof("Pushing keywords... %s", name)
	err = ingester.Push("files", "all", "key:"+key, keywords, "")
	if err != nil {
		return err
	}

	err = ingester.Push("files", m.Meta.Slug, "key:"+key, keywords, "")
	return err
}

func (m *Module) indexNews(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	key := getFileID(data)

	_, err = models.GetNews(key)
	if err == nil {
		color.Yellow("Skipped %s because it already exists", filepath.Base(path))
		return nil
	}

	println(key)

	reader := bufio.NewReader(strings.NewReader(string(data) + "\r\n"))
	tp := textproto.NewReader(reader)

	header, err := tp.ReadMIMEHeader()
	if err != nil {
		return err
	}

	// http.Header and textproto.MIMEHeader are both just a map[string][]string
	httpHeader := http.Header(header)

	date, err := dateparse.ParseAny(httpHeader["Date"][0])
	if err != nil {
		return err
	}

	content := preprocessNews(string(data))
	summary := content
	if len(content) > 100 {
		summary = content[:100]
	}

	dec := new(mime.WordDecoder)
	subject, err := dec.DecodeHeader(httpHeader.Get("Subject"))
	if err != nil {
		return err
	}

	news := models.News{
		ID:         key,
		Subject:    subject,
		From:       httpHeader["From"][0],
		Date:       date,
		Newsgroups: strings.Join(httpHeader["Newsgroups"], ","),
		MessageID:  httpHeader["Message-Id"][0],
		Size:       len(data),
		Summary:    summary,
	}

	zap.S().Infof("Writing file... %s", filepath.Base(path))
	if err := news.Save(); err != nil {
		return err
	}

	content = news.ID + " " + content
	zap.S().Infof("Pushing keywords... %s", filepath.Base(path))
	err = ingester.Push("news", "all", "key:"+key, content, "")
	if err != nil {
		return err
	}

	for _, ng := range httpHeader["Newsgroups"] {
		err = ingester.Push("news", ng, "key:"+key, content, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Module) indexFile(path string) error {
	fmt.Println("|" + getExtension(path) + "|")
	switch getExtension(path) {
	case "news":
		return m.indexNews(path)

	case "pdf", "docx", "doc", "pptx", "ppt", "odt", "odp":
		return m.indexDocument(path)
	}

	zap.S().Panicf("Unsupported file type: %s\n", path)

	return nil
}

func (m *Module) Index() error {
	p := filepath.Join(config.Cfg.Output, m.Meta.Slug)
	files := listIndexableFiles(p, m.IndexOptions.Files)

	total := len(files)
	i := 0

	for _, f := range files {
		fmt.Printf("\rIndexing %s (%d/%d)\n", f, i, total)
		err := m.indexFile(f)
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
