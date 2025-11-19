package config

import (
	"encoding/json"
	"fmt"
)

// Boilerplate represents a project template configuration
type Boilerplate struct {
	Name               string `json:"name"`
	Description        string `json:"description"`
	ArchiveUrl         string `json:"archiveUrl"`
	ModulePlaceholder  string `json:"modulePlaceholder"`
}

// Config represents the goinit.json configuration
type Config struct {
	AvailableArchitectures []string                `json:"availableArchitectures"`
	Boilerplates           map[string]Boilerplate  `json:"boilerplates"`
}

// Load parses the embedded config bytes and returns a Config struct
func Load(configBytes []byte) (*Config, error) {
	var config Config

	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// GetBoilerplate looks up a boilerplate by architecture name
func (c *Config) GetBoilerplate(architecture string) (*Boilerplate, error) {
	boilerplate, exists := c.Boilerplates[architecture]
	if !exists {
		return nil, fmt.Errorf("architecture '%s' not found in configuration", architecture)
	}

	return &boilerplate, nil
}

// ValidateArchitecture checks if an architecture exists in available architectures
func (c *Config) ValidateArchitecture(architecture string) error {
	for _, arch := range c.AvailableArchitectures {
		if arch == architecture {
			return nil
		}
	}

	return fmt.Errorf("architecture '%s' is not available. Available architectures: %v", architecture, c.AvailableArchitectures)
}
