package install

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cubancodepath/skild/internal/discovery"
)

func SkillByName(skills []discovery.Skill, name string) (discovery.Skill, bool) {
	for _, skill := range skills {
		if strings.EqualFold(skill.Name, name) {
			return skill, true
		}
	}

	return discovery.Skill{}, false
}

func InstallSkill(skill discovery.Skill, openCodeDir string) (string, error) {
	if err := os.MkdirAll(openCodeDir, 0o755); err != nil {
		return "", fmt.Errorf("create destination dir: %w", err)
	}

	dstName := sanitizeName(skill.Name)
	dst := filepath.Join(openCodeDir, dstName)

	if err := os.RemoveAll(dst); err != nil {
		return "", fmt.Errorf("remove previous installation: %w", err)
	}

	if err := copyDir(skill.Path, dst); err != nil {
		return "", fmt.Errorf("copy skill files: %w", err)
	}

	return dst, nil
}

func InstallAll(skills []discovery.Skill, openCodeDir string) ([]string, error) {
	installed := make([]string, 0, len(skills))
	for _, skill := range skills {
		dst, err := InstallSkill(skill, openCodeDir)
		if err != nil {
			return installed, fmt.Errorf("install %q: %w", skill.Name, err)
		}
		installed = append(installed, dst)
	}

	return installed, nil
}

func sanitizeName(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.' {
			return r
		}
		return '-'
	}, s)
	s = strings.Trim(s, ".-")
	if s == "" {
		return "unnamed-skill"
	}
	return s
}

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
			continue
		}

		if err := copyFile(srcPath, dstPath); err != nil {
			return err
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return nil
}
