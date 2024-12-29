package filesystem

import (
	"os"
	"path/filepath"
	"strings"

	"neploy.dev/pkg/websocket"
)

func HasDockerfile(projectDir string, client *websocket.Client) DockerfileStatus {
	status := DockerfileStatus{
		Exists: false,
	}

	// Send progress only if client is not nil
	if client != nil {
		client.SendProgress(0, "Searching for Dockerfile...")
	}

	// First check in root directory
	dockerfilePath := filepath.Join(projectDir, "Dockerfile")
	if FileExists(dockerfilePath) {
		status.Exists = true
		status.Path = dockerfilePath
		if client != nil {
			client.SendProgress(100, "Found Dockerfile in root directory")
		}
		return status
	}

	// Then check in docker directory
	dockerfilePath = filepath.Join(projectDir, "docker", "Dockerfile")
	if FileExists(dockerfilePath) {
		status.Exists = true
		status.Path = dockerfilePath
		if client != nil {
			client.SendProgress(100, "Found Dockerfile in docker directory")
		}
		return status
	}

	// Finally check in .docker directory
	dockerfilePath = filepath.Join(projectDir, ".docker", "Dockerfile")
	if FileExists(dockerfilePath) {
		status.Exists = true
		status.Path = dockerfilePath
		if client != nil {
			client.SendProgress(100, "Found Dockerfile in .docker directory")
		}
		return status
	}

	if client != nil {
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

func DockerfileHasExposedPort(projectDir string) bool {
	dockerfilePath := filepath.Join(projectDir, "Dockerfile")
	if FileExists(dockerfilePath) && FileContains(dockerfilePath, "EXPOSE") {
		return true
	}

	// Check in docker directory
	dockerfilePath = filepath.Join(projectDir, "docker", "Dockerfile")
	if FileExists(dockerfilePath) && FileContains(dockerfilePath, "EXPOSE") {
		return true
	}

	// Check in .docker directory
	dockerfilePath = filepath.Join(projectDir, ".docker", "Dockerfile")
	if FileExists(dockerfilePath) && FileContains(dockerfilePath, "EXPOSE") {
		return true
	}

	return false
}

func FileContains(filePath string, keyword string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	content, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}

	return strings.Contains(string(content), keyword)
}
