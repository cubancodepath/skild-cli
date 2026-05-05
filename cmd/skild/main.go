package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

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
	fmt.Println("  install   Install one skill or --all [--global]")
	fmt.Println("  list      List available skills")
	fmt.Println("  repo-sync Clone/update cached repository")
	fmt.Println("  setup     Configure global skild settings")
	fmt.Println("  update    Sync repo and reinstall all skills [--global]")
	fmt.Println("  version   Show skild version")
}

func parseGlobalFlag(args []string) (bool, []string) {
	filtered := make([]string, 0, len(args))
	global := false

	for _, arg := range args {
		if arg == "--global" || arg == "-g" {
			global = true
			continue
		}
		filtered = append(filtered, arg)
	}

	return global, filtered
}

func targetDir(cfg config.Config, global bool) string {
	if global {
		return cfg.GlobalOpenCodeDir
	}
	return cfg.OpenCodeDir
}

func handleSetupRequired(err error) error {
	if errors.Is(err, config.ErrSetupRequired) {
		return fmt.Errorf("no global config found. run: skild setup")
	}
	return err
}

func loadSkills() ([]discovery.Skill, config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, config.Config{}, handleSetupRequired(err)
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

	global, filteredArgs := parseGlobalFlag(args)
	installDir := targetDir(cfg, global)

	if len(filteredArgs) == 0 {
		return fmt.Errorf("install requires <skill-name> or --all")
	}

	if filteredArgs[0] == "--all" {
		if len(skills) == 0 {
			fmt.Println("No skills found.")
			return nil
		}

		installed, err := install.InstallAll(skills, installDir)
		if err != nil {
			return err
		}

		scope := "local"
		if global {
			scope = "global"
		}
		fmt.Printf("Installed %d skills into %s (%s)\n", len(installed), installDir, scope)
		return nil
	}

	selected, ok := install.SkillByName(skills, filteredArgs[0])
	if !ok {
		return fmt.Errorf("skill %q not found", filteredArgs[0])
	}

	dst, err := install.InstallSkill(selected, installDir)
	if err != nil {
		return err
	}

	scope := "local"
	if global {
		scope = "global"
	}
	fmt.Printf("Installed %s to %s (%s)\n", selected.Name, dst, scope)
	return nil
}

func runUpdate(args []string) error {
	skills, cfg, err := loadSkills()
	if err != nil {
		return err
	}

	global, _ := parseGlobalFlag(args)
	installDir := targetDir(cfg, global)

	if len(skills) == 0 {
		fmt.Println("No skills found.")
		return nil
	}

	installed, err := install.InstallAll(skills, installDir)
	if err != nil {
		return err
	}

	scope := "local"
	if global {
		scope = "global"
	}
	fmt.Printf("Updated %d skills into %s (%s)\n", len(installed), installDir, scope)
	return nil
}

func printConfig() error {
	cfg, err := config.Load()
	if err != nil {
		return handleSetupRequired(err)
	}

	path, err := config.ConfigPath()
	if err != nil {
		return err
	}

	fmt.Println("Resolved configuration:")
	fmt.Printf("  CONFIG_PATH=%s\n", path)
	fmt.Printf("  repoUrl=%s\n", cfg.RepoURL)
	fmt.Printf("  repoRef=%s\n", cfg.RepoRef)
	fmt.Printf("  cacheDir=%s\n", cfg.CacheDir)
	fmt.Printf("  rootPath=%s\n", cfg.RootPath)
	fmt.Printf("  openCodeDir=%s\n", cfg.OpenCodeDir)
	fmt.Printf("  globalOpenCodeDir=%s\n", cfg.GlobalOpenCodeDir)
	fmt.Printf("  installMode=%s\n", cfg.InstallMode)

	return nil
}

func promptWithDefault(reader *bufio.Reader, label, defaultValue string) (string, error) {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", label, defaultValue)
	} else {
		fmt.Printf("%s: ", label)
	}

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	value := strings.TrimSpace(input)
	if value == "" {
		return defaultValue, nil
	}
	return value, nil
}

func promptRepoURL(reader *bufio.Reader, defaultValue string) (string, error) {
	for {
		value, err := promptWithDefault(reader, "Repository URL", defaultValue)
		if err != nil {
			return "", err
		}
		if strings.TrimSpace(value) == "" {
			fmt.Println("Repository URL is required.")
			continue
		}
		return value, nil
	}
}

func promptInstallMode(reader *bufio.Reader, defaultValue string) (string, error) {
	for {
		value, err := promptWithDefault(reader, "Install mode (copy|symlink)", defaultValue)
		if err != nil {
			return "", err
		}
		if value == "copy" || value == "symlink" {
			return value, nil
		}
		fmt.Println("Invalid install mode. Use copy or symlink.")
	}
}

func runSetup() error {
	reader := bufio.NewReader(os.Stdin)

	path, err := config.ConfigPath()
	if err != nil {
		return err
	}

	existing, err := config.Default()
	if err != nil {
		return err
	}

	loaded, err := config.Load()
	if err != nil && !errors.Is(err, config.ErrSetupRequired) {
		return err
	}
	if err == nil {
		existing = loaded
	}

	fmt.Println("skild setup")
	fmt.Printf("Config file: %s\n\n", path)

	repoURL, err := promptRepoURL(reader, existing.RepoURL)
	if err != nil {
		return err
	}
	repoRef, err := promptWithDefault(reader, "Repository ref", existing.RepoRef)
	if err != nil {
		return err
	}
	cacheDir, err := promptWithDefault(reader, "Cache directory", existing.CacheDir)
	if err != nil {
		return err
	}
	rootPath, err := promptWithDefault(reader, "Root path in repository", existing.RootPath)
	if err != nil {
		return err
	}
	openCodeDir, err := promptWithDefault(reader, "OpenCode skills directory", existing.OpenCodeDir)
	if err != nil {
		return err
	}
	globalOpenCodeDir, err := promptWithDefault(reader, "Global OpenCode skills directory", existing.GlobalOpenCodeDir)
	if err != nil {
		return err
	}
	installMode, err := promptInstallMode(reader, existing.InstallMode)
	if err != nil {
		return err
	}

	cfg := config.Config{
		RepoURL:           repoURL,
		RepoRef:           repoRef,
		CacheDir:          cacheDir,
		RootPath:          rootPath,
		OpenCodeDir:       openCodeDir,
		GlobalOpenCodeDir: globalOpenCodeDir,
		InstallMode:       installMode,
	}

	fmt.Println("\nSummary:")
	fmt.Printf("  repoUrl=%s\n", cfg.RepoURL)
	fmt.Printf("  repoRef=%s\n", cfg.RepoRef)
	fmt.Printf("  cacheDir=%s\n", cfg.CacheDir)
	fmt.Printf("  rootPath=%s\n", cfg.RootPath)
	fmt.Printf("  openCodeDir=%s\n", cfg.OpenCodeDir)
	fmt.Printf("  globalOpenCodeDir=%s\n", cfg.GlobalOpenCodeDir)
	fmt.Printf("  installMode=%s\n", cfg.InstallMode)

	confirm, err := promptWithDefault(reader, "Save config? (y/N)", "N")
	if err != nil {
		return err
	}

	if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
		fmt.Println("Setup cancelled.")
		return nil
	}

	if err := config.Save(cfg); err != nil {
		return err
	}

	fmt.Printf("Saved config to %s\n", path)
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
	case "setup":
		if err := runSetup(); err != nil {
			fmt.Printf("Setup error: %v\n", err)
			os.Exit(1)
		}
	case "repo-sync":
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Config error: %v\n", handleSetupRequired(err))
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
	case "update":
		if err := runUpdate(os.Args[2:]); err != nil {
			fmt.Printf("Update error: %v\n", err)
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
