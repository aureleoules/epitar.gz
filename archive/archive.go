package archive

import (
	"bufio"
	"context"
	"time"

	"github.com/aureleoules/epitar/config"
	"github.com/briandowns/spinner"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
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

func BuildModules() error {
	for _, module := range modules {
		str := color.CyanString(" Building '%s' archive module...", module.Name)

		s := spinner.New(spinner.CharSets[4], 100*time.Millisecond, spinner.WithSuffix(str)) // Build our new spinner
		s.Start()

		tar, err := archive.TarWithOptions(module.Path, &archive.TarOptions{})
		if err != nil {
			return err
		}

		resp, err := dockerClient.ImageBuild(context.Background(), tar, types.ImageBuildOptions{
			Dockerfile: "Dockerfile",
			Tags:       []string{"epitar-module-" + module.Name},
		})

		if err != nil {
			return err
		}

		// Wait for build to finish
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
		}

		s.Stop()
		color.Cyan("Built.")
	}

	return nil
}

func (module *Module) Run() error {
	str := color.GreenString(" Running '%s' archive module...", module.Name)

	s := spinner.New(spinner.CharSets[4], 100*time.Millisecond, spinner.WithSuffix(str), spinner.WithFinalMSG(color.GreenString("Done.\n"))) // Build our new spinner
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

	cont, err := dockerClient.ContainerCreate(context.Background(), &container.Config{
		Image: "epitar-module-" + module.Name,
	}, nil, nil, nil, "")
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

	return nil
}

func RunModules() error {
	for _, module := range modules {
		if stop {
			break
		}
		module.Run()
	}

	return nil
}

func StopModules() error {
	stopCh <- true
	stop = true

	for _, module := range modules {
		if module.ContainerID == "" {
			continue
		}

		str := color.YellowString(" Stopping '%s' archive module...", module.Name)

		s := spinner.New(spinner.CharSets[4], 100*time.Millisecond, spinner.WithSuffix(str), spinner.WithFinalMSG(color.YellowString("Stopped.\n")))
		s.Start()

		err := dockerClient.ContainerRemove(context.Background(), module.ContainerID, types.ContainerRemoveOptions{
			Force: true,
		})

		if err != nil {
			color.Red("Error stopping container: %s", err)
			return err
		}

		s.Stop()
	}

	return nil
}
