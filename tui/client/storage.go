package client

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFile = "wordle-go.json"

type StoredSettings struct {
	Theme       string `json:"theme,omitempty"`
	Language    string `json:"language,omitempty"`
	DisplayMode string `json:"display_mode,omitempty"`
}

type StoredConfig struct {
	Token    string         `json:"token,omitempty"`
	Settings StoredSettings `json:"settings,omitempty"`
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configFile), nil
}

func LoadConfig() (*StoredConfig, error) {
	path, err := configPath()
	if err != nil {
		return &StoredConfig{}, nil
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &StoredConfig{}, nil
	}
	if err != nil {
		return nil, err
	}
	var cfg StoredConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return &StoredConfig{}, nil
	}
	return &cfg, nil
}

func SaveConfig(cfg *StoredConfig) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
