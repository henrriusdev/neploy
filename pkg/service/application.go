package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"neploy.dev/config"
	"neploy.dev/pkg/docker"
	"neploy.dev/pkg/filesystem"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
	"neploy.dev/pkg/websocket"
)

type Application interface {
	Create(ctx context.Context, app model.Application, techStack string) (string, error)
	Get(ctx context.Context, id string) (model.Application, error)
	GetAll(ctx context.Context) ([]model.FullApplication, error)
	Update(ctx context.Context, app model.Application) error
	GetStat(ctx context.Context, id string) (model.ApplicationStat, error)
	CreateStat(ctx context.Context, stat model.ApplicationStat) error
	UpdateStat(ctx context.Context, stat model.ApplicationStat) error
	GetHealthy(ctx context.Context) (uint, uint, error)
	Deploy(ctx context.Context, id string, repoURL string)
	Upload(ctx context.Context, id string, file *multipart.FileHeader) (string, error)
	Delete(ctx context.Context, id string) error
}

type application struct {
	repo repository.Application
	stat repository.ApplicationStat
	tech repository.TechStack
	hub  *websocket.Hub
}

func NewApplication(repo repository.Application, stat repository.ApplicationStat, tech repository.TechStack) Application {
	return &application{
		repo: repo,
		stat: stat,
		tech: tech,
		hub:  websocket.GetHub(),
	}
}

func (a *application) Create(ctx context.Context, app model.Application, techStack string) (string, error) {
	tech, err := a.tech.FindOrCreate(ctx, techStack)
	if err != nil {
		logger.Error("error finding or creating tech stack: %v", err)
		return "", err
	}

	app.TechStackID = tech.ID

	return a.repo.Insert(ctx, app)
}

func (a *application) Get(ctx context.Context, id string) (model.Application, error) {
	return a.repo.GetByID(ctx, id)
}

func (a *application) GetAll(ctx context.Context) ([]model.FullApplication, error) {
	apps, err := a.repo.GetAll(ctx)
	if err != nil {
		logger.Error("error getting all applications: %v", err)
		return nil, err
	}

	fullApps := make([]model.FullApplication, len(apps))
	for i, app := range apps {
		stats, err := a.stat.GetByApplicationID(ctx, app.ID)
		if err != nil {
			logger.Error("error getting application stats: %v", err)
			return nil, err
		}

		tech, err := a.tech.GetByID(ctx, app.TechStackID)
		if err != nil {
			logger.Error("error getting tech stack: %v", err)
			return nil, err
		}

		fullApps[i] = model.FullApplication{
			Application: app,
			Stats:       stats,
			TechStack:   tech,
		}
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

func (a *application) Deploy(ctx context.Context, id string, repoURL string) {
	app, err := a.repo.GetByID(ctx, id)
	if err != nil {
		logger.Error("error getting application: %v", err)
		return
	}

	appNameWithoutSpace := strings.ReplaceAll(app.AppName, " ", "-")
	appNameWithoutSpecialChars := regexp.MustCompile(`[^a-zA-Z0-9-]`).ReplaceAllString(appNameWithoutSpace, "")
	appName := strings.ToLower(appNameWithoutSpecialChars)

	path := filepath.Join(config.Env.UploadPath, appName)
	_, err = git.PlainCloneContext(ctx, path, false, &git.CloneOptions{
		URL:        repoURL,
		RemoteName: "origin",
	})
	if err != nil {
		logger.Error("error cloning repo: %v", err)
		return
	}
	logger.Info("repo cloned: %s", path)
	app.StorageLocation = path

	techStack, err := filesystem.DetectStack(path)
	if err != nil {
		logger.Error("error detecting tech stack: %v", err)
		return
	}

	tech, err := a.tech.FindOrCreate(ctx, techStack)
	if err != nil {
		logger.Error("error finding or creating tech stack: %v", err)
		return
	}

	a.hub.BroadcastProgress(0, "Checking for Dockerfile...")

	if !filesystem.HasDockerfile(path, a.hub.GetNotificationClient()).Exists {
		actionInput := websocket.NewSelectInput("action", []string{
			fmt.Sprintf("create_default_for_%s", techStack),
			"skip",
		})

		actionMsg := websocket.NewActionMessage(
			websocket.ActionTypeCritical,
			"Dockerfile Required",
			fmt.Sprintf("No Dockerfile found for %s application. Would you like to create one?", techStack),
			[]websocket.Input{actionInput},
		)

		a.hub.BroadcastInteractive(actionMsg)

		a.hub.BroadcastProgress(50, "Creating Dockerfile...")
		tmpl, ok := docker.GetDefaultTemplate(techStack)
		if !ok {
			logger.Error("no default template for tech stack: %s", techStack)
			a.hub.BroadcastProgress(100, "Error: No default template available for "+techStack)
			return
		}

		dockerfilePath := filepath.Join(path, "Dockerfile")
		if err := docker.WriteDockerfile(dockerfilePath, tmpl); err != nil {
			logger.Error("error writing dockerfile: %v", err)
			a.hub.BroadcastProgress(100, "Error creating Dockerfile")
			return
		}

		a.hub.BroadcastProgress(100, "Created default Dockerfile")
	}

	app.TechStackID = tech.ID
	if err := a.repo.Update(ctx, app); err != nil {
		logger.Error("error updating application: %v", err)
		return
	}

	logger.Info("application updated: %s", app.AppName)
	a.hub.BroadcastProgress(100, "Deployment complete!")
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

	app.StorageLocation = path
	app.TechStackID = tech.ID
	if err := a.repo.Update(ctx, app); err != nil {
		logger.Error("error updating application: %v", err)
		return "", err
	}

	return path, nil
}

func (a *application) Delete(ctx context.Context, id string) error {
	return a.repo.Delete(ctx, id)
}
