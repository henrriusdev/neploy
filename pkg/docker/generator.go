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
	StartCmd     []string
	Port         string
	EnvVars      map[string]string
}

var defaultTemplates = map[string]DockerfileTemplate{
	"Node": {
		BaseImage: "node:18-alpine",
		WorkDir:   "/app",
		Dependencies: []string{
			"COPY package*.json ./",
			"RUN npm install",
		},
		BuildCmd: "npm run build",
		StartCmd: []string{"npm", "start"},
		Port:     "3000",
		EnvVars: map[string]string{
			"NODE_ENV": "production",
			"PORT":     "3000",
		},
	},
	"Python": {
		BaseImage: "python:3.9-slim",
		WorkDir:   "/app",
		Dependencies: []string{
			"COPY requirements.txt .",
			"RUN pip install -r requirements.txt",
		},
		BuildCmd: "",
		StartCmd: []string{"python", "app.py"},
		Port:     "5000",
		EnvVars: map[string]string{
			"FLASK_ENV": "production",
			"PORT":      "5000",
		},
	},
	"Go": {
		BaseImage: "golang:1.23-alpine",
		WorkDir:   "/app",
		Dependencies: []string{
			"RUN ls -la",
			"COPY go.mod go.sum ./",
			"RUN go mod download",
			"COPY . .",
		},
		BuildCmd: "go build -o main .",
		StartCmd: []string{"./main"},
		Port:     "8080",
		EnvVars: map[string]string{
			"GO_ENV": "production",
			"PORT":   "8080",
		},
	},
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
	// Convert start command slice to properly formatted JSON array string
	cmdArgs := make([]string, len(tmpl.StartCmd))
	for i, arg := range tmpl.StartCmd {
		cmdArgs[i] = fmt.Sprintf("%q", arg)
	}

	// Build environment variables
	envVars := make([]string, 0, len(tmpl.EnvVars))
	for key, value := range tmpl.EnvVars {
		envVars = append(envVars, fmt.Sprintf("ENV %s=%s", key, value))
	}

	dockerfile := fmt.Sprintf(`FROM %s
WORKDIR %s
%s
%s
COPY . .
RUN %s
EXPOSE %s
%s
CMD [%s]
`,
		tmpl.BaseImage,
		tmpl.WorkDir,
		strings.Join(envVars, "\n"),
		strings.Join(tmpl.Dependencies, "\n"),
		tmpl.BuildCmd,
		tmpl.Port,
		"# Generated by Neploy",
		strings.Join(cmdArgs, ", "))

	return WriteFile(filePath, []byte(dockerfile))
}

func WriteFile(filePath string, content []byte) error {
	return os.WriteFile(filePath, content, 0o644)
}
