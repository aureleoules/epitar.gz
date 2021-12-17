package api

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"syscall"

	"github.com/aureleoules/epitar/config"
	"github.com/aureleoules/epitar/db"
	"github.com/aureleoules/epitar/models"
	"github.com/brpaz/echozap"
	"github.com/expectedsh/go-sonic/sonic"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

var searcher sonic.Searchable

func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nExiting...")

		os.Exit(1)
	}()
}

func Serve() {
	zap.S().Info("Starting server...")

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(echozap.ZapLogger(zap.L()))

	db.Init()
	defer db.Close()

	initSonic()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "salu")
	})

	var lock sync.Mutex
	e.GET("/search", func(c echo.Context) error {
		query := c.QueryParam("q")
		if len(query) == 0 {
			return c.JSON(http.StatusNotAcceptable, nil)
		}

		// TODO: Create pool of searchers
		lock.Lock()
		resp, err := searcher.Query("files", "default", query, 50, 0, "")
		lock.Unlock()
		if err != nil {
			c.NoContent(http.StatusInternalServerError)
			return err
		}

		var ids []string
		for _, r := range resp {
			keyStrSplit := strings.Split(r, ":")
			if len(keyStrSplit) != 2 {
				continue
			}
			keyStr := keyStrSplit[1]
			ids = append(ids, keyStr)
		}

		if len(ids) == 0 {
			return c.NoContent(http.StatusNotFound)
		}

		files, err := models.GetFilesMeta(ids)
		if err != nil {
			c.NoContent(http.StatusInternalServerError)
			return err
		}

		return c.JSON(200, files)
	})

	e.GET("/file/:id", func(c echo.Context) error {
		id := c.Param("id")
		meta, err := models.GetFileMeta(id)
		if err != nil {
			c.NoContent(http.StatusInternalServerError)
			return err
		}

		c.Response().Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", meta.Name))
		return c.File(path.Join(config.Cfg.Index.Store, "files", id))
	})

	e.Start(":1323")
}
