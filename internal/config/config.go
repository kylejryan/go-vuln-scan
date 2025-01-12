package config

import (
	"errors"
	"os"
)

type Config struct {
	HuggingFaceToken string
	HuggingFaceURL   string
}

// Will switch to ChatGPT/Claude soon probably
func LoadConfig() (*Config, error) {
	token := os.Getenv("HF_API_TOKEN")
	if token == "" {
		return nil, errors.New("HF_API_TOKEN environment variable not set")
	}

	url := os.Getenv("HF_API_URL")
	if url == "" {
		return nil, errors.New("HF_API_URL environment variable not set")
	}

	return &Config{
		HuggingFaceToken: token,
		HuggingFaceURL:   url,
	}, nil
}
