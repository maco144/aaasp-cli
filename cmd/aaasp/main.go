package main

import (
	"fmt"
	"os"

	"github.com/maco144/aaasp-cli/internal/config"
	"github.com/maco144/aaasp-cli/internal/output"
)

const version = "0.1.0"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		printHelp()
		os.Exit(0)
	}

	// Global flags
	filtered := args[:0]
	for _, a := range args {
		if a == "--json" {
			output.SetJSON(true)
		} else {
			filtered = append(filtered, a)
		}
	}
	args = filtered

	cmd := args[0]
	rest := args[1:]

	switch cmd {
	case "register":
		cmdRegister(rest)
	case "whoami":
		cmdWhoami(rest)
	case "deployments", "deployment":
		cmdDeployments(rest)
	case "runs", "run":
		cmdRuns(rest)
	case "credentials", "credential", "creds":
		cmdCredentials(rest)
	case "skills", "skill":
		cmdSkills(rest)
	case "version", "--version", "-v":
		fmt.Println(version)
	case "help", "--help", "-h":
		printHelp()
	default:
		output.Error("unknown command: %s", cmd)
		printHelp()
		os.Exit(1)
	}
}

func mustLoadConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		output.Fatal("failed to load config: %v", err)
	}
	if cfg.APIKey == "" {
		output.Fatal("no API key found\n\n  Set AAASP_API_KEY or run: aaasp register")
	}
	return cfg
}

func printHelp() {
	fmt.Print(`aaasp — Agent-as-a-Service Platform CLI

Usage:
  aaasp <command> [subcommand] [flags]

Commands:
  register              Create a new account
  whoami                Show current account info

  deployments list      List deployments
  deployments create    Create a deployment
  deployments delete    Delete a deployment

  runs list             List runs
  runs show <id>        Show run details
  runs create           Trigger a run
  runs cancel <id>      Cancel a run

  credentials list      List stored credentials
  credentials add       Add a credential
  credentials delete    Delete a credential

  skills list           List available skills

Flags:
  --json                Output raw JSON

Environment:
  AAASP_API_KEY         API key (required for most commands)
  AAASP_BASE_URL        Override API base URL (default: https://aaasp.ai)
`)
}
