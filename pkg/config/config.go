package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
}

type APIConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type CaddyConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	// whether to tunnel the caddy server fom the server.
	TunnelServer bool `json:"tunnelServer"`
}

type Config struct {
	Server ServerConfig `json:"server"`
	API    APIConfig    `json:"api"`
	Caddy  CaddyConfig  `json:"caddy"`
}

func ParseConfig() (*Config, error) {
	var config Config
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return nil, err
	}
	err = toml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
