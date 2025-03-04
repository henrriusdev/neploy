package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"neploy.dev/config"
	"neploy.dev/pkg/docker"
	"neploy.dev/pkg/filesystem"
	neployway "neploy.dev/pkg/gateway"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
	"neploy.dev/pkg/websocket"
)

type Application interface {
	Create(ctx context.Context, app model.Application) (string, error)
	Get(ctx context.Context, id string) (model.Application, error)
	GetAll(ctx context.Context) ([]model.FullApplication, error)
	Update(ctx context.Context, app model.Application) error
	GetStat(ctx context.Context, id string) (model.ApplicationStat, error)
	CreateStat(ctx context.Context, stat model.ApplicationStat) error
	UpdateStat(ctx context.Context, stat model.ApplicationStat) error
	GetHealthy(ctx context.Context) (uint, uint, error)
	Deploy(ctx context.Context, id string, repoURL string, branch string) error
	Upload(ctx context.Context, id string, file *multipart.FileHeader) (string, error)
	Delete(ctx context.Context, id string) error
	StartContainer(ctx context.Context, id string) error
	StopContainer(ctx context.Context, id string) error
	GetRepoBranches(ctx context.Context, repoURL string) ([]string, error)
}

type application struct {
	repos  repository.Repositories
	router *neployway.Router
	hub    *websocket.Hub
	docker *docker.Docker
}

func NewApplication(repos repository.Repositories, router *neployway.Router) Application {
	return &application{repos, router, websocket.GetHub(), docker.NewDocker()}
}

func (a *application) Create(ctx context.Context, app model.Application) (string, error) {
	return a.repos.Application.Insert(ctx, app)
}

func (a *application) Get(ctx context.Context, id string) (model.Application, error) {
	return a.repos.Application.GetByID(ctx, id)
}

func (a *application) GetAll(ctx context.Context) ([]model.FullApplication, error) {
	apps, err := a.repos.Application.GetAll(ctx)
	if err != nil {
		logger.Error("error getting applications: %v", err)
		return nil, err
	}

	var fullApps []model.FullApplication
	for _, app := range apps {
		stats, err := a.repos.ApplicationStat.GetByApplicationID(ctx, app.ID)
		if err != nil {
			logger.Error("error getting application stat: %v", err)
			return nil, err
		}

		var tech model.TechStack
		if app.TechStackID != nil {
			tech, err = a.repos.TechStack.GetByID(ctx, *app.TechStackID)
			if err != nil {
				logger.Error("error getting tech stack: %v", err)
				return nil, err
			}
		}

		appNameWithoutSpace := strings.ReplaceAll(app.AppName, " ", "-")
		appNameWithoutSpecialChars := regexp.MustCompile(`[^a-zA-Z0-9-]`).ReplaceAllString(appNameWithoutSpace, "")
		appName := strings.ToLower(appNameWithoutSpecialChars)
		containerName := "neploy-" + strings.ToLower(appName)
		imageName := "neploy/" + strings.ToLower(appName)

		status, err := a.docker.GetContainerStatus(ctx, containerName)
		if err != nil {
			logger.Error("error getting container status: %v", err)
			return nil, err
		}
		logger.Info("Container status: %s", status)

		if status == "Not created" {
			// get port from dockerfile
			go func() error {
				dockerfilePath := filepath.Join(app.StorageLocation, "Dockerfile")
				port, err := a.configurePort(dockerfilePath, false)
				if err != nil {
					logger.Error("error configuring port: %v", err)
					return err
				}

				// create container
				a.createAndStartContainer(ctx, imageName, containerName, app.StorageLocation, app.ID, port)

				status, err = a.docker.GetContainerStatus(ctx, containerName)
				if err != nil {
					logger.Error("error getting container status: %v", err)
					return err
				}

				return nil
			}()
		}

		fullApps = append(fullApps, model.FullApplication{
			Application: app,
			TechStack:   tech,
			Stats:       stats,
			Status:      status,
		})
	}

	return fullApps, nil
}

func (a *application) Update(ctx context.Context, app model.Application) error {
	return a.repos.Application.Update(ctx, app)
}

func (a *application) GetStat(ctx context.Context, id string) (model.ApplicationStat, error) {
	return a.repos.ApplicationStat.GetByID(ctx, id)
}

func (a *application) CreateStat(ctx context.Context, stat model.ApplicationStat) error {
	return a.repos.ApplicationStat.Insert(ctx, stat)
}

func (a *application) UpdateStat(ctx context.Context, stat model.ApplicationStat) error {
	return a.repos.ApplicationStat.Update(ctx, stat)
}

func (a *application) GetHealthy(ctx context.Context) (uint, uint, error) {
	apps, err := a.repos.ApplicationStat.GetAll(ctx)
	if err != nil {
		logger.Error("error getting all application stats: %v", err)
		return 0, 0, err
	}

	var healthy uint = 3
	for _, app := range apps {
		if app.Healthy {
			healthy++
		}
	}

	totalApps := uint(len(apps))
	return healthy, totalApps, nil
}

func (a *application) Deploy(ctx context.Context, id string, repoURL string, branch string) error {
	app, err := a.repos.Application.GetByID(ctx, id)
	if err != nil {
		logger.Error("error getting application: %v", err)
		return err
	}

	appNameWithoutSpace := strings.ReplaceAll(app.AppName, " ", "-")
	appNameWithoutSpecialChars := regexp.MustCompile(`[^a-zA-Z0-9-]`).ReplaceAllString(appNameWithoutSpace, "")
	appName := strings.ToLower(appNameWithoutSpecialChars)
	imageName := fmt.Sprintf("neploy/%s", appName)
	containerName := fmt.Sprintf("neploy-%s", appName)

	path := filepath.Join(config.Env.UploadPath, appName)
	repo := filesystem.NewGitRepo(repoURL)
	if err := repo.Clone(path, branch); err != nil {
		logger.Error("error cloning repository: %v", err)
		return err
	}

	techStack, err := filesystem.DetectStack(path)
	if err != nil {
		logger.Error("error detecting tech stack: %v", err)
		return err
	}

	tech, err := a.repos.TechStack.FindOrCreate(ctx, techStack)
	if err != nil {
		logger.Error("error finding or creating tech stack: %v", err)
		return err
	}

	// Broadcast progress if hub exists
	if a.hub != nil {
		a.hub.BroadcastProgress(0, "Checking for Dockerfile...")
	}

	// Get notification client (may be nil)
	var wsClient *websocket.Client
	if a.hub != nil {
		wsClient = a.hub.GetNotificationClient()
	}

	dockerStatus := filesystem.HasDockerfile(path, wsClient)
	if !dockerStatus.Exists {
		// Only send interactive message if hub exists
		if a.hub != nil {
			actionInput := websocket.NewSelectInput("action", []string{
				"create",
				"skip",
			})

			actionMsg := websocket.NewActionMessage(
				websocket.ActionTypeCritical,
				"Dockerfile Required",
				"No Dockerfile found for application. Would you like to create one?",
				[]websocket.Input{actionInput},
			)

			a.hub.BroadcastInteractive(actionMsg)
		}

		// Only broadcast progress if hub exists
		if a.hub != nil {
			a.hub.BroadcastProgress(50, "Creating Dockerfile...")
		}

		tmpl, ok := docker.GetDefaultTemplate(techStack)
		if !ok {
			logger.Error("no default template for tech stack: %s", techStack)
			if a.hub != nil {
				a.hub.BroadcastProgress(100, "Error: No default template available for "+techStack)
			}
			return err
		}

		dockerfilePath := filepath.Join(path, "Dockerfile")
		if err := docker.WriteDockerfile(dockerfilePath, tmpl); err != nil {
			logger.Error("error writing dockerfile: %v", err)
			if a.hub != nil {
				a.hub.BroadcastProgress(100, "Error creating Dockerfile")
			}
			return err
		}

		if a.hub != nil {
			a.hub.BroadcastProgress(100, "Created default Dockerfile")
		}
	}

	// Configure port
	dockerfilePath := filepath.Join(path, "Dockerfile")
	port, err := a.configurePort(dockerfilePath, true)
	if err != nil {
		logger.Error("error configuring port: %v", err)
		if a.hub != nil {
			a.hub.BroadcastProgress(100, fmt.Sprintf("Error configuring port: %v", err))
		}
		return err
	}

	// Update application
	app.TechStackID = &tech.ID
	if err := a.repos.Application.Update(ctx, app); err != nil {
		logger.Error("error updating application: %v", err)
		return err
	}

	a.createAndStartContainer(ctx, imageName, containerName, path, app.ID, port)

	logger.Info("application updated: %s", app.AppName)
	if a.hub != nil {
		a.hub.BroadcastProgress(100, "Deployment complete!")
	}
	return nil
}

func (a *application) createAndStartContainer(ctx context.Context, imageName, containerName, projectPath string, appID string, port string) {
	// First, build the Docker image
	if a.hub != nil {
		a.hub.BroadcastProgress(0, "Building Docker image...")
	}

	if err := a.docker.BuildImage(context.Background(), filepath.Join(projectPath, "Dockerfile"), imageName); err != nil {
		logger.Error("error building image: %v", err)
		if a.hub != nil {
			a.hub.BroadcastProgress(100, "Error building Docker image")
		}
		return
	}

	if a.hub != nil {
		a.hub.BroadcastProgress(50, "Docker image built successfully")
	}

	// Get a free port for the container
	hostConfig := &container.HostConfig{
		AutoRemove: true,
		PortBindings: nat.PortMap{
			nat.Port(port + "/tcp"): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: port,
				},
			},
		},
	}

	// Create container config with exposed port
	cfg := &container.Config{
		Image: imageName,
		Tty:   true,
		ExposedPorts: nat.PortSet{
			nat.Port(port + "/tcp"): struct{}{},
		},
	}

	if a.hub != nil {
		a.hub.BroadcastProgress(70, "Creating container...")
	}

	resp, err := a.docker.CreateContainer(context.Background(), cfg, hostConfig, containerName)
	if err != nil {
		logger.Error("error creating container: %v", err)
		if a.hub != nil {
			a.hub.BroadcastProgress(100, "Error creating container")
		}
		return
	}

	// Create default gateway for the application
	gateway := model.Gateway{
		Name:          containerName + "-gateway",
		EndpointType:  "subdomain",
		Domain:        config.Env.DefaultDomain,
		EndpointURL:   "/" + containerName,
		Subdomain:     strings.Replace(containerName, "neploy", "", -1),
		Port:          port,
		Path:          "/" + containerName,
		Status:        "inactive",
		ApplicationID: appID,
	}

	if err := a.repos.Gateway.Insert(ctx, gateway); err != nil {
		logger.Error("error creating gateway: %v", err)
		// Continue even if gateway creation fails
	}

	if a.hub != nil {
		a.hub.BroadcastProgress(90, "Starting container...")
	}

	if err := a.docker.StartContainer(ctx, resp.ID); err != nil {
		logger.Error("error starting container: %v", err)
		if a.hub != nil {
			a.hub.BroadcastProgress(100, "Error starting container")
		}
		return
	}

	// Update gateway status to active
	gateway.Status = "active"
	if err := a.repos.Gateway.Update(ctx, gateway); err != nil {
		logger.Error("error updating gateway status: %v", err)
	}

	route := neployway.Route{
		AppID:     appID,
		Port:      port,
		Domain:    config.Env.DefaultDomain,
		Path:      "/" + containerName,
		Subdomain: "",
	}
	if err := a.router.AddRoute(route); err != nil {
		logger.Error("Failed to add route: %v", err)
		return
	}

	if a.hub != nil {
		a.hub.BroadcastProgress(100, fmt.Sprintf("Container %s started successfully!", resp.ID[:12]))
	}
}

func (a *application) configurePort(dockerfilePath string, interactive bool) (string, error) {
	logger.Info("configuring port for Dockerfile: %s", dockerfilePath)
	content, err := os.ReadFile(dockerfilePath)
	if err != nil {
		logger.Error("error reading dockerfile: %v", err)
		return "", err
	}

	// Find EXPOSE directive and ask user for port
	port := "3000" // Default port if not found
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

	// Ask user to confirm or change port
	if a.hub != nil && interactive {
		// Wait for interactive client to be available
		for i := 0; i < 5; i++ {
			if a.hub.GetInteractiveClient() != nil {
				break
			}
			logger.Info("waiting for interactive client to be available...")
			time.Sleep(2 * time.Second)
		}

		if a.hub.GetInteractiveClient() == nil {
			return "", fmt.Errorf("no interactive client available")
		}

		response := a.hub.BroadcastInteractive(websocket.ActionMessage{
			Type:    "critical",
			Action:  "expose",
			Title:   "Port Configuration",
			Message: fmt.Sprintf("The application wants to expose port %s. You can change this if needed.", port),
			Inputs: []websocket.Input{
				{
					Name:        "port",
					Type:        "text",
					Placeholder: "Enter port number",
					Value:       port,
					Required:    true,
					Order:       1,
				},
			},
		})

		if response != nil && response.Data["port"] != "" {
			port = response.Data["port"].(string)
			logger.Info("using user-specified port: %s", port)

			// Update Dockerfile with new port
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

			// Add EXPOSE if it wasn't found
			if !foundExpose {
				newContent += fmt.Sprintf("\nEXPOSE %s\n", port)
			}

			// Write updated Dockerfile
			if err := os.WriteFile(dockerfilePath, []byte(newContent), 0o644); err != nil {
				logger.Error("error updating dockerfile: %v", err)
				return "", err
			}
		} else {
			return "", fmt.Errorf("no response received from interactive client")
		}
	}

	return port, nil
}

func (a *application) Upload(ctx context.Context, id string, file *multipart.FileHeader) (string, error) {
	app, err := a.repos.Application.GetByID(ctx, id)
	if err != nil {
		logger.Error("error getting application: %v", err)
		return "", err
	}

	zipPath, err := filesystem.UploadFile(file, app.AppName)
	if err != nil {
		logger.Error("error uploading file: %v", err)
		return "", err
	}

	path, err := filesystem.UnzipFile(zipPath, app.AppName)
	if err != nil {
		logger.Error("error unzipping file: %v", err)
		return "", err
	}

	techStack, err := filesystem.DetectStack(path)
	if err != nil {
		logger.Error("error detecting tech stack: %v", err)
		return "", err
	}

	tech, err := a.repos.TechStack.FindOrCreate(ctx, techStack)
	if err != nil {
		logger.Error("error finding or creating tech stack: %v", err)
		return "", err
	}

	// Delete zip file
	if err := os.Remove(zipPath); err != nil {
		logger.Error("error deleting zip file: %v", err)
		return "", err
	}

	if a.hub != nil {
		a.hub.BroadcastProgress(0, "Checking for Docker Compose...")
	}

	if filesystem.HasDockerCompose(path) {
		logger.Error("docker-compose file found, not supported")
		if a.hub != nil {
			a.hub.BroadcastProgress(100, "Error: Docker Compose files are not supported")
		}

		// Delete the application and notify
		if err := a.Delete(ctx, id); err != nil {
			logger.Error("error deleting application: %v", err)
		}

		if a.hub != nil {
			actionMsg := websocket.NewActionMessage(
				websocket.ActionTypeError,
				"Docker Compose Not Supported",
				"Docker Compose files are not supported. The application has been deleted.",
				nil,
			)
			a.hub.BroadcastInteractive(actionMsg)
		}
		return "", err
	}

	if a.hub != nil {
		a.hub.BroadcastProgress(0, "Checking for Dockerfile...")
	}

	var wsClient *websocket.Client
	if a.hub != nil {
		wsClient = a.hub.GetNotificationClient()
	}

	dockerStatus := filesystem.HasDockerfile(path, wsClient)
	if !dockerStatus.Exists {
		if a.hub != nil {
			actionInput := websocket.NewSelectInput("action", []string{
				"create",
				"skip",
			})

			actionMsg := websocket.NewActionMessage(
				websocket.ActionTypeCritical,
				"Dockerfile Required",
				"No Dockerfile found for application. Would you like to create one?",
				[]websocket.Input{actionInput},
			)

			a.hub.BroadcastInteractive(actionMsg)
		}

		if a.hub != nil {
			a.hub.BroadcastProgress(50, "Creating Dockerfile...")
		}

		tmpl, ok := docker.GetDefaultTemplate(techStack)
		if !ok {
			logger.Error("no default template for tech stack: %s", techStack)
			if a.hub != nil {
				a.hub.BroadcastProgress(100, "Error: No default template available for "+techStack)
			}
			return "", err
		}

		dockerfilePath := filepath.Join(path, "Dockerfile")
		if err := docker.WriteDockerfile(dockerfilePath, tmpl); err != nil {
			logger.Error("error writing dockerfile: %v", err)
			if a.hub != nil {
				a.hub.BroadcastProgress(100, "Error creating Dockerfile")
			}
			return "", err
		}

		if a.hub != nil {
			a.hub.BroadcastProgress(100, "Created default Dockerfile")
		}
	}

	// Configure port
	dockerfilePath := filepath.Join(path, "Dockerfile")
	port, err := a.configurePort(dockerfilePath, true)
	if err != nil {
		logger.Error("error configuring port: %v", err)
		if a.hub != nil {
			a.hub.BroadcastProgress(100, fmt.Sprintf("Error configuring port: %v", err))
		}
		return "", err
	}

	// Update application
	app.TechStackID = &tech.ID
	if err := a.repos.Application.Update(ctx, app); err != nil {
		logger.Error("error updating application: %v", err)
		return "", err
	}

	appNameWithoutSpace := strings.ReplaceAll(app.AppName, " ", "-")
	appNameWithoutSpecialChars := regexp.MustCompile(`[^a-zA-Z0-9-]`).ReplaceAllString(appNameWithoutSpace, "")
	appName := strings.ToLower(appNameWithoutSpecialChars)

	// Create Docker image name with neploy prefix
	imageName := fmt.Sprintf("neploy/%s", appName)
	containerName := fmt.Sprintf("neploy-%s", appName)

	// Start container creation in a separate goroutine
	go a.createAndStartContainer(ctx, imageName, containerName, path, app.ID, port)

	logger.Info("application updated: %s", app.AppName)
	if a.hub != nil {
		a.hub.BroadcastProgress(100, "Deployment complete!")
	}

	app.StorageLocation = path
	app.TechStackID = &tech.ID
	if err := a.repos.Application.Update(ctx, app); err != nil {
		logger.Error("error updating application: %v", err)
		return "", err
	}

	route := neployway.Route{
		AppID:     app.ID,
		Port:      port,
		Domain:    config.Env.DefaultDomain,
		Path:      "/" + containerName,
		Subdomain: "",
	}
	if err := a.router.AddRoute(route); err != nil {
		logger.Error("Failed to add route: %v", err)
		return "", err
	}

	return path, nil
}

func (a *application) Delete(ctx context.Context, id string) error {
	// Delete associated gateways first
	gateways, err := a.repos.Gateway.GetByApplicationID(ctx, id)
	if err != nil {
		logger.Error("error getting gateways: %v", err)
		return err
	}

	for _, gateway := range gateways {
		if err := a.repos.Gateway.Delete(ctx, gateway.ID); err != nil {
			logger.Error("error deleting gateway: %v", err)
			// Continue with other gateways
		}
	}

	// Get application details
	app, err := a.repos.Application.GetByID(ctx, id)
	if err != nil {
		logger.Error("error getting application: %v", err)
		return err
	}

	// Delete container if it exists
	appNameWithoutSpace := strings.ReplaceAll(app.AppName, " ", "-")
	appNameWithoutSpecialChars := regexp.MustCompile(`[^a-zA-Z0-9-]`).ReplaceAllString(appNameWithoutSpace, "")
	appName := strings.ToLower(appNameWithoutSpecialChars)
	containerName := fmt.Sprintf("neploy-%s", appName)

	if err := a.docker.RemoveContainer(ctx, containerName); err != nil && !strings.Contains(err.Error(), "No such container") {
		logger.Error("error removing container: %v", err)
		// Continue with deletion even if container removal fails
	}

	// Delete application files
	if app.StorageLocation != "" {
		if err := os.RemoveAll(app.StorageLocation); err != nil {
			logger.Error("error removing application files: %v", err)
			// Continue with deletion even if file removal fails
		}
	}

	a.router.RemoveRoute("/" + containerName)

	return a.repos.Application.Delete(ctx, id)
}

func (a *application) StartContainer(ctx context.Context, id string) error {
	app, err := a.repos.Application.GetByID(ctx, id)
	if err != nil {
		logger.Error("error getting application: %v", err)
		return err
	}

	appNameWithoutSpace := strings.ReplaceAll(app.AppName, " ", "-")
	appNameWithoutSpecialChars := regexp.MustCompile(`[^a-zA-Z0-9-]`).ReplaceAllString(appNameWithoutSpace, "")
	appName := strings.ToLower(appNameWithoutSpecialChars)
	containerId, err := a.docker.GetContainerID(ctx, appName)
	if err != nil {
		logger.Error("error getting container ID: %v", err)
		return err
	}

	return a.docker.StartContainer(ctx, containerId)
}

func (a *application) StopContainer(ctx context.Context, id string) error {
	app, err := a.repos.Application.GetByID(ctx, id)
	if err != nil {
		logger.Error("error getting application: %v", err)
		return err
	}

	appNameWithoutSpace := strings.ReplaceAll(app.AppName, " ", "-")
	appNameWithoutSpecialChars := regexp.MustCompile(`[^a-zA-Z0-9-]`).ReplaceAllString(appNameWithoutSpace, "")
	appName := strings.ToLower(appNameWithoutSpecialChars)

	containerId, err := a.docker.GetContainerID(ctx, appName)
	if err != nil {
		logger.Error("error getting container ID: %v", err)
		return err
	}

	return a.docker.StopContainer(ctx, containerId)
}

func (a *application) GetRepoBranches(ctx context.Context, repoURL string) ([]string, error) {
	repo := filesystem.NewGitRepo(repoURL)
	return repo.GetBranches()
}
