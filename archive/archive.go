package archive

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aureleoules/epitar/config"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
)

var (
	dockerClient *client.Client
	modules      []*Module

	stopCh = make(chan bool)
	stop   = false
)

func init() {
	var err error
	dockerClient, err = client.NewEnvClient()
	if err != nil {
		color.Red("Error initializing docker client: %s", err)
		panic(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if len(modules) > 0 {
			fmt.Println("\nExiting...")
			// Run Cleanup
			stopModules()

			os.Exit(1)
		}
	}()
}

func Register(config config.Module) {
	modules = append(modules, &Module{
		Module: config,
	})
}
