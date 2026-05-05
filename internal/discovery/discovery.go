package discovery

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type Skill struct {
	Name string
	Path string
}

func List(repoDir, rootPath string) ([]Skill, error) {
	base := filepath.Join(repoDir, rootPath)

	entries, err := os.ReadDir(base)
	if err != nil {
		return nil, fmt.Errorf("read skills root %q: %w", base, err)
	}

	skills := make([]Skill, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillName := entry.Name()
		skillDir := filepath.Join(base, skillName)
		skillFile := filepath.Join(skillDir, "SKILL.md")

		info, err := os.Stat(skillFile)
		if err != nil || info.IsDir() {
			continue
		}

		skills = append(skills, Skill{Name: skillName, Path: skillDir})
	}

	sort.Slice(skills, func(i, j int) bool {
		return skills[i].Name < skills[j].Name
	})

	return skills, nil
}
