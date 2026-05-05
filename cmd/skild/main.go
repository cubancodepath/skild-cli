package main

import (
	"fmt"
	"os"
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
	fmt.Println("  version   Show skild version")
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "help", "--help", "-h":
		printHelp()
	case "version", "--version", "-v":
		fmt.Println(version)
	default:
		fmt.Printf("Unknown command: %s\n\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}
