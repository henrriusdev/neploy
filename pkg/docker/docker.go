package docker

import (
	"context"
	"io"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Docker struct {
	cli *client.Client
}

func NewDocker() *Docker {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	return &Docker{cli: cli}
}

func (d *Docker) ListContainers(ctx context.Context) ([]types.Container, error) {
	return d.cli.ContainerList(ctx, container.ListOptions{})
}

func (d *Docker) CreateContainer(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, name string) (container.CreateResponse, error) {
	return d.cli.ContainerCreate(ctx, config, hostConfig, nil, nil, name)
}

func (d *Docker) StartContainer(ctx context.Context, containerName string) error {
	containers, err := d.ListContainers(ctx)
	if err != nil {
		return err
	}

	var containerID string

	for _, container := range containers {
		if container.Names[0] == "/"+containerName {
			containerID = container.ID
			break
		}
	}

	return d.cli.ContainerStart(ctx, containerID, container.StartOptions{})
}

func (d *Docker) StopContainer(ctx context.Context, containerName string) error {
	containers, err := d.ListContainers(ctx)
	if err != nil {
		return err
	}

	var containerID string

	for _, container := range containers {
		if container.Names[0] == "/"+containerName {
			containerID = container.ID
			break
		}
	}

	return d.cli.ContainerStop(ctx, containerID, container.StopOptions{})
}

func (d *Docker) PauseContainer(ctx context.Context, containerName string) error {
	containers, err := d.ListContainers(ctx)
	if err != nil {
		return err
	}

	var containerID string

	for _, container := range containers {
		if container.Names[0] == "/"+containerName {
			containerID = container.ID
			break
		}
	}

	return d.cli.ContainerPause(ctx, containerID)
}

func (d *Docker) RemoveContainer(ctx context.Context, containerID string) error {
	return d.cli.ContainerRemove(ctx, containerID, container.RemoveOptions{})
}

func (d *Docker) ContainerLogs(ctx context.Context, containerID string) (io.ReadCloser, error) {
	return d.cli.ContainerLogs(ctx, containerID, container.LogsOptions{})
}

func (d *Docker) GetContainerID(ctx context.Context, containerName string) (string, error) {
	containers, err := d.ListContainers(ctx)
	if err != nil {
		return "", err
	}

	for _, container := range containers {
		if container.Names[0] == "/"+containerName {
			return container.ID, nil
		}
	}

	return "", nil
}

func (d *Docker) BuildImage(ctx context.Context, dockerfilePath string, tag string) error {
	// Create build context from the directory containing the Dockerfile
	path := filepath.Dir(dockerfilePath)

	options := types.ImageBuildOptions{
		Dockerfile: path,
		Tags:       []string{tag},
		Remove:     true,
	}

	// BuildContext is the directory where the Dockerfile is located
	response, err := d.cli.ImageBuild(ctx, nil, options)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Read the response to ensure the build completes
	_, err = io.Copy(io.Discard, response.Body)
	return err
}
