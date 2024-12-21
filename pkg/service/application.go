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
}

type application struct {
	repo          repository.Application
	stat          repository.ApplicationStat
	tech          repository.TechStack
	notifications *websocket.Client
	interactive   *websocket.Client
}

func NewApplication(repo repository.Application, stat repository.ApplicationStat, tech repository.TechStack) Application {
	hub := websocket.GetHub()

	return &application{
		repo:          repo,
		stat:          stat,
		tech:          tech,
		notifications: hub.GetNotificationsClient(),
		interactive:   hub.GetInteractiveClient(),
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

	if !filesystem.HasDockerfile(path, nil).Exists {
		// Only proceed with interactive flow if we have a client
		if a.interactive != nil {
			// Ask user what to do via WebSocket
			if err := a.interactive.SendProgress(0, "No Dockerfile found. Waiting for user action..."); err != nil {
				logger.Error("error sending progress: %v", err)
				return
			}

			// Send options to the user
			question := map[string]interface{}{
				"type":    "dockerfile_action",
				"message": fmt.Sprintf("No Dockerfile found. What would you like to do? Your project is using %s.", techStack),
				"options": []string{
					fmt.Sprintf("create_default_for_%s", techStack),
					"skip",
				},
			}

			if err := a.interactive.SendJSON(question); err != nil {
				logger.Error("error sending dockerfile question: %v", err)
				return
			}

			// Wait for user response
			var action struct {
				Type   string `json:"type"`
				Action string `json:"action"`
			}
			if err := a.interactive.ReadJSON(&action); err != nil {
				logger.Error("error reading user response: %v", err)
				return
			}

			expectedAction := fmt.Sprintf("create_default_for_%s", techStack)
			if action.Action == expectedAction {
				if a.notifications != nil {
					if err := a.notifications.SendProgress(50, fmt.Sprintf("Creating default Dockerfile for %s...", techStack)); err != nil {
						logger.Error("error sending progress: %v", err)
					}
				}

				tmpl, ok := docker.GetDefaultTemplate(techStack)
				if !ok {
					logger.Error("no default template for tech stack: %s", techStack)
					if a.notifications != nil {
						if err := a.notifications.SendProgress(100, "Error: No default template available for "+techStack); err != nil {
							logger.Error("error sending progress: %v", err)
						}
					}
					return
				}

				dockerfilePath := filepath.Join(path, "Dockerfile")
				if err := docker.WriteDockerfile(dockerfilePath, tmpl); err != nil {
					logger.Error("error writing dockerfile: %v", err)
					if a.notifications != nil {
						if err := a.notifications.SendProgress(100, "Error creating Dockerfile"); err != nil {
							logger.Error("error sending progress: %v", err)
						}
					}
					return
				}

				if a.notifications != nil {
					if err := a.notifications.SendProgress(100, "Created default Dockerfile"); err != nil {
						logger.Error("error sending progress: %v", err)
					}
				}
			}
		}
	}

	app.TechStackID = tech.ID
	if err := a.repo.Update(ctx, app); err != nil {
		logger.Error("error updating application: %v", err)
		return
	}

	logger.Info("application updated: %s", app.AppName)
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
