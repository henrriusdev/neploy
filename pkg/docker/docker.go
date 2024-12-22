package docker

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"neploy.dev/pkg/logger"
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
	// Use the directory containing the Dockerfile as build context
	contextDir := filepath.Dir(dockerfilePath)
	logger.Info("Building image from context: %s", contextDir)

	// Create tar archive of the build context
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	defer tw.Close()

	// Walk through the directory and add files to tar
	err := filepath.Walk(contextDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path from context directory
		relPath, err := filepath.Rel(contextDir, path)
		if err != nil {
			return err
		}

		// Skip if path is outside context
		if strings.HasPrefix(relPath, "..") {
			return nil
		}

		logger.Info("Adding file to tar: %s", path)

		// Create tar header
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}
		header.Name = relPath

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if _, err := tw.Write(data); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	options := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{tag},
		Remove:     true,
	}

	// Build the image using the tar context
	response, err := d.cli.ImageBuild(ctx, bytes.NewReader(buf.Bytes()), options)
	if err != nil {
		return errors.New("error building image: " + err.Error())
	}
	defer response.Body.Close()

	// Read the build output to check for errors
	var lastError string
	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		var output map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &output); err != nil {
			continue
		}

		if errMsg, ok := output["error"]; ok {
			lastError = errMsg.(string)
		}

		if stream, ok := output["stream"]; ok {
			logger.Info("Build: %v", stream)
		}
	}

	if lastError != "" {
		return errors.New("build failed: " + lastError)
	}

	if err := scanner.Err(); err != nil {
		return errors.New("error reading build output: " + err.Error())
	}

	return nil
}
