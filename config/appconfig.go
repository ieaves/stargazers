package config

import (
	"fmt"
	"os"
	"strings"
)

type AppConfig struct {
	Repo     string
	Token    string
	CacheDir string
	Mode     string
	Username string
}

// LoadAppConfig merges config.yaml, env, and CLI flags (in that order of increasing priority)
func LoadAppConfig(cliRepo, cliToken, cliCacheDir, cliMode, cliUsername string) (*AppConfig, error) {
	// 1. Load config.yaml (lowest priority)
	var fileCfg *Config
	fileCfg, _ = Load() // ignore error, treat as empty if missing/invalid

	cfg := &AppConfig{}
	if fileCfg != nil {
		cfg.Username = fileCfg.GitHub.Username
		cfg.CacheDir = "./stargazer_cache"
		if fileCfg.GitHub.Token != "" {
			cfg.Token = fileCfg.GitHub.Token
		}
		if fileCfg.Repository != "" {
			cfg.Repo = fileCfg.GetRepositoryPath()
		}
	}

	// 2. Overwrite with env vars if set
	if v := os.Getenv("GH_TOKEN"); v != "" {
		cfg.Token = v
	}
	if v := os.Getenv("STARGAZERS_REPO"); v != "" {
		cfg.Repo = v
	}
	if v := os.Getenv("STARGAZERS_CACHE"); v != "" {
		cfg.CacheDir = v
	}
	if v := os.Getenv("STARGAZERS_MODE"); v != "" {
		cfg.Mode = v
	}
	if v := os.Getenv("STARGAZERS_USERNAME"); v != "" {
		cfg.Username = v
	}

	// 3. Overwrite with CLI flags if set
	if cliRepo != "" {
		cfg.Repo = cliRepo
	}
	if cliToken != "" {
		cfg.Token = cliToken
	}
	if cliCacheDir != "" {
		cfg.CacheDir = cliCacheDir
	}
	if cliMode != "" {
		cfg.Mode = cliMode
	}
	if cliUsername != "" {
		cfg.Username = cliUsername
	}

	// Set defaults if still empty
	if cfg.CacheDir == "" {
		cfg.CacheDir = "./stargazer_cache"
	}
	if cfg.Mode == "" {
		cfg.Mode = "basic"
	}

	// Validate required fields
	missing := []string{}
	if cfg.Repo == "" {
		missing = append(missing, "repo")
	}
	if cfg.Token == "" {
		missing = append(missing, "token (GH_TOKEN)")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required config: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}
