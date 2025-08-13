package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type Config struct {
	OpenRouter struct {
		ApiKey string `json:"api_key"`
		Model  string `json:"model"`
	} `json:"openrouter"`
	Commands struct {
		AskToRun bool `json:"ask_to_run"`
	} `json:"commands"`
}

func loadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	var config Config

	// defaults
	config.OpenRouter.ApiKey = "your-api-key"
	config.OpenRouter.Model = "openai/gpt-5-chat"
	config.Commands.AskToRun = true

	// ~/.idk.yml
	path := filepath.Join(home, ".idk.yml")

	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		if os.IsNotExist(err) {
			b, err := yaml.Marshal(config)
			if err != nil {
				return nil, err
			}

			err = os.WriteFile(path, b, 0644)
			if err != nil {
				return nil, err
			}

			fmt.Printf("created config at %s\n", path)

			os.Exit(1)
		}

		return nil, err
	}

	defer file.Close()

	err = yaml.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
