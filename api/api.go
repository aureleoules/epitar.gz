package api

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
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

func handleInterrupt() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nExiting...")

		os.Exit(1)
	}()
}

func Serve() {
	handleInterrupt()
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
		const limit = 150
		query := c.QueryParam("q")
		if len(query) == 0 {
			return c.JSON(http.StatusNotAcceptable, nil)
		}

		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			page = 1
		}

		// TODO: Create pool of searchers
		lock.Lock()
		resp, err := searcher.Query("files", "default", query, limit, (page-1)*limit, "")
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

		var filtered []models.FileMeta
		if c.QueryParam("module") == "" {
			filtered = files
		} else {
			for _, f := range files {
				for _, o := range f.Origins {
					if o.Module == c.QueryParam("module") {
						filtered = append(filtered, f)
						break
					}
				}
			}
		}

		return c.JSON(200, filtered)
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

	e.GET("/modules", func(c echo.Context) error {
		var modules []config.ModuleMeta

		for _, m := range config.Cfg.Modules {
			if m.Enable {
				modules = append(modules, m.Meta)
			}
		}

		return c.JSON(http.StatusOK, modules)
	})

	e.GET("/stats", func(c echo.Context) error {
		// TODO: CACHE
		stats, err := models.GetStats()
		if err != nil {
			c.NoContent(http.StatusInternalServerError)
			return err
		}

		return c.JSON(200, stats)
	})

	e.Start(":1323")
}
