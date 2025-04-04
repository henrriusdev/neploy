package service

import (
	"context"
	"fmt"
	"mime/multipart"
	neploker "neploy.dev/pkg/docker"
	"neploy.dev/pkg/filesystem"
	neployway "neploy.dev/pkg/gateway"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/websocket"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Application interface {
	Create(ctx context.Context, app model.Application) (string, error)
	Get(ctx context.Context, id string) (model.ApplicationDockered, error)
	GetAll(ctx context.Context) ([]model.FullApplication, error)
	Update(ctx context.Context, app model.Application) error
	GetStat(ctx context.Context, id string) (model.ApplicationStat, error)
	CreateStat(ctx context.Context, stat model.ApplicationStat) error
	UpdateStat(ctx context.Context, stat model.ApplicationStat) error
	GetHealthy(ctx context.Context) (uint, uint, error)
	Delete(ctx context.Context, id string) error
	StartContainer(ctx context.Context, id string) error
	StopContainer(ctx context.Context, id string) error
	GetRepoBranches(ctx context.Context, repoURL string) ([]string, error)
	Deploy(ctx context.Context, id string, repoURL string, branch string) error
	Upload(ctx context.Context, id string, file *multipart.FileHeader) (string, error)
}

type application struct {
	repos             repository.Repositories
	hub               *websocket.Hub
	docker            *neploker.Docker
	router            *neployway.Router
	versioningService Versioning
	dockerService     Docker
}

func NewApplication(repos repository.Repositories, router *neployway.Router) Application {
	hub := websocket.GetHub()
	dockerClient := neploker.NewDocker()
	return &application{
		repos:             repos,
		hub:               hub,
		docker:            dockerClient,
		router:            router,
		versioningService: NewVersioning(repos, hub, dockerClient),
		dockerService:     NewDocker(repos, hub, dockerClient, router),
	}
}

func (a *application) Create(ctx context.Context, app model.Application) (string, error) {
	return a.repos.Application.Insert(ctx, app)
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
		return 0, 0, err
	}

	var healthy uint
	for _, app := range apps {
		if app.Healthy {
			healthy++
		}
	}
	return healthy, uint(len(apps)), nil
}

func (a *application) StartContainer(ctx context.Context, id string) error {
	return a.dockerService.StartContainer(ctx, id)
}

func (a *application) StopContainer(ctx context.Context, id string) error {
	return a.dockerService.StopContainer(ctx, id)
}

func (a *application) GetRepoBranches(ctx context.Context, repoURL string) ([]string, error) {
	repo := filesystem.NewGitRepo(repoURL)
	return repo.GetBranches()
}

func (a *application) Deploy(ctx context.Context, id string, repoURL string, branch string) error {
	return a.versioningService.Deploy(ctx, id, repoURL, branch)
}

func (a *application) Upload(ctx context.Context, id string, file *multipart.FileHeader) (string, error) {
	return a.versioningService.Upload(ctx, id, file)
}

func (a *application) Get(ctx context.Context, id string) (model.ApplicationDockered, error) {
	app, err := a.repos.Application.GetByID(ctx, id)
	if err != nil {
		logger.Error("error getting app %v", err)
		return model.ApplicationDockered{}, err
	}

	appNameWithoutSpace := strings.ReplaceAll(app.AppName, " ", "-")
	appNameWithoutSpecialChars := regexp.MustCompile(`[^a-zA-Z0-9-]`).ReplaceAllString(appNameWithoutSpace, "")
	appName := strings.ToLower(appNameWithoutSpecialChars)
	containerId, err := a.docker.GetContainerID(ctx, appName)
	if err != nil {
		logger.Error("errpr getting container ID: %v", err)
		return model.ApplicationDockered{}, err
	}

	cpu, ram, err := a.docker.GetUsage(ctx, containerId)
	if err != nil {
		logger.Error("error getting container usage: %v", err)
		return model.ApplicationDockered{}, err
	}
	fmt.Println(cpu, ram)

	uptime, err := a.docker.GetUptime(ctx, containerId)
	if err != nil {
		logger.Error("error getting container uptime: %v", err)
		return model.ApplicationDockered{}, err
	}

	stats, err := a.repos.ApplicationStat.GetByApplicationID(ctx, app.ID)
	if err != nil {
		logger.Error("error getting app stats: %v", err)
		return model.ApplicationDockered{}, err
	}

	var requestsPerMin int
	for _, stat := range stats {
		requestsPerMin += stat.Requests
	}

	logs, err := a.docker.GetLogs(ctx, containerId, false)
	if err != nil {
		logger.Error("error getting container logs: %v", err)
		return model.ApplicationDockered{}, err
	}

	upTime := uptime.String()

	versions, err := a.repos.ApplicationVersion.GetAll(ctx, filters.IsSelectFilter("application_id", app.ID))
	if err != nil {
		logger.Error("error getting application versions: %v", err)
		return model.ApplicationDockered{}, err
	}

	appDockered := model.ApplicationDockered{
		Application:    app,
		CpuUsage:       cpu,
		MemoryUsage:    ram,
		Uptime:         upTime,
		RequestsPerMin: requestsPerMin,
		Logs:           logs,
		Versions:       versions,
	}

	return appDockered, nil
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

		versions, err := a.repos.ApplicationVersion.GetAll(ctx, filters.IsSelectFilter("application_id", app.ID))
		if err != nil {
			logger.Error("error getting application versions: %v", err)
			return nil, err
		}

		var appVersion model.ApplicationVersion
		if len(versions) > 0 {
			appVersion = versions[0]
		} else {
			appVersion, err = a.repos.ApplicationVersion.InsertOne(ctx, model.ApplicationVersion{
				ApplicationID:   app.ID,
				VersionTag:      "v1.0.0",
				Description:     "Initial version",
				Status:          "active",
				StorageLocation: filepath.Join(app.StorageLocation, "v1.0.0"),
			})
			if err != nil {
				logger.Error("error inserting application version: %v", err)
				return nil, err
			}
		}

		appName := strings.ToLower(regexp.MustCompile(`[^a-zA-Z0-9-]`).ReplaceAllString(strings.ReplaceAll(app.AppName, " ", "-"), ""))
		containerName := "neploy-" + appName

		status, err := a.docker.GetContainerStatus(ctx, containerName)
		if err != nil {
			logger.Error("error getting container status: %v", err)
			return nil, err
		}

		if status == "Not created" {
			go func() error {
				dockerfilePath := filepath.Join(appVersion.StorageLocation, "Dockerfile")
				port, err := a.dockerService.ConfigurePort(dockerfilePath, false)
				if err != nil {
					logger.Error("error configuring port: %v", err)
					return err
				}
				return a.dockerService.CreateAndStartContainer(ctx, app, appVersion, port)
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
