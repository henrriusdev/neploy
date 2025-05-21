package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"neploy.dev/config"
	neploker "neploy.dev/pkg/docker"
	neployway "neploy.dev/pkg/gateway"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
	"neploy.dev/pkg/websocket"
)

type Docker interface {
	CreateAndStartContainer(ctx context.Context, app model.Application, version model.ApplicationVersion, port string) error
	StartContainer(ctx context.Context, id, versionId string) error
	StopContainer(ctx context.Context, id, versionId string) error
	ConfigurePort(dockerfilePath string, interactive bool) (string, error)
}

type docker struct {
	repos  repository.Repositories
	hub    *websocket.Hub
	docker *neploker.Docker
	router *neployway.Router
}

func NewDocker(repos repository.Repositories, hub *websocket.Hub, dckr *neploker.Docker, router *neployway.Router) Docker {
	return &docker{repos, hub, dckr, router}
}

func (d *docker) CreateAndStartContainer(ctx context.Context, app model.Application, version model.ApplicationVersion, port string) error {
	appName := sanitizeAppName(app.AppName)
	imageName := fmt.Sprintf("neploy/%s:%s", appName, version.VersionTag)
	containerName := getContainerName(appName, version.VersionTag)

	if d.hub != nil {
		d.hub.BroadcastProgress(0, "Building Docker image...")
	}

	dockerfile := filepath.Join(version.StorageLocation, "Dockerfile")
	if err := d.docker.BuildImage(ctx, dockerfile, imageName); err != nil {
		logger.Error("error building image: %v", err)
		return err
	}

	hostConfig := &container.HostConfig{
		AutoRemove: true,
		PortBindings: nat.PortMap{
			nat.Port(port + "/tcp"): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: port}},
		},
	}

	cfg := &container.Config{
		Image: imageName,
		Tty:   true,
		ExposedPorts: nat.PortSet{
			nat.Port(port + "/tcp"): struct{}{},
		},
	}

	resp, err := d.docker.CreateContainer(context.Background(), cfg, hostConfig, containerName)
	if err != nil && !strings.Contains(err.Error(), "already in use") {
		logger.Error("error creating container: %v", err)
		return err
	}

	if err := d.docker.StartContainer(ctx, resp.ID); err != nil {
		logger.Error("error starting container: %v", err)
		return err
	}

	gateway := model.Gateway{
		Name:          appName + "-gateway",
		EndpointType:  "path",
		Domain:        config.Env.DefaultDomain,
		EndpointURL:   "/" + appName,
		Port:          port,
		Path:          "/" + appName,
		Status:        "active",
		ApplicationID: app.ID,
	}
	if _, err := d.repos.Gateway.UpsertOneDoUpdate(ctx, gateway, "name"); err != nil {
		logger.Error("error creating gateway: %v", err)
	}

	route := neployway.Route{
		AppID:  app.ID,
		Port:   port,
		Domain: config.Env.DefaultDomain,
		Path:   "/" + appName,
	}
	if err := d.router.AddRoute(route); err != nil {
		logger.Error("Failed to add route: %v", err)
		return err
	}

	if d.hub != nil {
		d.hub.BroadcastProgress(100, fmt.Sprintf("Container %s started successfully!", resp.ID[:12]))
	}

	return nil
}

func (d *docker) StartContainer(ctx context.Context, id, versionId string) error {
	app, err := d.repos.Application.GetByID(ctx, id)
	if err != nil {
		return err
	}

	version, err := d.repos.ApplicationVersion.GetOneById(ctx, versionId)
	if err != nil {
		return err
	}

	appName := getContainerName(app.AppName, version.VersionTag)
	containerId, err := d.docker.GetContainerID(ctx, appName)
	if err != nil {
		return err
	}
	if err := d.docker.StartContainer(ctx, containerId); err != nil {
		return err
	}

	version.Status = "active"
	if _, err := d.repos.ApplicationVersion.UpdateOneById(ctx, version.ID, version); err != nil {
		logger.Error("error updating application version: %v", err)
	}

	return nil
}

func (d *docker) StopContainer(ctx context.Context, id, versionId string) error {
	app, err := d.repos.Application.GetByID(ctx, id)
	if err != nil {
		return err
	}

	version, err := d.repos.ApplicationVersion.GetOneById(ctx, versionId)
	if err != nil {
		return err
	}

	appName := getContainerName(app.AppName, version.VersionTag)
	containerId, err := d.docker.GetContainerID(ctx, appName)
	if err != nil {
		return err
	}
	err = d.docker.StopContainer(ctx, containerId)
	if err != nil {
		return err
	}

	version.Status = "inactive"
	if _, err := d.repos.ApplicationVersion.UpdateOneById(ctx, version.ID, version); err != nil {
		logger.Error("error updating application version: %v", err)
	}

	return nil
}

func (d *docker) ConfigurePort(dockerfilePath string, interactive bool) (string, error) {
	logger.Info("configuring port for Dockerfile: %s", dockerfilePath)
	content, err := os.ReadFile(dockerfilePath)
	if err != nil {
		return "", err
	}
	port := "3000"
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "EXPOSE") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				port = parts[1]
				break
			}
		}
	}

	if d.hub != nil && interactive {
		for i := 0; i < 5; i++ {
			if d.hub.GetInteractiveClient() != nil {
				break
			}
			time.Sleep(2 * time.Second)
		}
		client := d.hub.GetInteractiveClient()
		if client == nil {
			return "", fmt.Errorf("no interactive client available")
		}

		response := d.hub.BroadcastInteractive(websocket.ActionMessage{
			Type:    "critical",
			Action:  "expose",
			Title:   "Port Configuration",
			Message: fmt.Sprintf("The application wants to expose port %s. You can change this if needed.", port),
			Inputs: []websocket.Input{{
				Name:        "port",
				Type:        "text",
				Placeholder: "Enter port number",
				Value:       port,
				Required:    true,
				Order:       1,
			}},
		})
		if response != nil && response.Data["port"] != "" {
			port = response.Data["port"].(string)

			newContent := ""
			foundExpose := false
			for _, line := range lines {
				if strings.HasPrefix(strings.TrimSpace(line), "EXPOSE") {
					newContent += fmt.Sprintf("EXPOSE %s\n", port)
					foundExpose = true
				} else {
					newContent += line + "\n"
				}
			}
			if !foundExpose {
				newContent += fmt.Sprintf("\nEXPOSE %s\n", port)
			}
			if err := os.WriteFile(dockerfilePath, []byte(newContent), 0o644); err != nil {
				return "", err
			}
		} else {
			return "", fmt.Errorf("no response received from interactive client")
		}
	}
	return port, nil
}

func getContainerName(appName, versionTag string) string {
	safeApp := sanitizeAppName(appName)
	safeTag := strings.ReplaceAll(versionTag, ".", "-") // Opcional: evita puntos
	return fmt.Sprintf("neploy-%s_v%s", safeApp, strings.TrimPrefix(safeTag, "v"))
}
