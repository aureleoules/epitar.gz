package archive

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/aureleoules/epitar/config"
	"github.com/briandowns/spinner"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/fatih/color"
)

func (module *Module) Run() error {
	str := color.GreenString(" Running '%s' archive module...", module.Meta.Slug)

	s := spinner.New(spinner.CharSets[4], 100*time.Millisecond, spinner.WithSuffix(str), spinner.WithFinalMSG(color.GreenString("Done.\n")))
	s.Start()

	go func(c *spinner.Spinner) {
		for {
			select {
			case <-stopCh:
				s.Stop()
				return
			}
		}
	}(s)

	// Resolve absolute path of module
	output, err := filepath.Abs(config.Cfg.Output)
	if err != nil {
		return err
	}

	source := path.Join(output, module.Meta.Slug)
	err = os.MkdirAll(source, 0755)
	if err != nil {
		return err
	}

	cont, err := dockerClient.ContainerCreate(context.Background(), &container.Config{
		Image: "epitar-module-" + module.Meta.Slug,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: source,
				Target: "/output",
			},
		},
	}, nil, nil, "")
	if err != nil {
		color.Red("Error creating container: %s", err)
		return err
	}

	err = dockerClient.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	if err != nil {
		color.Red("Error starting container: %s", err)
		return err
	}

	// Set container ID
	module.ContainerID = cont.ID

	// Wait for container to finish
	statusCh, errCh := dockerClient.ContainerWait(context.Background(), cont.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			color.Red("Error waiting for container: %s", err)
			return err
		}
	case <-statusCh:
	}

	s.Stop()

	time.Sleep(time.Second * 1)

	err = dockerClient.ContainerRemove(context.Background(), cont.ID, types.ContainerRemoveOptions{
		Force: true,
	})

	if err != nil {
		color.Red("Error removing container: %s", err)
		return err
	}

	module.ContainerID = ""

	return nil
}

func RunModules() error {
	for _, module := range modules {
		if stop {
			break
		}

		if err := module.Run(); err != nil {
			return err
		}
	}

	return nil
}

func (module *Module) Stop() error {
	if module.ContainerID == "" {
		return nil
	}

	str := color.YellowString(" Stopping '%s' archive module...", module.Meta.Name)

	s := spinner.New(spinner.CharSets[4], 100*time.Millisecond, spinner.WithSuffix(str), spinner.WithFinalMSG(color.YellowString("Stopped.\n")))
	s.Start()

	err := dockerClient.ContainerRemove(context.Background(), module.ContainerID, types.ContainerRemoveOptions{
		Force: true,
	})

	if err != nil {
		color.Red("Error stopping container: %s", err)
		return err
	}

	module.ContainerID = ""

	s.Stop()
	return nil
}

func stopModules() error {
	stopCh <- true
	stop = true

	for _, module := range modules {
		if err := module.Stop(); err != nil {
			return err
		}
	}

	return nil
}
