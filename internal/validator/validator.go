package validator

import (
	"fmt"
	"os"
	"path/filepath"
)

// ValidateProjectDirectory checks if the target directory is suitable for project initialization
// It creates the directory if it doesn't exist, and ensures it's empty
func ValidateProjectDirectory(projectPath string) error {
	// Convert to absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Check if directory exists
	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		// Directory doesn't exist, create it
		if err := os.MkdirAll(absPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to access directory: %w", err)
	}

	// Check if it's a directory
	if !info.IsDir() {
		return fmt.Errorf("path '%s' exists but is not a directory", absPath)
	}

	// Check if directory is empty
	entries, err := os.ReadDir(absPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	if len(entries) > 0 {
		return fmt.Errorf("directory '%s' is not empty. goinit should be used for creating new projects only", absPath)
	}

	return nil
}
