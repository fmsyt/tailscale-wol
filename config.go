package main

import (
	"encoding/json"
	"os"
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
	Mac  string  `json:"mac"`
	Port *int    `json:"port"`
	Ip   *string `json:"ip"`
}

func getConfig() (*Config, error) {
	bytes, err := os.ReadFile("config.json")
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
