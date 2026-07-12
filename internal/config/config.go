package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const DefaultBaseURL = "https://aaasp.ai"

type Config struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".aaasp", "config.json"), nil
}

func Load() (*Config, error) {
	cfg := &Config{BaseURL: DefaultBaseURL}

	// Env var takes priority
	if key := os.Getenv("AAASP_API_KEY"); key != "" {
		cfg.APIKey = key
	}
	if url := os.Getenv("AAASP_BASE_URL"); url != "" {
		cfg.BaseURL = url
	}

	path, err := configPath()
	if err != nil {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}

	var file Config
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, err
	}

	// File values only fill in if env didn't set them
	if cfg.APIKey == "" {
		cfg.APIKey = file.APIKey
	}
	if cfg.BaseURL == DefaultBaseURL && file.BaseURL != "" {
		cfg.BaseURL = file.BaseURL
	}

	return cfg, nil
}

func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
