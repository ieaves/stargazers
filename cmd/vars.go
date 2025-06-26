package cmd

import (
	"fmt"
	"stargazers/config"
)

// Global variables for command flags
var (
    Repo        string
    AccessToken string
    CacheDir    string
    Mode        string
)

// LoadConfig loads the configuration and sets defaults
func LoadConfig() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Set default repo from config if not provided
	if Repo == "" {
		Repo = cfg.GetRepositoryPath()
	}

	// Set default token from config if not provided
	if AccessToken == "" {
		AccessToken = cfg.GetToken()
	}

	return cfg, nil
}