package repo

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cubancodepath/skild/internal/config"
)

func Prepare(cfg config.Config) (string, error) {
	repoDir := RepoDir(cfg)

	if err := os.MkdirAll(cfg.CacheDir, 0o755); err != nil {
		return "", fmt.Errorf("create cache dir: %w", err)
	}

	gitDir := filepath.Join(repoDir, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		if err := runGit(repoDir, "fetch", "origin"); err != nil {
			return "", fmt.Errorf("fetch repository: %w", err)
		}
		if err := runGit(repoDir, "checkout", cfg.RepoRef); err != nil {
			return "", fmt.Errorf("checkout ref %q: %w", cfg.RepoRef, err)
		}
		if err := runGit(repoDir, "pull", "--ff-only", "origin", cfg.RepoRef); err != nil {
			return "", fmt.Errorf("pull ref %q: %w", cfg.RepoRef, err)
		}
		return repoDir, nil
	}

	if err := runGit("", "clone", "--branch", cfg.RepoRef, "--single-branch", cfg.RepoURL, repoDir); err != nil {
		return "", fmt.Errorf("clone repository: %w", err)
	}

	return repoDir, nil
}

func RepoDir(cfg config.Config) string {
	h := sha1.Sum([]byte(cfg.RepoURL + "@" + cfg.RepoRef))
	name := hex.EncodeToString(h[:8])
	return filepath.Join(cfg.CacheDir, name)
}

func runGit(workdir string, args ...string) error {
	cmd := exec.Command("git", args...)
	if workdir != "" {
		cmd.Dir = workdir
	}

	output, err := cmd.CombinedOutput()
	if isVerbose() && len(output) > 0 {
		fmt.Print(string(output))
	}

	if err != nil {
		msg := strings.TrimSpace(string(output))
		if msg == "" {
			return fmt.Errorf("git %v: %w", args, err)
		}
		return fmt.Errorf("git %v failed: %s", args, msg)
	}

	return nil
}

func isVerbose() bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("SKILD_VERBOSE")))
	return v == "1" || v == "true" || v == "yes"
}
