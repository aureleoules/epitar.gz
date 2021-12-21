package api

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/aureleoules/epitar/config"
	"github.com/aureleoules/epitar/db"
	"github.com/aureleoules/epitar/models"
	"github.com/brpaz/echozap"
	"github.com/expectedsh/go-sonic/sonic"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

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

	e.GET("/documents/search", func(c echo.Context) error {
		const limit = 150
		query := c.QueryParam("q")
		if len(query) == 0 {
			return c.JSON(http.StatusNotAcceptable, nil)
		}

		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			page = 1
		}

		searcher, err := sonicPool.Get()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return err
		}

		module := strings.ToLower(c.QueryParam("module"))
		if module == "" {
			module = "all"
		}
		resp, err := searcher.(sonic.Searchable).Query("files", module, query, limit, (page-1)*limit, "")
		if err != nil {
			c.NoContent(http.StatusInternalServerError)
			return err
		}

		sonicPool.Put(searcher)

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

	e.GET("/news/search", func(c echo.Context) error {
		const limit = 150
		query := c.QueryParam("q")
		if len(query) == 0 {
			return c.JSON(http.StatusNotAcceptable, nil)
		}

		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			page = 1
		}

		searcher, err := sonicPool.Get()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return err
		}

		bucket := strings.ToLower(c.QueryParam("newsgroup"))
		if bucket == "" {
			bucket = "all"
		}
		resp, err := searcher.(sonic.Searchable).Query("news", bucket, query, limit, (page-1)*limit, "")
		if err != nil {
			c.NoContent(http.StatusInternalServerError)
			return err
		}

		sonicPool.Put(searcher)

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

		news, err := models.GetNewsList(ids)
		if err != nil {
			c.NoContent(http.StatusInternalServerError)
			return err
		}

		return c.JSON(200, news)
	})

	groups, err := models.GetUniqueNewsgroups()
	if err != nil {
		zap.S().Fatal("Failed to get newsgroups", zap.Error(err))
	}
	e.GET("/newsgroups", func(c echo.Context) error {
		return c.JSON(http.StatusOK, groups)
	})

	// e.GET("/file/:id", func(c echo.Context) error {
	// 	id := c.Param("id")
	// 	meta, err := models.GetFileMeta(id)
	// 	if err != nil {
	// 		c.NoContent(http.StatusInternalServerError)
	// 		return err
	// 	}

	// 	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", meta.Name))
	// 	return c.File(path.Join(config.Cfg.Index.Store, "files", id))
	// })

	e.GET("/modules", func(c echo.Context) error {
		var modules []config.ModuleMeta

		for _, m := range config.Cfg.Modules {
			if m.Enable && !m.Hide {
				modules = append(modules, m.Meta)
			}
		}

		return c.JSON(http.StatusOK, modules)
	})

	var stats models.Stats
	go func() {
		for {
			stats, err = models.GetStats()
			if err != nil {
				zap.S().Error("Failed to get stats", zap.Error(err))
			}
			time.Sleep(time.Minute)
		}
	}()

	e.GET("/stats", func(c echo.Context) error {
		return c.JSON(200, stats)
	})

	e.Start(":1323")
}
