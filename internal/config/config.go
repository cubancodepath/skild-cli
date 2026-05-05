package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ErrSetupRequired = errors.New("setup required")

type Config struct {
	RepoURL           string `json:"repoUrl"`
	RepoRef           string `json:"repoRef"`
	CacheDir          string `json:"cacheDir"`
	RootPath          string `json:"rootPath"`
	OpenCodeDir       string `json:"openCodeDir"`
	GlobalOpenCodeDir string `json:"globalOpenCodeDir"`
	InstallMode       string `json:"installMode"`
}

func defaults(home string) Config {
	return Config{
		RepoRef:           "main",
		CacheDir:          filepath.Join(home, ".cache", "skild"),
		RootPath:          "skills",
		OpenCodeDir:       filepath.Join(".", ".opencode", "skills"),
		GlobalOpenCodeDir: defaultOpenCodeGlobalSkillsDir(home),
		InstallMode:       "copy",
	}
}

func defaultOpenCodeGlobalSkillsDir(home string) string {
	xdgConfigHome := strings.TrimSpace(os.Getenv("XDG_CONFIG_HOME"))
	if xdgConfigHome == "" {
		xdgConfigHome = filepath.Join(home, ".config")
	}

	return filepath.Join(xdgConfigHome, "opencode", "skills")
}

func Default() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("resolve home dir: %w", err)
	}
	return defaults(home), nil
}

func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".config", "skild", "config.json"), nil
}

func Load() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("resolve home dir: %w", err)
	}

	path, err := ConfigPath()
	if err != nil {
		return Config{}, err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, ErrSetupRequired
		}
		return Config{}, fmt.Errorf("read config file: %w", err)
	}

	if len(strings.TrimSpace(string(b))) == 0 {
		return Config{}, ErrSetupRequired
	}

	cfg := defaults(home)
	if err := json.Unmarshal(b, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config file: %w", err)
	}

	applyDefaults(&cfg, defaults(home))
	cfg.CacheDir = expandHome(cfg.CacheDir, home)
	cfg.GlobalOpenCodeDir = expandHome(cfg.GlobalOpenCodeDir, home)

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func Save(cfg Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("resolve home dir: %w", err)
	}

	applyDefaults(&cfg, defaults(home))
	cfg.CacheDir = expandHome(cfg.CacheDir, home)
	cfg.GlobalOpenCodeDir = expandHome(cfg.GlobalOpenCodeDir, home)

	if err := cfg.Validate(); err != nil {
		return err
	}

	path, err := ConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	b = append(b, '\n')

	if err := os.WriteFile(path, b, 0o644); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}

	return nil
}

func (c Config) Validate() error {
	if c.RepoURL == "" {
		return fmt.Errorf("missing required field repoUrl")
	}

	if c.InstallMode != "copy" && c.InstallMode != "symlink" {
		return fmt.Errorf("invalid installMode %q: expected copy or symlink", c.InstallMode)
	}

	return nil
}

func applyDefaults(cfg *Config, d Config) {
	if strings.TrimSpace(cfg.RepoRef) == "" {
		cfg.RepoRef = d.RepoRef
	}
	if strings.TrimSpace(cfg.CacheDir) == "" {
		cfg.CacheDir = d.CacheDir
	}
	if strings.TrimSpace(cfg.RootPath) == "" {
		cfg.RootPath = d.RootPath
	}
	if strings.TrimSpace(cfg.OpenCodeDir) == "" {
		cfg.OpenCodeDir = d.OpenCodeDir
	}
	if strings.TrimSpace(cfg.GlobalOpenCodeDir) == "" {
		cfg.GlobalOpenCodeDir = d.GlobalOpenCodeDir
	}
	if strings.TrimSpace(cfg.InstallMode) == "" {
		cfg.InstallMode = d.InstallMode
	}
}

func expandHome(pathValue, home string) string {
	if pathValue == "~" {
		return home
	}
	if strings.HasPrefix(pathValue, "~/") {
		return filepath.Join(home, strings.TrimPrefix(pathValue, "~/"))
	}
	return pathValue
}
