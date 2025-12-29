package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	LastListenPort   string   `json:"lastListenPort"`
	LastTargetAddr   string   `json:"lastTargetAddr"`
	AutoOpenFirewall bool     `json:"autoOpenFirewall"`
	AutoStart        bool     `json:"autoStart"`
	History          []string `json:"history"`       // Store as "Port|Target"
	TargetHistory    []string `json:"targetHistory"` // Store as "Target"
}

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".tunf")
	_ = os.MkdirAll(configDir, 0755)
	return filepath.Join(configDir, "config.json")
}

func LoadConfig() Config {
	configPath := GetConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{
			LastListenPort:   "5678",
			LastTargetAddr:   "localhost:1234",
			AutoOpenFirewall: false,
			AutoStart:        false,
			History:          []string{},
			TargetHistory:    []string{},
		}
	}

	var config Config
	_ = json.Unmarshal(data, &config)
	return config
}

func SaveConfig(config Config) error {
	configPath := GetConfigPath()
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func AddToHistory(history []string, port, target string) []string {
	entry := port + "|" + target
	return AddValueToHistory(history, entry)
}

func AddValueToHistory(history []string, value string) []string {
	// Remove if already exists
	for i, p := range history {
		if p == value {
			history = append(history[:i], history[i+1:]...)
			break
		}
	}
	// Add to front
	history = append([]string{value}, history...)
	// Limit to 10
	if len(history) > 10 {
		history = history[:10]
	}
	return history
}
