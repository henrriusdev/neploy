package docker

import (
	"fmt"
	"os"
	"path/filepath"
)

func GenerateDockerfile(projectDir string, stack string) (string, error) {
	tmpl, ok := defaultTemplates[stack]
	if !ok {
		return "", fmt.Errorf("no default template for stack: %s", stack)
	}

	dockerfilePath := filepath.Join(projectDir, "Dockerfile")

	err := WriteDockerfile(dockerfilePath, tmpl)
	if err != nil {
		return "", err
	}

	return dockerfilePath, nil
}

func GetDefaultTemplate(stack string) (DockerfileTemplate, bool) {
	tmpl, ok := defaultTemplates[stack]
	return tmpl, ok
}

func WriteDockerfile(filePath string, tmpl DockerfileTemplate) error {
	return WriteFile(filePath, []byte(tmpl.GetDockerfile()))
}

func WriteFile(filePath string, content []byte) error {
	return os.WriteFile(filePath, content, 0o644)
}
