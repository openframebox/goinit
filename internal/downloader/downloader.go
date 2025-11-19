package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Download fetches an archive from the given URL and saves it to a temporary file
// Returns the path to the downloaded file
func Download(url string) (string, error) {
	fmt.Printf("Downloading template from %s...\n", url)

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "goinit-*.tar.gz")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tmpFile.Close()

	tmpPath := tmpFile.Name()

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		os.Remove(tmpPath) // Clean up on error
		return "", fmt.Errorf("failed to download archive: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		os.Remove(tmpPath) // Clean up on error
		return "", fmt.Errorf("failed to download archive: HTTP %d %s", resp.StatusCode, resp.Status)
	}

	// Write response body to file
	written, err := io.Copy(tmpFile, resp.Body)
	if err != nil {
		os.Remove(tmpPath) // Clean up on error
		return "", fmt.Errorf("failed to save archive: %w", err)
	}

	fmt.Printf("Downloaded %d bytes\n", written)

	return tmpPath, nil
}

// DetectArchiveFormat determines the archive format from the URL
func DetectArchiveFormat(url string) string {
	ext := filepath.Ext(url)

	// Handle .tar.gz case
	if ext == ".gz" && len(url) > 7 && url[len(url)-7:] == ".tar.gz" {
		return "tar.gz"
	}

	// Handle .zip case
	if ext == ".zip" {
		return "zip"
	}

	// Default to tar.gz
	return "tar.gz"
}
