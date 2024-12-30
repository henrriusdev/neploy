package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"neploy.dev/config"
	"neploy.dev/pkg/docker"
	"neploy.dev/pkg/filesystem"
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
	repo    repository.Application
	stat    repository.ApplicationStat
	tech    repository.TechStack
	gateway repository.Gateway
	hub     *websocket.Hub
	docker  *docker.Docker
}

func NewApplication(repo repository.Application, stat repository.ApplicationStat, tech repository.TechStack, gateway repository.Gateway) Application {
	return &application{
		repo:    repo,
		stat:    stat,
		tech:    tech,
		gateway: gateway,
		hub:     websocket.GetHub(),
		docker:  docker.NewDocker(),
	}
}

func (a *application) Create(ctx context.Context, app model.Application) (string, error) {
	return a.repo.Insert(ctx, app)
}

func (a *application) Get(ctx context.Context, id string) (model.Application, error) {
	return a.repo.GetByID(ctx, id)
}

func (a *application) GetAll(ctx context.Context) ([]model.FullApplication, error) {
	apps, err := a.repo.GetAll(ctx)
	if err != nil {
		logger.Error("error getting applications: %v", err)
		return nil, err
	}

	var fullApps []model.FullApplication
	for _, app := range apps {
		stats, err := a.stat.GetByApplicationID(ctx, app.ID)
		if err != nil {
			logger.Error("error getting application stat: %v", err)
			return nil, err
		}

		var tech model.TechStack
		if app.TechStackID != nil {
			tech, err = a.tech.GetByID(ctx, *app.TechStackID)
			if err != nil {
				logger.Error("error getting tech stack: %v", err)
				return nil, err
			}
		}

		fullApps = append(fullApps, model.FullApplication{
			Application: app,
			TechStack:   tech,
			Stats:       stats,
		})
	}

	return fullApps, nil
}

func (a *application) Update(ctx context.Context, app model.Application) error {
	return a.repo.Update(ctx, app)
}

func (a *application) GetStat(ctx context.Context, id string) (model.ApplicationStat, error) {
	return a.stat.GetByID(ctx, id)
}

func (a *application) CreateStat(ctx context.Context, stat model.ApplicationStat) error {
	return a.stat.Insert(ctx, stat)
}

func (a *application) UpdateStat(ctx context.Context, stat model.ApplicationStat) error {
	return a.stat.Update(ctx, stat)
}

func (a *application) GetHealthy(ctx context.Context) (uint, uint, error) {
	apps, err := a.stat.GetAll(ctx)
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
	app, err := a.repo.GetByID(ctx, id)
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

	tech, err := a.tech.FindOrCreate(ctx, techStack)
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

	// Check if Dockerfile has exposed port
	if !filesystem.DockerfileHasExposedPort(path) {
		// Find the Dockerfile path
		dockerStatus := filesystem.HasDockerfile(path, nil)
		if !dockerStatus.Exists {
			logger.Error("no dockerfile found")
			return err
		}

		if a.hub != nil && a.hub.GetInteractiveClient() != nil {
			portInput := websocket.NewTextInput("port", "Enter the port number (e.g. 3000)")
			actionInput := websocket.NewSelectInput("action", []string{
				"expose",
				"skip",
			})

			actionMsg := websocket.NewActionMessage(
				websocket.ActionTypeCritical,
				"Port Required",
				"No exposed port found in Dockerfile. The application needs to expose a port to be accessible. What port would you like to expose?",
				[]websocket.Input{portInput, actionInput},
			)

			response := a.hub.BroadcastInteractive(actionMsg)
			if response != nil && response.Action == "expose" {
				port := response.Data["port"]
				if port == "" {
					port = "3000" // fallback to default if somehow empty
				}

				content, err := os.ReadFile(dockerStatus.Path)
				if err != nil {
					logger.Error("error reading dockerfile: %v", err)
					return err
				}

				// Add EXPOSE with user-specified port before the last line (usually CMD or ENTRYPOINT)
				lines := strings.Split(string(content), "\n")
				if len(lines) > 0 {
					newLines := append(lines[:len(lines)-1], fmt.Sprintf("EXPOSE %s", port), lines[len(lines)-1])
					newContent := strings.Join(newLines, "\n")
					if err := os.WriteFile(dockerStatus.Path, []byte(newContent), 0o644); err != nil {
						logger.Error("error writing dockerfile: %v", err)
						return err
					}
				}
			}
		} else {
			logger.Info("no interactive client connected, using default port 3000")
			// Add default port 3000
			content, err := os.ReadFile(dockerStatus.Path)
			if err != nil {
				logger.Error("error reading dockerfile: %v", err)
				return err
			}

			// Add EXPOSE with default port before the last line
			lines := strings.Split(string(content), "\n")
			if len(lines) > 0 {
				newLines := append(lines[:len(lines)-1], "EXPOSE 3000", lines[len(lines)-1])
				newContent := strings.Join(newLines, "\n")
				if err := os.WriteFile(dockerStatus.Path, []byte(newContent), 0o644); err != nil {
					logger.Error("error writing dockerfile: %v", err)
					return err
				}
			}
		}
	}

	app.TechStackID = &tech.ID
	if err := a.repo.Update(ctx, app); err != nil {
		logger.Error("error updating application: %v", err)
		return err
	}

	a.createAndStartContainer(ctx, imageName, containerName, path, app.ID)

	logger.Info("application updated: %s", app.AppName)
	if a.hub != nil {
		a.hub.BroadcastProgress(100, "Deployment complete!")
	}
	return nil
}

func (a *application) createAndStartContainer(ctx context.Context, imageName, containerName, projectPath string, appID string) {
	// First, build the Docker image
	if a.hub != nil {
		a.hub.BroadcastProgress(0, "Building Docker image...")
	}

	dockerfilePath := filepath.Join(projectPath, "Dockerfile")
	if err := a.docker.BuildImage(context.Background(), dockerfilePath, imageName); err != nil {
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
	port := "3000" // Default port
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
		Subdomain:     containerName,
		Port:          port,
		Status:        "inactive",
		ApplicationID: appID,
	}

	if err := a.gateway.Insert(ctx, gateway); err != nil {
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
	if err := a.gateway.Update(ctx, gateway); err != nil {
		logger.Error("error updating gateway status: %v", err)
	}

	if a.hub != nil {
		a.hub.BroadcastProgress(100, fmt.Sprintf("Container %s started successfully!", resp.ID[:12]))
	}
}

func (a *application) Upload(ctx context.Context, id string, file *multipart.FileHeader) (string, error) {
	app, err := a.repo.GetByID(ctx, id)
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

	tech, err := a.tech.FindOrCreate(ctx, techStack)
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

	// Check if Dockerfile has exposed port
	if !filesystem.DockerfileHasExposedPort(path) {
		// Find the Dockerfile path
		dockerStatus := filesystem.HasDockerfile(path, nil)
		if !dockerStatus.Exists {
			logger.Error("no dockerfile found")
			return "", err
		}

		if a.hub != nil {
			portInput := websocket.NewTextInput("port", "Enter the port number (e.g. 3000)")
			actionInput := websocket.NewSelectInput("action", []string{
				"expose",
				"skip",
			})

			actionMsg := websocket.NewActionMessage(
				websocket.ActionTypeCritical,
				"Port Required",
				"No exposed port found in Dockerfile. The application needs to expose a port to be accessible. What port would you like to expose?",
				[]websocket.Input{portInput, actionInput},
			)

			response := a.hub.BroadcastInteractive(actionMsg)
			if response != nil && response.Action == "expose" {
				port := response.Data["port"]
				if port == "" {
					port = "3000" // fallback to default if somehow empty
				}

				content, err := os.ReadFile(dockerStatus.Path)
				if err != nil {
					logger.Error("error reading dockerfile: %v", err)
					return "", err
				}

				// Add EXPOSE with user-specified port before the last line (usually CMD or ENTRYPOINT)
				lines := strings.Split(string(content), "\n")
				if len(lines) > 0 {
					newLines := append(lines[:len(lines)-1], fmt.Sprintf("EXPOSE %s", port), lines[len(lines)-1])
					newContent := strings.Join(newLines, "\n")
					if err := os.WriteFile(dockerStatus.Path, []byte(newContent), 0o644); err != nil {
						logger.Error("error writing dockerfile: %v", err)
						return "", err
					}
				}
			}
		} else {
			logger.Info("no interactive client connected, using default port 3000")
			// Add default port 3000
			content, err := os.ReadFile(dockerStatus.Path)
			if err != nil {
				logger.Error("error reading dockerfile: %v", err)
				return "", err
			}

			// Add EXPOSE with default port before the last line
			lines := strings.Split(string(content), "\n")
			if len(lines) > 0 {
				newLines := append(lines[:len(lines)-1], "EXPOSE 3000", lines[len(lines)-1])
				newContent := strings.Join(newLines, "\n")
				if err := os.WriteFile(dockerStatus.Path, []byte(newContent), 0o644); err != nil {
					logger.Error("error writing dockerfile: %v", err)
					return "", err
				}
			}
		}
	}

	app.TechStackID = &tech.ID
	if err := a.repo.Update(ctx, app); err != nil {
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
	go a.createAndStartContainer(ctx, imageName, containerName, path, app.ID)

	logger.Info("application updated: %s", app.AppName)
	if a.hub != nil {
		a.hub.BroadcastProgress(100, "Deployment complete!")
	}

	app.StorageLocation = path
	app.TechStackID = &tech.ID
	if err := a.repo.Update(ctx, app); err != nil {
		logger.Error("error updating application: %v", err)
		return "", err
	}

	return path, nil
}

func (a *application) Delete(ctx context.Context, id string) error {
	// Delete associated gateways first
	// gateways, err := a.gateway.GetByApplicationID(ctx, id)
	// if err != nil {
	// 	logger.Error("error getting gateways: %v", err)
	// 	return err
	// }

	// for _, gateway := range gateways {
	// 	if err := a.gateway.Delete(ctx, gateway.ID); err != nil {
	// 		logger.Error("error deleting gateway: %v", err)
	// 		// Continue with other gateways
	// 	}
	// }

	// Get application details
	app, err := a.repo.GetByID(ctx, id)
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

	return a.repo.Delete(ctx, id)
}

func (a *application) StartContainer(ctx context.Context, id string) error {
	app, err := a.repo.GetByID(ctx, id)
	if err != nil {
		logger.Error("error getting application: %v", err)
		return err
	}

	containerId, err := a.docker.GetContainerID(ctx, app.AppName)
	if err != nil {
		logger.Error("error getting container ID: %v", err)
		return err
	}

	return a.docker.StartContainer(ctx, containerId)
}

func (a *application) StopContainer(ctx context.Context, id string) error {
	app, err := a.repo.GetByID(ctx, id)
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
