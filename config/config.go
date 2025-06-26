package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Repository string `yaml:"repository"`
	GitHub     struct {
		Username string `yaml:"username"`
		Token    string `yaml:"token,omitempty"`
	} `yaml:"github"`
}

func Load() (*Config, error) {
	configPath := "config.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config.yaml not found. Please create it with your repository and GitHub username")
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config.yaml: %w", err)
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config.yaml: %w", err)
	}
	if config.Repository == "" {
		return nil, fmt.Errorf("repository is required in config.yaml")
	}
	if config.GitHub.Username == "" {
		return nil, fmt.Errorf("github.username is required in config.yaml")
	}
	if err := config.validateRepositoryURL(); err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) GetToken() string {
	if c.GitHub.Token != "" {
		return c.GitHub.Token
	}
	if token := os.Getenv("GH_TOKEN"); token != "" {
		return token
	}
	return ""
}

func (c *Config) validateRepositoryURL() error {
	if !strings.HasPrefix(c.Repository, "https://github.com/") {
		return fmt.Errorf("repository must be a GitHub URL (e.g., https://github.com/owner/repo)")
	}
	parsedURL, err := url.Parse(c.Repository)
	if err != nil {
		return fmt.Errorf("invalid repository URL: %w", err)
	}
	pathParts := strings.Split(strings.TrimPrefix(parsedURL.Path, "/"), "/")
	if len(pathParts) != 2 {
		return fmt.Errorf("repository URL must be in format https://github.com/owner/repo")
	}
	return nil
}

func (c *Config) GetRepositoryOwner() string {
	parsedURL, _ := url.Parse(c.Repository)
	pathParts := strings.Split(strings.TrimPrefix(parsedURL.Path, "/"), "/")
	if len(pathParts) >= 1 {
		return pathParts[0]
	}
	return ""
}

func (c *Config) GetRepositoryName() string {
	parsedURL, _ := url.Parse(c.Repository)
	pathParts := strings.Split(strings.TrimPrefix(parsedURL.Path, "/"), "/")
	if len(pathParts) >= 2 {
		return pathParts[1]
	}
	return ""
}

func (c *Config) GetModulePath() string {
	return fmt.Sprintf("github.com/%s/%s", c.GitHub.Username, c.GetRepositoryName())
}

func (c *Config) GetRepositoryPath() string {
	return fmt.Sprintf("%s/%s", c.GetRepositoryOwner(), c.GetRepositoryName())
}

func (c *Config) GetCloneURL() string {
	return fmt.Sprintf("https://github.com/%s/%s.git", c.GitHub.Username, c.GetRepositoryName())
}

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

func CreateDefaultEnvFile() error {
	defaultEnv := `# GitHub Token (optional - can also be set in config.yaml)
# Uncomment and set your token here:
# GH_TOKEN=your_github_token_here
`
	return os.WriteFile(".env", []byte(defaultEnv), 0644)
}
