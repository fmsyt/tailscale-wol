package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Hosts []ConfigHost `json:"hosts"`
}

type ConfigHost struct {
	Host     string  `json:"host"`
	User     string  `json:"user"`
	Port     *int    `json:"port"`
	Password *string `json:"password"`
	Identity *string `json:"identity"`
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
