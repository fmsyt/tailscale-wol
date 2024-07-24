package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Hosts   []ConfigHost   `json:"hosts"`
	Targets []ConfigTarget `json:"targets"`
}

type ConfigHost struct {
	Host     string  `json:"host"`
	User     string  `json:"user"`
	Port     *int    `json:"port"`
	Password *string `json:"password"`
	Identity *string `json:"identity"`
}

type ConfigTarget struct {
	Mac              string  `json:"mac"`
	Port             *int    `json:"port"`
	Ip               *string `json:"ip"`
	PreferredCommand *string `json:"preferredCommand"`
}

func getConfig() (*Config, error) {
	execDir, err := appPath()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(execDir, "config.json")

	bytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(bytes, &config)

	return &config, err
}

func findTarget(mac string, config *Config) *ConfigTarget {
	for _, target := range config.Targets {
		if target.Mac == mac {
			return &target
		}
	}

	return nil
}
