package main

import (
	"fmt"

	"github.com/maco144/aaasp-cli/internal/api"
	"github.com/maco144/aaasp-cli/internal/output"
)

func cmdDeployments(args []string) {
	sub := "list"
	if len(args) > 0 {
		sub = args[0]
		args = args[1:]
	}

	cfg := mustLoadConfig()
	client := api.New(cfg.BaseURL, cfg.APIKey)

	switch sub {
	case "list", "ls":
		deploymentslist(client)
	case "show", "get":
		if len(args) == 0 {
			output.Fatal("usage: aaasp deployments show <id>")
		}
		deploymentShow(client, args[0])
	case "delete", "rm":
		if len(args) == 0 {
			output.Fatal("usage: aaasp deployments delete <id>")
		}
		deploymentDelete(client, args[0])
	default:
		output.Fatal("unknown subcommand: %s", sub)
	}
}

func deploymentslist(client *api.Client) {
	var result map[string]any
	if err := client.Get("/deployments", &result); err != nil {
		output.Fatal("%v", err)
	}
	if output.IsJSON() {
		output.JSON(result)
		return
	}

	deps, _ := result["deployments"].([]any)
	if len(deps) == 0 {
		fmt.Println("No deployments found.")
		return
	}
	fmt.Printf("  %-36s  %-36s  %-10s  %s\n", "ID", "AGENT", "STATUS", "CREATED")
	fmt.Printf("  %-36s  %-36s  %-10s  %s\n", "---", "-----", "------", "-------")
	for _, d := range deps {
		m, _ := d.(map[string]any)
		fmt.Printf("  %-36s  %-36s  %-10s  %s\n",
			str(m["id"]),
			str(m["agent_def_id"]),
			str(m["status"]),
			str(m["created_at"]),
		)
	}
}

func deploymentShow(client *api.Client, id string) {
	var result map[string]any
	if err := client.Get("/deployments/"+id, &result); err != nil {
		output.Fatal("%v", err)
	}
	if output.IsJSON() {
		output.JSON(result)
		return
	}
	output.KV([][2]string{
		{"id", str(result["id"])},
		{"agent", str(result["agent_def_id"])},
		{"status", str(result["status"])},
		{"schedule", str(result["schedule"])},
	})
}

func deploymentDelete(client *api.Client, id string) {
	if err := client.Delete("/deployments/" + id); err != nil {
		output.Fatal("%v", err)
	}
	fmt.Printf("Deleted deployment %s\n", id)
}
