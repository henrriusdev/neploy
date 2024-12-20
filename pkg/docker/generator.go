package docker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DockerfileTemplate struct {
	BaseImage    string
	WorkDir      string
	Dependencies []string
	BuildCmd     string
	StartCmd     string
}

var defaultTemplates = map[string]DockerfileTemplate{
	"Node.js": {
		BaseImage: "node:18-alpine",
		WorkDir:   "/app",
		Dependencies: []string{
			"COPY package*.json ./",
			"RUN npm install",
		},
		BuildCmd: "npm run build",
		StartCmd: "npm start",
	},
	"Python": {
		BaseImage: "python:3.9-slim",
		WorkDir:   "/app",
		Dependencies: []string{
			"COPY requirements.txt .",
			"RUN pip install -r requirements.txt",
		},
		BuildCmd: "",
		StartCmd: "python app.py",
	},
	"Go": {
		BaseImage: "golang:1.21-alpine",
		WorkDir:   "/app",
		Dependencies: []string{
			"COPY go.mod go.sum ./",
			"RUN go mod download",
		},
		BuildCmd: "go build -o main .",
		StartCmd: "./main",
	},
	// Add more templates for other stacks
}

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
	return WriteFile(filePath, []byte(fmt.Sprintf(`FROM %s
WORKDIR %s
%s
COPY . .
%s
%s
`, tmpl.BaseImage, tmpl.WorkDir, strings.Join(tmpl.Dependencies, "\n"), tmpl.BuildCmd, tmpl.StartCmd)))
}

func WriteFile(filePath string, content []byte) error {
	return os.WriteFile(filePath, content, 0o644)
}
