package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type (
	Server struct {
		Port string `yaml:"port"`
	}
	Database struct {
		URL string `yaml:"url"`
	}
	Logging struct {
		Level string `yaml:"level"`
	}
	Telemetry struct {
		Enabled  bool   `yaml:"enabled"`
		Endpoint string `yaml:"endpoint"`
	}
	AppConfig struct {
		Env       string    `yaml:"env"`
		Service   string    `yaml:"service"`
		Version   string    `yaml:"version"`
		Server    Server    `yaml:"server"`
		Database  Database  `yaml:"database"`
		Logging   Logging   `yaml:"logging"`
		Telemetry Telemetry `yaml:"telemetry"`
	}
)

func Load(env string) (*AppConfig, error) {
	var configFile map[string]*AppConfig
	file, _ := os.Open("config.yml")
	defer file.Close() //nolint:errcheck
	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&configFile); err != nil {
		return nil, err
	}

	appConfig, ok := configFile[env]
	if !ok {
		return nil, fmt.Errorf("no such environment: %s", env)
	}

	return appConfig, nil
}
