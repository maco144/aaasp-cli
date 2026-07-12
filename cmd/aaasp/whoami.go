package main

import (
	"fmt"

	"github.com/maco144/aaasp-cli/internal/api"
	"github.com/maco144/aaasp-cli/internal/output"
)

func cmdWhoami(args []string) {
	cfg := mustLoadConfig()
	client := api.New(cfg.BaseURL, cfg.APIKey)

	var account map[string]any
	if err := client.Get("/account", &account); err != nil {
		output.Fatal("%v", err)
	}

	if output.IsJSON() {
		output.JSON(account)
		return
	}

	output.KV([][2]string{
		{"email", str(account["email"])},
		{"plan", str(account["plan"])},
		{"credits", fmt.Sprintf("%v", account["credit_balance"])},
		{"active", fmt.Sprintf("%v", account["active"])},
	})
}

func str(v any) string {
	if v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}
