package main

import (
	"fmt"

	"github.com/maco144/aaasp-cli/internal/api"
	"github.com/maco144/aaasp-cli/internal/output"
)

func cmdSkills(args []string) {
	sub := "list"
	if len(args) > 0 {
		sub = args[0]
		args = args[1:]
	}

	cfg := mustLoadConfig()
	client := api.New(cfg.BaseURL, cfg.APIKey)

	switch sub {
	case "list", "ls":
		skillsList(client)
	case "show", "get":
		if len(args) == 0 {
			output.Fatal("usage: aaasp skills show <id>")
		}
		skillShow(client, args[0])
	default:
		output.Fatal("unknown subcommand: %s", sub)
	}
}

func skillsList(client *api.Client) {
	var result map[string]any
	if err := client.Get("/skills", &result); err != nil {
		output.Fatal("%v", err)
	}
	if output.IsJSON() {
		output.JSON(result)
		return
	}

	skills, _ := result["skills"].([]any)
	if len(skills) == 0 {
		fmt.Println("No skills found.")
		return
	}
	fmt.Printf("  %-36s  %-24s  %s\n", "ID", "NAME", "DESCRIPTION")
	fmt.Printf("  %-36s  %-24s  %s\n", "---", "----", "-----------")
	for _, s := range skills {
		m, _ := s.(map[string]any)
		desc := str(m["description"])
		if len(desc) > 50 {
			desc = desc[:47] + "..."
		}
		fmt.Printf("  %-36s  %-24s  %s\n", str(m["id"]), str(m["name"]), desc)
	}
}

func skillShow(client *api.Client, id string) {
	var result map[string]any
	if err := client.Get("/skills/"+id, &result); err != nil {
		output.Fatal("%v", err)
	}
	if output.IsJSON() {
		output.JSON(result)
		return
	}
	output.KV([][2]string{
		{"id", str(result["id"])},
		{"name", str(result["name"])},
		{"description", str(result["description"])},
		{"source", str(result["source_url"])},
	})
}
