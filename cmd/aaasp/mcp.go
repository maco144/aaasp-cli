package main

import (
	"os"

	"github.com/maco144/aaasp-cli/internal/api"
	"github.com/maco144/aaasp-cli/internal/mcp"
	"github.com/maco144/aaasp-cli/internal/output"
)

func cmdMcp(args []string) {
	sub := "serve"
	if len(args) > 0 {
		sub = args[0]
	}

	switch sub {
	case "serve":
		mcpServe()
	default:
		output.Fatal("unknown subcommand: %s\n\nusage: aaasp mcp serve", sub)
	}
}

// mcpServe runs an MCP (Model Context Protocol) server over stdio, exposing
// the AAASP core loop — dispatch a goal, check a run, list deployments/agents
// — as MCP tools for any MCP client (Claude Desktop, Claude Code, etc.).
func mcpServe() {
	cfg := mustLoadConfig()
	client := api.New(cfg.BaseURL, cfg.APIKey)

	server := mcp.NewServer(client, version, os.Stderr)
	if err := server.Serve(os.Stdin, os.Stdout); err != nil {
		output.Fatal("mcp server error: %v", err)
	}
}
