package filesystem

import (
	"os"
	"path/filepath"

	"neploy.dev/pkg/websocket"
)

type DockerfileStatus struct {
	Exists bool   `json:"exists"`
	Path   string `json:"path,omitempty"`
}

func HasDockerfile(projectDir string, client *websocket.Client) DockerfileStatus {
	status := DockerfileStatus{
		Exists: false,
	}

	// Send initial progress
	client.SendProgress(0, "Searching for Dockerfile...")

	// First check in root directory
	dockerfilePath := filepath.Join(projectDir, "Dockerfile")
	if FileExists(dockerfilePath) {
		status.Exists = true
		status.Path = dockerfilePath
		client.SendProgress(100, "Found Dockerfile in root directory")
		return status
	}

	client.SendProgress(30, "Searching in subdirectories...")

	// Then search in subdirectories
	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "Dockerfile" {
			status.Exists = true
			status.Path = path
			client.SendProgress(100, "Found Dockerfile in subdirectory")
			return filepath.SkipDir
		}

		return nil
	})
	if err != nil {
		client.SendProgress(100, "Error searching for Dockerfile")
		return status
	}

	if !status.Exists {
		client.SendProgress(100, "No Dockerfile found")
	}

	return status
}

func HasDockerCompose(projectDir string) bool {
	return FileExists(filepath.Join(projectDir, "docker-compose.yml"))
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
