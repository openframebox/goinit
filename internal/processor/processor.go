package processor

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ProcessAndCopy processes all files in the source directory, replacing module placeholders,
// and copies them to the destination directory
func ProcessAndCopy(srcDir, destDir, modulePlaceholder, moduleName string) error {
	fmt.Println("Processing template files...")
	fmt.Println()

	// Compile regex pattern for placeholder replacement
	pattern := regexp.MustCompile(regexp.QuoteMeta(modulePlaceholder))

	fileCount := 0
	replacementCount := 0

	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path: %w", err)
		}

		// Skip root directory
		if relPath == "." {
			return nil
		}

		// Skip hidden directories (like .git)
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			fmt.Printf("  Skipping hidden directory: %s\n", relPath)
			return filepath.SkipDir
		}

		destPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			// Create directory
			fmt.Printf("  Creating directory: %s\n", relPath)
			return os.MkdirAll(destPath, info.Mode())
		}

		// Process file
		fmt.Printf("  Processing file: %s", relPath)
		replaced, err := processFile(path, destPath, pattern, moduleName)
		if err != nil {
			fmt.Println()
			return fmt.Errorf("failed to process file %s: %w", relPath, err)
		}

		if replaced {
			fmt.Printf(" (replaced module name)\n")
			replacementCount++
		} else {
			fmt.Println()
		}

		fileCount++
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("✓ Processed %d files (%d files with module name replacements)\n", fileCount, replacementCount)
	return nil
}

// processFile reads a file, replaces placeholders, and writes to destination
// Returns true if replacements were made, false otherwise
func processFile(srcPath, destPath string, pattern *regexp.Regexp, replacement string) (bool, error) {
	// Read source file
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return false, fmt.Errorf("failed to read file: %w", err)
	}

	// Check if file is likely binary
	if isBinary(content) {
		// For binary files, just copy without modification
		return false, copyFile(srcPath, destPath)
	}

	// Check if pattern exists in content
	hasMatch := pattern.Match(content)

	// Replace placeholders in text files
	modified := pattern.ReplaceAll(content, []byte(replacement))

	// Get source file info for permissions
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return false, fmt.Errorf("failed to stat source file: %w", err)
	}

	// Write to destination with same permissions
	if err := os.WriteFile(destPath, modified, srcInfo.Mode()); err != nil {
		return false, fmt.Errorf("failed to write file: %w", err)
	}

	return hasMatch, nil
}

// copyFile copies a file from src to dest preserving permissions
func copyFile(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
	}

	destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

// isBinary checks if content appears to be binary data
// Uses a simple heuristic: if there are null bytes or high percentage of non-printable chars, it's binary
func isBinary(content []byte) bool {
	if len(content) == 0 {
		return false
	}

	// Check first 512 bytes (or less)
	checkSize := len(content)
	if checkSize > 512 {
		checkSize = 512
	}

	sample := content[:checkSize]

	// Check for null bytes
	for _, b := range sample {
		if b == 0 {
			return true
		}
	}

	// Count non-printable characters (excluding common whitespace)
	nonPrintable := 0
	for _, b := range sample {
		if b < 32 && b != '\n' && b != '\r' && b != '\t' {
			nonPrintable++
		}
	}

	// If more than 10% non-printable, consider it binary
	return float64(nonPrintable)/float64(checkSize) > 0.1
}
