package filesystem

// DockerfileStatus represents the status of a Dockerfile in a project
type DockerfileStatus struct {
	Exists bool   `json:"exists"`
	Path   string `json:"path,omitempty"`
}

// Map of indicators to tech stacks
var StackIndicators = map[string]string{
	".vue":             "Vue",
	".js":              "JavaScript",
	".ts":              "TypeScript",
	".jsx":             "React",
	".tsx":             "React",
	".svelte":          "Svelte",
	"app.module.ts":    "Angular",
	"requirements.txt": "Python",
	"go.mod":           "Go",
	"pom.xml":          "Java",
	"build.gradle":     "Java",
	"Cargo.toml":       "Rust",
	"composer.json":    "PHP",
	"server.js":        "Node.js",
	"app.js":           "Node.js",
	"index.js":         "Node.js",
	".py":              "Python",
	".go":              "Go",
	".java":            "Java",
	".rs":              "Rust",
	".php":             "PHP",
}
