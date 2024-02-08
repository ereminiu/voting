package config

import (
	"errors"
	"os"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	SSLMode  string `yaml:"sslmode"`
}

const (
	testPath = "file://config/test/config.yaml"
	prodPath = "file://config/prod/config.yaml"
)

func getPath(mode string) (string, error) {
	envPath := os.Getenv("CONFIG_PATH")
	if envPath != "" {
		return envPath, nil
	}
	if mode == "test" {
		return testPath, nil
	} else if mode == "prod" {
		return prodPath, nil
	}
	return "", errors.New("Unknown config type")
}

func LoadConfigs(mode string) (*Config, error) {
	path, err := getPath(mode)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
