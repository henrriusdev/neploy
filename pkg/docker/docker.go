package docker

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
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
	// List all containers (including stopped ones)
	return d.cli.ContainerList(ctx, container.ListOptions{All: true})
}

func (d *Docker) CreateContainer(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, name string) (container.CreateResponse, error) {
	return d.cli.ContainerCreate(ctx, config, hostConfig, nil, nil, name)
}

func (d *Docker) StartContainer(ctx context.Context, containerID string) error {
	return d.cli.ContainerStart(ctx, containerID, container.StartOptions{})
}

func (d *Docker) StopContainer(ctx context.Context, containerID string) error {
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
	return d.cli.ContainerRemove(ctx, containerID, container.RemoveOptions{
		Force: true, // Force remove the container even if it is running
	})
}

func (d *Docker) ContainerLogs(ctx context.Context, containerID string) (io.ReadCloser, error) {
	return d.cli.ContainerLogs(ctx, containerID, container.LogsOptions{})
}

func (d *Docker) GetContainerID(ctx context.Context, containerName string) (string, error) {
	containers, err := d.ListContainers(ctx)
	if err != nil {
		return "", err
	}

	containerName = "neploy-" + containerName

	logger.Info("Container: %s", containerName)
	for _, container := range containers {
		logger.Info("Container: %v", container.Names)
		if container.Names[0] == "/"+containerName {
			return container.ID, nil
		}
	}

	return "", nil
}

func (d *Docker) GetContainerStatus(ctx context.Context, containerName string) (string, error) {
	containers, err := d.ListContainers(ctx)
	if err != nil {
		return "", err
	}

	for _, container := range containers {
		if container.Names[0] == "/"+containerName {
			// Use Status instead of State for more detailed information
			if strings.HasPrefix(container.Status, "Up") {
				return "Running", nil
			} else if strings.HasPrefix(container.Status, "Exited") {
				return "Stopped", nil
			} else if strings.Contains(container.Status, "Created") {
				return "Created", nil
			} else if strings.Contains(container.Status, "Paused") {
				return "Paused", nil
			} else if strings.Contains(container.Status, "Restarting") {
				return "Restarting", nil
			} else if strings.Contains(container.Status, "Removing") {
				return "Removing", nil
			} else if strings.Contains(container.Status, "Dead") {
				return "Error", nil
			}
			return container.Status, nil
		}
	}

	return "Not created", nil
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
			logger.Error("Error walking the path: %v", err)
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

		// Skip .git directory
		if info.IsDir() && (info.Name() == ".git" || info.Name() == "node_modules") {
			return filepath.SkipDir
		}

		// Create tar header
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		// Use relative path for the header name
		header.Name = filepath.ToSlash(relPath)

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
		NoCache:    true,
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

func (d *Docker) RemoveImage(ctx context.Context, imageName string) error {
	_, err := d.cli.ImageRemove(ctx, imageName, image.RemoveOptions{
		Force: true, // Force remove the image even if it is in use
	})
	return err
}

func (d *Docker) GetExposedPorts(ctx context.Context, containerId string) ([]string, error) {
	inspect, err := d.cli.ContainerInspect(ctx, containerId)
	if err != nil {
		return nil, err
	}

	ports := inspect.Config.ExposedPorts
	var exposedPorts []string
	for port := range ports {
		exposedPorts = append(exposedPorts, port.Port())
	}
	return exposedPorts, nil
}

func (d *Docker) GetUsage(ctx context.Context, containerId string) (float64, float64, error) {
	stats, err := d.cli.ContainerStats(ctx, containerId, false)
	if err != nil {
		return 0, 0, err
	}

	defer stats.Body.Close()

	// Leer y decodificar JSON
	var statsData types.StatsJSON
	err = json.NewDecoder(stats.Body).Decode(&statsData)
	if err != nil && err != io.EOF {
		return 0, 0, err
	}

	// Cálculo del uso de CPU
	cpuDelta := float64(statsData.CPUStats.CPUUsage.TotalUsage - statsData.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(statsData.CPUStats.SystemUsage - statsData.PreCPUStats.SystemUsage)
	var cpuPercent float64
	if systemDelta > 0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(statsData.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}

	// Cálculo del uso de RAM
	memUsage := float64(statsData.MemoryStats.Usage)
	memLimit := float64(statsData.MemoryStats.Limit)
	fmt.Println(memUsage, memLimit)
	memPercent := (memUsage / memLimit) * 100.0

	if math.IsNaN(memPercent) {
		memPercent = 0.0
	}

	return cpuPercent, memPercent, nil
}

func (d *Docker) GetUptime(ctx context.Context, containerId string) (time.Duration, error) {
	inspect, err := d.cli.ContainerInspect(ctx, containerId)
	if err != nil {
		return 0, err
	}

	// Parsear el tiempo de inicio del contenedor
	startTime, err := time.Parse(time.RFC3339Nano, inspect.State.StartedAt)
	if err != nil {
		return 0, err
	}

	// Calcular uptime
	uptime := time.Since(startTime)
	return uptime, nil
}

func (d *Docker) GetLogs(ctx context.Context, containerId string, stream bool) ([]string, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     stream,
		Timestamps: false,
	}

	logs, err := d.cli.ContainerLogs(ctx, containerId, options)
	if err != nil {
		return nil, err
	}
	defer logs.Close()

	var logLines []string
	reader := bufio.NewReader(logs)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		logLines = append(logLines, strings.TrimSpace(line)) // Guardar cada línea en el slice
	}

	return logLines, nil
}
