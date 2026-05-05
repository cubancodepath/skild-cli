package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	RepoURL      string
	RepoRef      string
	CacheDir     string
	RootPath     string
	OpenCodeDir  string
	InstallMode  string
}

func Load() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("resolve home dir: %w", err)
	}

	cfg := Config{
		RepoURL:     strings.TrimSpace(os.Getenv("SKILD_REPO_URL")),
		RepoRef:     getEnvOrDefault("SKILD_REPO_REF", "main"),
		CacheDir:    getEnvOrDefault("SKILD_CACHE_DIR", filepath.Join(home, ".cache", "skild")),
		RootPath:    getEnvOrDefault("SKILD_ROOT_PATH", "skills"),
		OpenCodeDir: getEnvOrDefault("SKILD_OPENCODE_DIR", filepath.Join(".", ".opencode", "skills")),
		InstallMode: getEnvOrDefault("SKILD_INSTALL_MODE", "copy"),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	if c.RepoURL == "" {
		return fmt.Errorf("missing required env var SKILD_REPO_URL")
	}

	if c.InstallMode != "copy" && c.InstallMode != "symlink" {
		return fmt.Errorf("invalid SKILD_INSTALL_MODE %q: expected copy or symlink", c.InstallMode)
	}

	return nil
}

func getEnvOrDefault(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	return v
}
