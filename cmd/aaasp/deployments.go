package main

import (
	"fmt"
	"io"
	"os"

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
	case "create":
		if len(args) == 0 {
			output.Fatal("usage: aaasp deployments create <agent_def_id>")
		}
		deploymentCreate(client, args[0])
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

func deploymentCreate(client *api.Client, agentDefID string) {
	body := map[string]string{"agent_def_id": agentDefID}

	var result map[string]any
	if err := client.Post("/deployments", body, &result); err != nil {
		output.Fatal("%v", err)
	}
	if output.IsJSON() {
		output.JSON(result)
		return
	}

	fmt.Printf("Created deployment %s\n", str(result["id"]))
	output.KV([][2]string{
		{"agent_def_id", str(result["agent_def_id"])},
		{"status", str(result["status"])},
		{"runnable", readyLabel(result)},
	})
	if reason := str(result["readiness_error"]); reason != "" {
		fmt.Fprintf(os.Stderr, "\nwarning: this deployment cannot run yet: %s\n", reason)
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
	renderDeploymentList(os.Stdout, deps)
}

func renderDeploymentList(w io.Writer, deps []any) {
	fmt.Fprintf(w, "  %-36s  %-36s  %-10s  %-5s  %s\n", "ID", "AGENT", "STATUS", "READY", "CREATED")
	fmt.Fprintf(w, "  %-36s  %-36s  %-10s  %-5s  %s\n", "---", "-----", "------", "-----", "-------")

	var blocked [][2]string
	for _, d := range deps {
		m, _ := d.(map[string]any)
		fmt.Fprintf(w, "  %-36s  %-36s  %-10s  %-5s  %s\n",
			str(m["id"]),
			str(m["agent_def_id"]),
			str(m["status"]),
			readyLabel(m),
			str(m["created_at"]),
		)
		if reason := str(m["readiness_error"]); reason != "" {
			blocked = append(blocked, [2]string{str(m["id"]), reason})
		}
	}

	if len(blocked) > 0 {
		fmt.Fprintln(w, "\nNot runnable:")
		for _, b := range blocked {
			fmt.Fprintf(w, "  %s: %s\n", b[0], b[1])
		}
	}
}

// readyLabel renders the API's runnable field; "?" for servers that predate it.
func readyLabel(m map[string]any) string {
	switch v := m["runnable"].(type) {
	case bool:
		if v {
			return "yes"
		}
		return "no"
	default:
		return "?"
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
	renderDeploymentDetail(os.Stdout, result)
}

func renderDeploymentDetail(w io.Writer, m map[string]any) {
	pairs := [][2]string{
		{"id", str(m["id"])},
		{"agent", str(m["agent_def_id"])},
		{"status", str(m["status"])},
		{"credential", str(m["credential_id"])},
		{"created", str(m["created_at"])},
		{"runnable", readyLabel(m)},
	}
	if reason := str(m["readiness_error"]); reason != "" {
		pairs = append(pairs, [2]string{"reason", reason})
	}
	output.KVTo(w, pairs)
}

func deploymentDelete(client *api.Client, id string) {
	if err := client.Delete("/deployments/" + id); err != nil {
		output.Fatal("%v", err)
	}
	fmt.Printf("Deleted deployment %s\n", id)
}
