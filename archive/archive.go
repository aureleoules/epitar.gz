package archive

import (
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
}

func Register(config config.Module) {
	modules = append(modules, &Module{
		Module: config,
	})
}
