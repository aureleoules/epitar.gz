package cmd

import (
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/aureleoules/epitar/archive"
	"github.com/aureleoules/epitar/config"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// Run Cleanup
		archive.StopModules()

		os.Exit(1)
	}()
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
					yamlFile, err := ioutil.ReadFile(c.String("config"))
					if err != nil {
						color.Red("Error reading config file")
						return err
					}

					err = yaml.Unmarshal(yamlFile, &config.Cfg)
					if err != nil {
						color.Red("Could not decode config file: %s\n%s", c.String("config"), err)
						return err
					}

					for _, module := range config.Cfg.Modules {
						if module.Enabled {
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

					return nil
				},
			},
		},
	}

	app.Run(os.Args)
}
