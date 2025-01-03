package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

// DetectStack scans the directory and determines the main tech stack
func DetectStack(projectDir string) (string, error) {
	detectedStacks := make(map[string]int)

	// Walk through the project directory
	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check by file name
		if stack, ok := StackIndicators[info.Name()]; ok {
			detectedStacks[stack]++
			return nil
		}

		// Check by file extension
		ext := filepath.Ext(info.Name())
		if stack, ok := StackIndicators[ext]; ok {
			detectedStacks[stack]++
			return nil
		}

		// Check MIME type
		mime, err := mimetype.DetectFile(path)
		if err == nil && strings.HasPrefix(mime.String(), "text") {
			// Example: Additional checks for Rust and PHP
			if ext == ".json" && strings.Contains(info.Name(), "composer") {
				detectedStacks["PHP"]++
			}
			if ext == ".toml" && strings.Contains(info.Name(), "Cargo") {
				detectedStacks["Rust"]++
			}
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	// Determine the most probable stack
	var mainStack string
	var maxCount int
	for stack, count := range detectedStacks {
		if count > maxCount {
			mainStack = stack
			maxCount = count
		}
	}

	if mainStack == "" {
		return "", fmt.Errorf("no recognizable tech stack found")
	}

	return mainStack, nil
}
