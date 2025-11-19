package extractor

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Extract extracts an archive to a temporary directory
// Returns the path to the extracted content (with root directory stripped if applicable)
func Extract(archivePath, format string) (string, error) {
	fmt.Println("Extracting template...")

	// Create temporary directory for extraction
	tmpDir, err := os.MkdirTemp("", "goinit-extract-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}

	switch format {
	case "tar.gz":
		if err := extractTarGz(archivePath, tmpDir); err != nil {
			os.RemoveAll(tmpDir)
			return "", err
		}
	case "zip":
		if err := extractZip(archivePath, tmpDir); err != nil {
			os.RemoveAll(tmpDir)
			return "", err
		}
	default:
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("unsupported archive format: %s", format)
	}

	// Strip root directory if all files are in a single subdirectory (common for GitHub archives)
	strippedDir, err := stripRootDirectory(tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}

	fmt.Println("Extraction complete")
	return strippedDir, nil
}

// extractTarGz extracts a .tar.gz archive to the destination directory
func extractTarGz(archivePath, destDir string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		target := filepath.Join(destDir, header.Name)

		// Security check: prevent path traversal
		if !strings.HasPrefix(target, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path in archive: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			// Create parent directories
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}

			// Create file
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}

			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to write file: %w", err)
			}

			outFile.Close()
		}
	}

	return nil
}

// extractZip extracts a .zip archive to the destination directory
func extractZip(archivePath, destDir string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open zip archive: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		target := filepath.Join(destDir, f.Name)

		// Security check: prevent path traversal
		if !strings.HasPrefix(target, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path in archive: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}

		// Create file
		outFile, err := os.Create(target)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("failed to open file in archive: %w", err)
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()

		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
	}

	return nil
}

// stripRootDirectory checks if all files are in a single subdirectory and returns that path
// This is common for GitHub archives which have a root directory like "repo-name-v0.0.0/"
func stripRootDirectory(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("failed to read extracted directory: %w", err)
	}

	// If there's exactly one entry and it's a directory, use it as the root
	if len(entries) == 1 && entries[0].IsDir() {
		return filepath.Join(dir, entries[0].Name()), nil
	}

	// Otherwise, return the original directory
	return dir, nil
}
