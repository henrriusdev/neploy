package filesystem

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"neploy.dev/pkg/logger"
)

// GitRepo represents a Git repository
type GitRepo struct {
	URL     string
	Branch  string
	BaseDir string
}

// NewGitRepo creates a new GitRepo instance
func NewGitRepo(url string) *GitRepo {
	return &GitRepo{
		URL: url,
	}
}

// GetBranches returns a list of available branches in the repository
func (g *GitRepo) GetBranches() ([]string, error) {
	// Use git ls-remote to list remote branches
	cmd := exec.Command("git", "ls-remote", "--heads", g.URL)
	output, err := cmd.Output()
	if err != nil {
		logger.Error("error listing remote branches: %v", err)
		return nil, err
	}

	// Parse the output to get branch names
	branches := []string{}
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		// Format is: <commit-hash>\trefs/heads/<branch-name>
		parts := strings.Split(line, "\t")
		if len(parts) == 2 {
			branchName := strings.TrimPrefix(parts[1], "refs/heads/")
			branches = append(branches, branchName)
		}
	}

	return branches, nil
}

// Clone clones the repository to a local directory
func (g *GitRepo) Clone(destDir string, branch string) error {
	if branch == "" {
		return errors.New("branch is required")
	}

	// Remove the directory if it exists
	if _, err := os.Stat(destDir); err == nil {
		if err := os.RemoveAll(destDir); err != nil {
			logger.Error("error removing existing directory: %v", err)
			return fmt.Errorf("failed to remove existing directory: %v", err)
		}
	}

	// Create the destination directory
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		logger.Error("error creating destination directory: %v", err)
		return err
	}

	// Clone the repository
	cmd := exec.Command("git", "clone", "-b", branch, "--single-branch", g.URL, destDir)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		errMsg := stderr.String()
		logger.Error("error cloning repository: %v - %s", err, errMsg)
		return fmt.Errorf("git clone failed: %v - %s", err, errMsg)
	}

	g.BaseDir = destDir
	g.Branch = branch
	return nil
}

// CleanUp removes the cloned repository
func (g *GitRepo) CleanUp() error {
	if g.BaseDir != "" {
		return os.RemoveAll(g.BaseDir)
	}
	return nil
}
