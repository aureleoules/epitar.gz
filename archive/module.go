package archive

import "github.com/aureleoules/epitar/config"

type Module struct {
	config.Module

	ContainerID string
}
