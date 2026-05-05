package main

import (
	"fmt"
	"os"

	"github.com/cubancodepath/skild/internal/config"
	"github.com/cubancodepath/skild/internal/discovery"
	"github.com/cubancodepath/skild/internal/install"
	"github.com/cubancodepath/skild/internal/repo"
)

const version = "dev"

func printHelp() {
	fmt.Println("skild - Manage OpenCode skills")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  skild <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  help      Show this help message")
	fmt.Println("  config    Show resolved configuration")
	fmt.Println("  install   Install one skill or --all")
	fmt.Println("  list      List available skills")
	fmt.Println("  repo-sync Clone/update cached repository")
	fmt.Println("  version   Show skild version")
}

func loadSkills() ([]discovery.Skill, config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, config.Config{}, err
	}

	repoDir, err := repo.Prepare(cfg)
	if err != nil {
		return nil, config.Config{}, err
	}

	skills, err := discovery.List(repoDir, cfg.RootPath)
	if err != nil {
		return nil, config.Config{}, err
	}

	return skills, cfg, nil
}

func runList() error {
	skills, _, err := loadSkills()
	if err != nil {
		return err
	}

	if len(skills) == 0 {
		fmt.Println("No skills found.")
		return nil
	}

	fmt.Println("Available skills:")
	for _, skill := range skills {
		fmt.Printf("- %s\n", skill.Name)
	}

	return nil
}

func runInstall(args []string) error {
	skills, cfg, err := loadSkills()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return fmt.Errorf("install requires <skill-name> or --all")
	}

	if args[0] == "--all" {
		if len(skills) == 0 {
			fmt.Println("No skills found.")
			return nil
		}

		installed, err := install.InstallAll(skills, cfg.OpenCodeDir)
		if err != nil {
			return err
		}

		fmt.Printf("Installed %d skills into %s\n", len(installed), cfg.OpenCodeDir)
		return nil
	}

	selected, ok := install.SkillByName(skills, args[0])
	if !ok {
		return fmt.Errorf("skill %q not found", args[0])
	}

	dst, err := install.InstallSkill(selected, cfg.OpenCodeDir)
	if err != nil {
		return err
	}

	fmt.Printf("Installed %s to %s\n", selected.Name, dst)
	return nil
}

func printConfig() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Println("Resolved configuration:")
	fmt.Printf("  SKILD_REPO_URL=%s\n", cfg.RepoURL)
	fmt.Printf("  SKILD_REPO_REF=%s\n", cfg.RepoRef)
	fmt.Printf("  SKILD_CACHE_DIR=%s\n", cfg.CacheDir)
	fmt.Printf("  SKILD_ROOT_PATH=%s\n", cfg.RootPath)
	fmt.Printf("  SKILD_OPENCODE_DIR=%s\n", cfg.OpenCodeDir)
	fmt.Printf("  SKILD_INSTALL_MODE=%s\n", cfg.InstallMode)

	return nil
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "help", "--help", "-h":
		printHelp()
	case "config":
		if err := printConfig(); err != nil {
			fmt.Printf("Config error: %v\n", err)
			os.Exit(1)
		}
	case "repo-sync":
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Config error: %v\n", err)
			os.Exit(1)
		}

		repoDir, err := repo.Prepare(cfg)
		if err != nil {
			fmt.Printf("Repo error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Repository ready at %s\n", repoDir)
	case "list":
		if err := runList(); err != nil {
			fmt.Printf("List error: %v\n", err)
			os.Exit(1)
		}
	case "install":
		if err := runInstall(os.Args[2:]); err != nil {
			fmt.Printf("Install error: %v\n", err)
			os.Exit(1)
		}
	case "version", "--version", "-v":
		fmt.Println(version)
	default:
		fmt.Printf("Unknown command: %s\n\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}
