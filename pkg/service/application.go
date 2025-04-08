package service

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
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

var globalSemaphore = semaphore.NewWeighted(4)

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
	StartContainer(ctx context.Context, id, versionID string) error
	StopContainer(ctx context.Context, id, versionID string) error
	GetRepoBranches(ctx context.Context, repoURL string) ([]string, error)
	Deploy(ctx context.Context, id string, repoURL string, branch string) error
	Upload(ctx context.Context, id string, file *multipart.FileHeader) (string, error)
	DeleteVersion(ctx context.Context, appID string, versionID string) error
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

func (a *application) StartContainer(ctx context.Context, id, versionId string) error {
	return a.dockerService.StartContainer(ctx, id, versionId)
}

func (a *application) StopContainer(ctx context.Context, id, versionId string) error {
	return a.dockerService.StopContainer(ctx, id, versionId)
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

func (a *application) ensureContainerRunning(ctx context.Context, app model.Application, version model.ApplicationVersion) error {
	if version.Status == "paused" || version.Status == "inactive" {
		return nil
	}

	if err := globalSemaphore.Acquire(ctx, 1); err != nil {
		logger.Error("semaphore timeout for version %s: %v", version.VersionTag, err)
		return err
	}
	defer globalSemaphore.Release(1)

	containerName := getContainerName(app.AppName, version.VersionTag)
	status, err := a.docker.GetContainerStatus(ctx, containerName)
	if err != nil {
		logger.Error("error getting container status: %v", err)
		return err
	}

	if status == "Not created" {
		dockerfilePath := filepath.Join(version.StorageLocation, "Dockerfile")
		port, err := a.dockerService.ConfigurePort(dockerfilePath, false)
		if err != nil {
			logger.Error("error configuring port: %v", err)
			return err
		}
		return a.dockerService.CreateAndStartContainer(ctx, app, version, port)
	} else if status == "Stopped" {
		return a.dockerService.StartContainer(ctx, app.ID, version.VersionTag)
	}

	return nil
}

func (a *application) Get(ctx context.Context, id string) (model.ApplicationDockered, error) {
	app, err := a.repos.Application.GetByID(ctx, id)
	if err != nil {
		logger.Error("error getting app %v", err)
		return model.ApplicationDockered{}, err
	}

	stats, err := a.repos.ApplicationStat.GetByApplicationID(ctx, app.ID)
	if err != nil {
		logger.Error("error getting app stats: %v", err)
		return model.ApplicationDockered{}, err
	}

	versions, err := a.repos.ApplicationVersion.GetAll(ctx, filters.IsSelectFilter("application_id", app.ID))
	if err != nil {
		logger.Error("error getting application versions: %v", err)
		return model.ApplicationDockered{}, err
	}

	for _, version := range versions {
		v := version
		go func() {
			if err := a.ensureContainerRunning(ctx, app, v); err != nil {
				logger.Error("error ensuring container for version %s: %v", v.VersionTag, err)
			}
		}()
	}

	if len(versions) == 0 {
		return model.ApplicationDockered{
			Application:    app,
			CpuUsage:       0,
			MemoryUsage:    0,
			Uptime:         "0s",
			RequestsPerMin: 0,
			Logs:           []string{},
			Versions:       versions,
		}, nil
	}

	latest := versions[len(versions)-1]
	containerID, err := a.docker.GetContainerID(ctx, getContainerName(app.AppName, latest.VersionTag))
	if err != nil {
		logger.Error("error getting container ID: %v", err)
		return model.ApplicationDockered{}, err
	}

	cpu, ram, err := a.docker.GetUsage(ctx, containerID)
	if err != nil {
		logger.Error("error getting container usage: %v", err)
		return model.ApplicationDockered{}, err
	}

	uptime, err := a.docker.GetUptime(ctx, containerID)
	if err != nil {
		logger.Error("error getting container uptime: %v", err)
		return model.ApplicationDockered{}, err
	}

	logs, err := a.docker.GetLogs(ctx, containerID, false)
	if err != nil {
		logger.Error("error getting container logs: %v", err)
		return model.ApplicationDockered{}, err
	}

	var requestsPerMin int
	for _, stat := range stats {
		requestsPerMin += stat.Requests
	}

	return model.ApplicationDockered{
		Application:    app,
		CpuUsage:       cpu,
		MemoryUsage:    ram,
		Uptime:         uptime.String(),
		RequestsPerMin: requestsPerMin,
		Logs:           logs,
		Versions:       versions,
	}, nil
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

		for _, version := range versions {
			v := version
			go func() {
				if err := a.ensureContainerRunning(context.Background(), app, v); err != nil {
					logger.Error("error ensuring container for version %s: %v", v.VersionTag, err)
				}
			}()
		}

		// Usar el status de la primera versión activa, si hay
		appStatus := "unknown"
		for _, v := range versions {
			if v.Status != "paused" && v.Status != "inactive" {
				containerName := getContainerName(app.AppName, v.VersionTag)
				status, err := a.docker.GetContainerStatus(ctx, containerName)
				if err == nil {
					appStatus = status
					break
				}
			}
		}

		fullApps = append(fullApps, model.FullApplication{
			Application: app,
			TechStack:   tech,
			Stats:       stats,
			Status:      appStatus,
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

func (a *application) DeleteVersion(ctx context.Context, appID string, versionID string) error {
	return a.versioningService.DeleteVersion(ctx, appID, versionID)
}
