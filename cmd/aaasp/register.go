package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/maco144/aaasp-cli/internal/api"
	"github.com/maco144/aaasp-cli/internal/config"
	"github.com/maco144/aaasp-cli/internal/output"
	"golang.org/x/term"
)

func cmdRegister(args []string) {
	cfg, _ := config.Load()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Password: ")
	pwBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		output.Fatal("failed to read password: %v", err)
	}
	password := string(pwBytes)

	client := api.New(cfg.BaseURL, "")
	result, err := client.Register(email, password)
	if err != nil {
		output.Fatal("%v", err)
	}

	apiKey, _ := result["api_key"].(string)
	if apiKey == "" {
		output.Fatal("registration succeeded but no API key returned")
	}

	// Save to config file
	cfg.APIKey = apiKey
	if err := config.Save(cfg); err != nil {
		output.Error("could not save config: %v", err)
	}

	fmt.Printf("\nRegistered successfully!\n\n")
	fmt.Printf("  API key: %s\n\n", apiKey)
	fmt.Printf("Your key has been saved to ~/.aaasp/config.json\n")
	fmt.Printf("You can also set it via: export AAASP_API_KEY=%s\n", apiKey)
}
