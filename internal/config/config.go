package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Language string `json:"language"`
}

func configDir() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "jotoba-tui")
}

func configPath() string {
	d := configDir()
	if d == "" {
		return ""
	}
	return filepath.Join(d, "config.json")
}

func Load() *Config {
	cfg := &Config{Language: "English"}

	path := configPath()
	if path == "" {
		return cfg
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}

	var loaded Config
	if err := json.Unmarshal(data, &loaded); err != nil {
		return cfg
	}

	if loaded.Language != "" {
		cfg.Language = loaded.Language
	}

	return cfg
}

func Save(cfg *Config) error {
	path := configPath()
	if path == "" {
		return nil
	}

	dir := configDir()
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
