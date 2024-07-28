package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ConfigSchema struct {
	Hosts   []ConnectionHostSchema `json:"hosts"`
	Targets *[]WoLTargetSchema     `json:"targets"`
}

type Config struct {
	Hosts   []ConnectionHost
	Targets []WoLTarget
}

func (c ConfigSchema) toConfig() Config {
	hosts := make([]ConnectionHost, 0, len(c.Hosts))
	for _, host := range c.Hosts {
		hosts = append(hosts, host.toConnectionHost())
	}

	var targets []WoLTarget
	if c.Targets != nil {
		targets = make([]WoLTarget, 0, len(*c.Targets))
		for _, target := range *c.Targets {
			targets = append(targets, target.toWoLTarget())
		}
	}

	return Config{
		Hosts:   hosts,
		Targets: targets,
	}
}

type ConnectionHostSchema struct {
	Host         string  `json:"host"`
	User         string  `json:"user"`
	Port         *int    `json:"port"`
	Password     *string `json:"password"`
	IdentityFile *string `json:"identityFile"`
	Timeout      *int    `json:"timeout"`
}

type ConnectionHost struct {
	Host     string
	User     string
	Timeout  int
	Port     int
	Password *string
	Identity *string
}

func (c ConnectionHostSchema) toConnectionHost() ConnectionHost {
	port := 22
	if c.Port != nil {
		port = *c.Port
	}

	timeout := 300
	if c.Timeout != nil {
		timeout = *c.Timeout
	}

	return ConnectionHost{
		Host:     c.Host,
		User:     c.User,
		Timeout:  timeout,
		Port:     port,
		Password: c.Password,
		Identity: c.IdentityFile,
	}
}

type WoLTargetSchema struct {
	Mac              string  `json:"mac"`
	Port             *int    `json:"port"`
	Ip               *string `json:"ip"`
	PreferredCommand *string `json:"preferredCommand"`
}

type WoLTarget struct {
	Mac              string
	Port             int
	Ip               string
	PreferredCommand string
}

func (w WoLTargetSchema) toWoLTarget() WoLTarget {
	port := 9
	if w.Port != nil {
		port = *w.Port
	}

	ip := "255.255.255.255"
	if w.Ip != nil {
		ip = *w.Ip
	}

	command := "wol"
	if w.PreferredCommand != nil {
		command = *w.PreferredCommand
	}

	return WoLTarget{
		Mac:              w.Mac,
		Port:             port,
		Ip:               ip,
		PreferredCommand: command,
	}
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

	var cs ConfigSchema
	err = json.Unmarshal(bytes, &cs)

	c := cs.toConfig()

	return &c, err
}

func findHost(mac string, config Config) *WoLTarget {
	for _, target := range config.Targets {
		if target.Mac == mac {
			return &target
		}
	}

	return nil
}
