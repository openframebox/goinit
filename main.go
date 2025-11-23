package main

import (
	"fmt"
	"os"
	"strings"

	_ "embed"

	"github.com/openframebox/goinit/internal/config"
	"github.com/openframebox/goinit/internal/downloader"
	"github.com/openframebox/goinit/internal/extractor"
	"github.com/openframebox/goinit/internal/processor"
	"github.com/openframebox/goinit/internal/validator"
	"github.com/spf13/cobra"
)

//go:embed goinit.json
var configBytes []byte

const version = "1.1.1"

func main() {
	cmd := &cobra.Command{
		Use:   "goinit",
		Short: "Initialize a new project",
		Long:  `Initialize a new project`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	cmd.AddCommand(createCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(versionCmd())

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new project",
		Long:  `Create a new project from a template architecture`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get flags
			projectDir, err := cmd.Flags().GetString("project")
			if err != nil {
				return err
			}
			moduleName, err := cmd.Flags().GetString("module")
			if err != nil {
				return err
			}
			architecture, err := cmd.Flags().GetString("architecture")
			if err != nil {
				return err
			}

			// Load configuration
			cfg, err := config.Load(configBytes)
			if err != nil {
				return fmt.Errorf("failed to load configuration: %w", err)
			}

			// Validate architecture exists
			if err := cfg.ValidateArchitecture(architecture); err != nil {
				return err
			}

			// Get boilerplate configuration
			boilerplate, err := cfg.GetBoilerplate(architecture)
			if err != nil {
				return err
			}

			fmt.Printf("Creating new project: %s\n", moduleName)
			fmt.Printf("Architecture: %s\n", boilerplate.Name)
			fmt.Printf("Target directory: %s\n", projectDir)
			fmt.Println()

			// Validate project directory
			if err := validator.ValidateProjectDirectory(projectDir); err != nil {
				return err
			}

			// Download archive
			archivePath, err := downloader.Download(boilerplate.ArchiveUrl)
			if err != nil {
				return err
			}
			defer os.Remove(archivePath) // Clean up downloaded archive

			// Detect archive format
			format := downloader.DetectArchiveFormat(boilerplate.ArchiveUrl)

			// Extract archive
			extractedDir, err := extractor.Extract(archivePath, format)
			if err != nil {
				return err
			}
			defer os.RemoveAll(extractedDir) // Clean up extracted files

			// Process and copy files to target directory
			if err := processor.ProcessAndCopy(extractedDir, projectDir, boilerplate.ModulePlaceholder, moduleName); err != nil {
				return err
			}

			fmt.Println()
			fmt.Printf("✓ Project created successfully at %s\n", projectDir)
			fmt.Println()
			fmt.Println("Next steps:")
			fmt.Printf("  cd %s\n", projectDir)
			fmt.Println("  go mod tidy")

			return nil
		},
	}

	cmd.Flags().StringP("project", "p", "", "project directory path")
	cmd.MarkFlagRequired("project")
	cmd.Flags().StringP("module", "m", "", "module name")
	cmd.MarkFlagRequired("module")
	cmd.Flags().StringP("architecture", "a", "", "architecture name (e.g., layered@v0)")
	cmd.MarkFlagRequired("architecture")

	return cmd
}

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available architectures",
		Long:  `List all available project architectures and their details`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load configuration
			cfg, err := config.Load(configBytes)
			if err != nil {
				return fmt.Errorf("failed to load configuration: %w", err)
			}

			fmt.Println("Available Architectures:")
			fmt.Println()

			// Display each available architecture
			for i, arch := range cfg.AvailableArchitectures {
				boilerplate, err := cfg.GetBoilerplate(arch)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: Could not load details for %s: %v\n", arch, err)
					continue
				}

				// Architecture name with visual separator
				fmt.Printf("┌─ %s\n", arch)
				fmt.Printf("│  Name: %s\n", boilerplate.Name)
				fmt.Printf("│  Description:\n")

				// Word wrap the description at 70 characters
				wrapped := wrapText(boilerplate.Description, 70)
				for _, line := range wrapped {
					fmt.Printf("│    %s\n", line)
				}

				// Add separator between items, but not after the last one
				if i < len(cfg.AvailableArchitectures)-1 {
					fmt.Println("│")
				} else {
					fmt.Println("└─")
				}
			}

			fmt.Println()
			fmt.Println("Usage:")
			fmt.Printf("  goinit create --project <dir> --module <name> --architecture <arch>\n")

			return nil
		},
	}

	return cmd
}

// wrapText wraps text at the specified width
func wrapText(text string, width int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}

	var lines []string
	currentLine := words[0]

	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	lines = append(lines, currentLine)
	return lines
}

func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  `Display the current version of goinit`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("goinit version %s\n", version)
		},
	}

	return cmd
}
