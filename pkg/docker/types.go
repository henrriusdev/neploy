package docker

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
