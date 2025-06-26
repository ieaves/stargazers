package config

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Repository string `yaml:"repository"`
	GitHub     struct {
		Username string `yaml:"username"`
		Token    string `yaml:"token,omitempty"` // Optional token in config file
	} `yaml:"github"`
}

// Load reads the configuration from config.yaml
func Load() (*Config, error) {
	configPath := "config.yaml"
	
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config.yaml not found. Please create it with your repository and GitHub username")
	}
	
	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config.yaml: %w", err)
	}
	
	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config.yaml: %w", err)
	}
	
	// Validate required fields
	if config.Repository == "" {
		return nil, fmt.Errorf("repository is required in config.yaml")
	}
	if config.GitHub.Username == "" {
		return nil, fmt.Errorf("github.username is required in config.yaml")
	}
	
	// Validate repository URL format
	if err := config.validateRepositoryURL(); err != nil {
		return nil, err
	}
	
	// Load environment variables and .env file
	if err := config.loadEnvironment(); err != nil {
		return nil, err
	}
	
	return &config, nil
}

// loadEnvironment loads environment variables and .env file
func (c *Config) loadEnvironment() error {
	// Load .env file if it exists
	if err := c.loadEnvFile(); err != nil {
		return err
	}
	
	// Load token from GH_TOKEN environment variable
	if token := os.Getenv("GH_TOKEN"); token != "" {
		c.GitHub.Token = token
	}
	
	return nil
}

// loadEnvFile loads environment variables from .env file
func (c *Config) loadEnvFile() error {
	envPath := ".env"
	
	// Check if .env file exists
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return nil // .env file is optional
	}
	
	// Read .env file
	file, err := os.Open(envPath)
	if err != nil {
		return fmt.Errorf("failed to open .env file: %w", err)
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// Parse key=value pairs
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				
				// Remove quotes if present
				value = strings.Trim(value, `"'`)
				
				// Set environment variable
				os.Setenv(key, value)
			}
		}
	}
	
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %w", err)
	}
	
	return nil
}

// GetToken returns the GitHub token from config or environment
func (c *Config) GetToken() string {
	// Priority order:
	// 1. Token from config file
	// 2. GH_TOKEN environment variable
	
	if c.GitHub.Token != "" {
		return c.GitHub.Token
	}
	
	if token := os.Getenv("GH_TOKEN"); token != "" {
		return token
	}
	
	return ""
}

// validateRepositoryURL validates that the repository URL is in the correct format
func (c *Config) validateRepositoryURL() error {
	// Check if it's a GitHub URL
	if !strings.HasPrefix(c.Repository, "https://github.com/") {
		return fmt.Errorf("repository must be a GitHub URL (e.g., https://github.com/owner/repo)")
	}
	
	// Parse the URL to extract owner and repo
	parsedURL, err := url.Parse(c.Repository)
	if err != nil {
		return fmt.Errorf("invalid repository URL: %w", err)
	}
	
	// Check if the path has the correct format (owner/repo)
	pathParts := strings.Split(strings.TrimPrefix(parsedURL.Path, "/"), "/")
	if len(pathParts) != 2 {
		return fmt.Errorf("repository URL must be in format https://github.com/owner/repo")
	}
	
	return nil
}

// GetRepositoryOwner extracts the owner from the repository URL
func (c *Config) GetRepositoryOwner() string {
	parsedURL, _ := url.Parse(c.Repository)
	pathParts := strings.Split(strings.TrimPrefix(parsedURL.Path, "/"), "/")
	if len(pathParts) >= 1 {
		return pathParts[0]
	}
	return ""
}

// GetRepositoryName extracts the repository name from the repository URL
func (c *Config) GetRepositoryName() string {
	parsedURL, _ := url.Parse(c.Repository)
	pathParts := strings.Split(strings.TrimPrefix(parsedURL.Path, "/"), "/")
	if len(pathParts) >= 2 {
		return pathParts[1]
	}
	return ""
}

// GetModulePath returns the full module path for the repository
func (c *Config) GetModulePath() string {
	return fmt.Sprintf("github.com/%s/%s", c.GitHub.Username, c.GetRepositoryName())
}

// GetRepositoryPath returns the repository path in owner/repo format
func (c *Config) GetRepositoryPath() string {
	return fmt.Sprintf("%s/%s", c.GetRepositoryOwner(), c.GetRepositoryName())
}

// GetCloneURL returns the GitHub clone URL
func (c *Config) GetCloneURL() string {
	return fmt.Sprintf("https://github.com/%s/%s.git", c.GitHub.Username, c.GetRepositoryName())
}

// CreateDefaultConfig creates a default config.yaml file
func CreateDefaultConfig() error {
	defaultConfig := `# Stargazers Configuration
# Update these values with your repository and GitHub username

repository: "https://github.com/your-username/your-repo"  # The repository you want to analyze

github:
  username: "your-username"  # Your GitHub username (for module path)
  # token: "your-token"     # Optional: GitHub token (can also use GH_TOKEN environment variable)
`

	return os.WriteFile("config.yaml", []byte(defaultConfig), 0644)
}

// CreateDefaultEnvFile creates a default .env file
func CreateDefaultEnvFile() error {
	defaultEnv := `# GitHub Token (optional - can also be set in config.yaml)
# Uncomment and set your token here:
# GH_TOKEN=your_github_token_here
`

	return os.WriteFile(".env", []byte(defaultEnv), 0644)
} 