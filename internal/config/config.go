package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Name   string `json:"name"`
	Db_url string `json:"db_url"`
}

func Read() (Config, error) {
	var cfg Config

	path, err := getConfigFilePath()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, fmt.Errorf("Error reading the config file: %w", err)
	}

	if len(data) == 0 {
		return cfg, nil
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("Error while unmarshaling config file in Read: %w", err)
	}

	return cfg, nil
}

func getConfigFilePath() (path string, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Could not get user home direcotry: %w", err)
	}
	return filepath.Join(homeDir, configFileName), nil
}

func (c *Config) Write() error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
