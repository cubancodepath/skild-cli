package main

import (
	"fmt"
	"os"

	"github.com/cubancodepath/skild/internal/config"
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
	fmt.Println("  version   Show skild version")
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
	case "version", "--version", "-v":
		fmt.Println(version)
	default:
		fmt.Printf("Unknown command: %s\n\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}
