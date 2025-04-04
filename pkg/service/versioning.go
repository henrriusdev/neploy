package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	gitconfig "gopkg.in/src-d/go-git.v4/config"

	"neploy.dev/config"
	neploker "neploy.dev/pkg/docker"
	"neploy.dev/pkg/filesystem"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
	"neploy.dev/pkg/websocket"
)

type Versioning interface {
	Deploy(ctx context.Context, id string, repoURL string, branch string) error
	Upload(ctx context.Context, id string, file *multipart.FileHeader) (string, error)
}

type versioning struct {
	repos  repository.Repositories
	hub    *websocket.Hub
	docker *neploker.Docker
}

func NewVersioning(repos repository.Repositories, hub *websocket.Hub, docker *neploker.Docker) Versioning {
	return &versioning{repos, hub, docker}
}

func (v *versioning) Deploy(ctx context.Context, id string, repoURL string, branch string) error {
	app, err := v.repos.Application.GetByID(ctx, id)
	if err != nil {
		logger.Error("error getting application: %v", err)
		return err
	}

	appName := sanitizeAppName(app.AppName)
	basePath := filepath.Join(config.Env.UploadPath, appName)

	versionTag, err := getLatestGitTag(repoURL)
	if err != nil {
		logger.Error("error fetching tags for app %s, using default tag: %s", app.AppName, err)
		versionTag = "v1.0.0"
	}
	logger.Info("detected latest tag: %s", versionTag)

	versionPath := filepath.Join(basePath, versionTag)
	if err := os.MkdirAll(versionPath, os.ModePerm); err != nil {
		logger.Error("error creating version directory: %v", err)
		return err
	}

	repo := filesystem.NewGitRepo(repoURL)
	if err := repo.Clone(versionPath, branch); err != nil {
		logger.Error("error cloning repository: %v", err)
		return err
	}

	techStack, err := filesystem.DetectStack(versionPath)
	if err != nil {
		logger.Error("error detecting tech stack: %v", err)
		return err
	}

	tech, err := v.repos.TechStack.FindOrCreate(ctx, techStack)
	if err != nil {
		logger.Error("error finding or creating tech stack: %v", err)
		return err
	}

	if v.hub != nil {
		v.hub.BroadcastProgress(0, "Checking for Dockerfile...")
	}

	var wsClient *websocket.Client
	if v.hub != nil {
		wsClient = v.hub.GetNotificationClient()
	}

	dockerStatus := filesystem.HasDockerfile(versionPath, wsClient)
	if !dockerStatus.Exists {
		if v.hub != nil {
			v.hub.BroadcastProgress(50, "Creating Dockerfile...")
		}

		tmpl, ok := neploker.GetDefaultTemplate(techStack)
		if !ok {
			logger.Error("no default template for tech stack: %s", techStack)
			return fmt.Errorf("no template for tech stack")
		}

		dockerfilePath := filepath.Join(versionPath, "Dockerfile")
		if err := neploker.WriteDockerfile(dockerfilePath, tmpl); err != nil {
			logger.Error("error writing dockerfile: %v", err)
			return err
		}
	}

	app.TechStackID = &tech.ID
	if err := v.repos.Application.Update(ctx, app); err != nil {
		logger.Error("error updating application: %v", err)
		return err
	}

	version := model.ApplicationVersion{
		VersionTag:      versionTag,
		Description:     fmt.Sprintf("Auto deployed from repo (%s)", versionTag),
		Status:          "active",
		StorageLocation: versionPath,
		ApplicationID:   app.ID,
	}
	if err := v.repos.ApplicationVersion.Insert(ctx, version); err != nil {
		logger.Error("error inserting application version: %v", err)
		return err
	}

	logger.Info("application deployed: %s - version %s", app.AppName, versionTag)
	if v.hub != nil {
		v.hub.BroadcastProgress(100, "Deployment complete!")
	}
	return nil
}

func (v *versioning) Upload(ctx context.Context, id string, file *multipart.FileHeader) (string, error) {
	app, err := v.repos.Application.GetByID(ctx, id)
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

	tech, err := v.repos.TechStack.FindOrCreate(ctx, techStack)
	if err != nil {
		logger.Error("error finding or creating tech stack: %v", err)
		return "", err
	}

	if err := os.Remove(zipPath); err != nil {
		logger.Error("error deleting zip file: %v", err)
	}

	if v.hub != nil {
		v.hub.BroadcastProgress(0, "Checking for Dockerfile...")
	}

	dockerStatus := filesystem.HasDockerfile(path, v.hub.GetNotificationClient())
	if !dockerStatus.Exists {
		tmpl, ok := neploker.GetDefaultTemplate(techStack)
		if !ok {
			logger.Error("no default template for tech stack: %s", techStack)
			return "", fmt.Errorf("no template for tech stack")
		}
		dockerfilePath := filepath.Join(path, "Dockerfile")
		if err := neploker.WriteDockerfile(dockerfilePath, tmpl); err != nil {
			logger.Error("error writing dockerfile: %v", err)
			return "", err
		}
	}

	app.TechStackID = &tech.ID
	app.StorageLocation = path
	if err := v.repos.Application.Update(ctx, app); err != nil {
		logger.Error("error updating application: %v", err)
		return "", err
	}

	if v.hub != nil {
		v.hub.BroadcastProgress(100, "Deployment complete!")
	}

	return path, nil
}

func (v *versioning) DeleteVersion(ctx context.Context, appID string, versionID string) error {
	version, err := v.repos.ApplicationVersion.GetOneById(ctx, versionID)
	if err != nil {
		return fmt.Errorf("version not found: %w", err)
	}
	if version.ApplicationID != appID {
		return fmt.Errorf("version does not belong to the application")
	}
	if version.StorageLocation != "" {
		if err := os.RemoveAll(version.StorageLocation); err != nil {
			logger.Error("failed to delete storage folder: %v", err)
			// not fatal, continue
		}
	}
	return v.repos.ApplicationVersion.Delete(ctx, versionID)
}

func getLatestGitTag(repoURL string) (string, error) {
	repo := git.NewRemote(nil, &gitconfig.RemoteConfig{
		Name: "origin",
		URLs: []string{repoURL},
	})

	refs, err := repo.List(&git.ListOptions{})
	if err != nil {
		return "", err
	}

	var tags []string
	for _, ref := range refs {
		if ref.Name().IsTag() {
			tags = append(tags, ref.Name().Short())
		}
	}

	if len(tags) == 0 {
		return "", fmt.Errorf("no valid tags found")
	}

	return tags[len(tags)-1], nil // último tag asumido más reciente
}

func sanitizeAppName(name string) string {
	safe := strings.ReplaceAll(name, " ", "-")
	safe = regexp.MustCompile(`[^a-zA-Z0-9-]`).ReplaceAllString(safe, "")
	return strings.ToLower(safe)
}
