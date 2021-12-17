package cmd

import (
	"io/ioutil"
	"os"

	"github.com/aureleoules/epitar/api"
	"github.com/aureleoules/epitar/archive"
	"github.com/aureleoules/epitar/config"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
)

var logger *zap.Logger

func init() {
	c := zap.NewDevelopmentConfig()
	c.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ = c.Build()
	zap.ReplaceGlobals(logger)
}

func loadConfig(path string) error {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		color.Red("Error reading config file")
		return err
	}

	err = yaml.Unmarshal(yamlFile, &config.Cfg)
	if err != nil {
		color.Red("Could not decode config file: %s\n%s", path, err)
		return err
	}

	return nil
}

func Run() {
	app := &cli.App{
		Name:  "epitar",
		Usage: "archive epita services",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      "config",
				Value:     "config.yml",
				TakesFile: true,
			},
		},
		Commands: []*cli.Command{
			{
				Name: "start",
				Action: func(c *cli.Context) error {
					err := loadConfig(c.String("config"))
					if err != nil {
						return err
					}

					for _, module := range config.Cfg.Modules {
						if module.Enable {
							archive.Register(module)
						}
					}

					err = archive.BuildModules()
					if err != nil {
						color.Red("Error building modules: %s", err)
						return err
					}

					err = archive.RunModules()
					if err != nil {
						color.Red("Error running modules: %s", err)
						return err
					}

					color.Green("Done archiving.")

					err = archive.IndexModules()
					if err != nil {
						color.Red("Error indexing modules: %s", err)
						return err
					}

					return nil
				},
			},
			{
				Name: "serve",
				Action: func(c *cli.Context) error {
					err := loadConfig(c.String("config"))
					if err != nil {
						return err
					}

					api.Serve()
					return nil
				},
			},
		},
	}

	app.Run(os.Args)
}
