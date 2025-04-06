package git

import (
	"fmt"
	"os"
	"os/exec"
)

// CloneRepository clones a Git repository into a temporary directory and returns the path
func CloneRepository(repoURL, branch string) (string, error) {
	// Create a temporary directory for the clone
	tempDir, err := os.MkdirTemp("", "git-clone-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}

	// Prepare clone arguments
	args := []string{"clone", "--depth=1"}

	// Add branch if specified
	if branch != "" {
		args = append(args, "-b", branch)
	}

	// Add repository URL and destination
	args = append(args, repoURL, tempDir)

	// Execute git clone command
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Clean up the temporary directory if clone fails
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to clone repository: %w", err)
	}

	return tempDir, nil
}
