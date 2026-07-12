package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/maco144/aaasp-cli/internal/api"
	"github.com/maco144/aaasp-cli/internal/output"
	"golang.org/x/term"
)

func cmdCredentials(args []string) {
	sub := "list"
	if len(args) > 0 {
		sub = args[0]
		args = args[1:]
	}

	cfg := mustLoadConfig()
	client := api.New(cfg.BaseURL, cfg.APIKey)

	switch sub {
	case "list", "ls":
		credentialsList(client)
	case "add", "create":
		credentialsAdd(client)
	case "delete", "rm":
		if len(args) == 0 {
			output.Fatal("usage: aaasp credentials delete <id>")
		}
		credentialsDelete(client, args[0])
	default:
		output.Fatal("unknown subcommand: %s", sub)
	}
}

func credentialsList(client *api.Client) {
	var result map[string]any
	if err := client.Get("/credentials", &result); err != nil {
		output.Fatal("%v", err)
	}
	if output.IsJSON() {
		output.JSON(result)
		return
	}

	creds, _ := result["credentials"].([]any)
	if len(creds) == 0 {
		fmt.Println("No credentials stored.")
		fmt.Println("\n  Add one: aaasp credentials add")
		return
	}
	fmt.Printf("  %-36s  %-20s  %s\n", "ID", "PROVIDER", "LABEL")
	fmt.Printf("  %-36s  %-20s  %s\n", "---", "--------", "-----")
	for _, c := range creds {
		m, _ := c.(map[string]any)
		fmt.Printf("  %-36s  %-20s  %s\n",
			str(m["id"]),
			str(m["provider"]),
			str(m["label"]),
		)
	}
}

func credentialsAdd(client *api.Client) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Provider (e.g. openai, anthropic): ")
	provider, _ := reader.ReadString('\n')
	provider = strings.TrimSpace(provider)

	fmt.Print("Label (optional, e.g. my-key): ")
	label, _ := reader.ReadString('\n')
	label = strings.TrimSpace(label)

	fmt.Print("API key: ")
	keyBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		output.Fatal("failed to read key: %v", err)
	}

	body := map[string]string{
		"provider":  provider,
		"vault_key": string(keyBytes),
	}
	if label != "" {
		body["label"] = label
	}

	var result map[string]any
	if err := client.Post("/credentials", body, &result); err != nil {
		output.Fatal("%v", err)
	}

	fmt.Printf("Credential added (id: %s)\n", str(result["id"]))
}

func credentialsDelete(client *api.Client, id string) {
	if err := client.Delete("/credentials/" + id); err != nil {
		output.Fatal("%v", err)
	}
	fmt.Printf("Deleted credential %s\n", id)
}
