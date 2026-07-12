package main

import (
	"fmt"
	"strings"

	"github.com/maco144/aaasp-cli/internal/api"
	"github.com/maco144/aaasp-cli/internal/output"
)

func cmdRuns(args []string) {
	sub := "list"
	if len(args) > 0 {
		sub = args[0]
		args = args[1:]
	}

	cfg := mustLoadConfig()
	client := api.New(cfg.BaseURL, cfg.APIKey)

	switch sub {
	case "list", "ls":
		runsList(client)
	case "create":
		runCreate(client, args)
	case "show", "get":
		if len(args) == 0 {
			output.Fatal("usage: aaasp runs show <id>")
		}
		runShow(client, args[0])
	case "cancel":
		if len(args) == 0 {
			output.Fatal("usage: aaasp runs cancel <id>")
		}
		runCancel(client, args[0])
	default:
		output.Fatal("unknown subcommand: %s", sub)
	}
}

func runCreate(client *api.Client, args []string) {
	sync := false
	positional := args[:0]
	for _, a := range args {
		if a == "--sync" {
			sync = true
		} else {
			positional = append(positional, a)
		}
	}

	if len(positional) == 0 {
		output.Fatal("usage: aaasp runs create <deployment_id> [prompt...] [--sync]")
	}

	deploymentID := positional[0]
	prompt := strings.Join(positional[1:], " ")

	body := map[string]any{"deployment_id": deploymentID, "prompt": prompt}
	if sync {
		body["sync"] = true
	}

	var result map[string]any
	if err := client.Post("/runs", body, &result); err != nil {
		output.Fatal("%v", err)
	}
	if output.IsJSON() {
		output.JSON(result)
		return
	}

	fmt.Printf("Created run %s\n", str(result["id"]))
	output.KV([][2]string{
		{"status", str(result["status"])},
		{"deployment", str(result["deployment_id"])},
	})
}

func runsList(client *api.Client) {
	var result map[string]any
	if err := client.Get("/runs", &result); err != nil {
		output.Fatal("%v", err)
	}
	if output.IsJSON() {
		output.JSON(result)
		return
	}

	runs, _ := result["runs"].([]any)
	if len(runs) == 0 {
		fmt.Println("No runs found.")
		return
	}
	fmt.Printf("  %-36s  %-12s  %s\n", "ID", "STATUS", "CREATED")
	fmt.Printf("  %-36s  %-12s  %s\n", "---", "------", "-------")
	for _, r := range runs {
		m, _ := r.(map[string]any)
		fmt.Printf("  %-36s  %-12s  %s\n",
			str(m["id"]),
			str(m["status"]),
			str(m["inserted_at"]),
		)
	}
}

func runShow(client *api.Client, id string) {
	var result map[string]any
	if err := client.Get("/runs/"+id, &result); err != nil {
		output.Fatal("%v", err)
	}
	if output.IsJSON() {
		output.JSON(result)
		return
	}
	output.KV([][2]string{
		{"id", str(result["id"])},
		{"status", str(result["status"])},
		{"deployment", str(result["deployment_id"])},
		{"created", str(result["inserted_at"])},
	})
}

func runCancel(client *api.Client, id string) {
	var result map[string]any
	if err := client.Post("/runs/"+id+"/cancel", nil, &result); err != nil {
		output.Fatal("%v", err)
	}
	fmt.Printf("Cancelled run %s\n", id)
}
