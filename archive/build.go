package archive

import (
	"bufio"
	"context"
	"time"

	"github.com/briandowns/spinner"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"github.com/fatih/color"
)

func (module *Module) Build() error {
	str := color.CyanString(" Building '%s' archive module...", module.Meta.Slug)

	s := spinner.New(spinner.CharSets[4], 100*time.Millisecond, spinner.WithSuffix(str))
	s.Start()

	tar, err := archive.TarWithOptions(module.Path, &archive.TarOptions{})
	if err != nil {
		return err
	}

	resp, err := dockerClient.ImageBuild(context.Background(), tar, types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{"epitar-module-" + module.Meta.Slug},
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

	return nil
}

func BuildModules() error {
	for _, module := range modules {
		if err := module.Build(); err != nil {
			return err
		}
	}

	return nil
}
