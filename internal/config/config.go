package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the configuration values.
type Config struct {
	HuggingFaceURL   string
	HuggingFaceToken string
}

// LoadConfig reads configuration from a .env file and environment variables.
func LoadConfig() (*Config, error) {
	// Load variables from .env file, if it exists.
	// Ignore the error if the .env file is not present.
	_ = godotenv.Load()

	// Fetch required environment variables.
	hfURL := os.Getenv("HUGGING_FACE_URL")
	hfToken := os.Getenv("HUGGING_FACE_TOKEN")

	// Validate that required variables are set.
	var missing []string
	if hfURL == "" {
		missing = append(missing, "HUGGING_FACE_URL")
	}
	if hfToken == "" {
		missing = append(missing, "HUGGING_FACE_TOKEN")
	}
	if len(missing) > 0 {
		return nil, errors.New("missing required environment variables: " + joinStrings(missing, ", "))
	}

	// Initialize and return the Config struct.
	cfg := &Config{
		HuggingFaceURL:   hfURL,
		HuggingFaceToken: hfToken,
	}
	return cfg, nil
}

// joinStrings joins a slice of strings with a given separator.
func joinStrings(items []string, sep string) string {
	result := ""
	for i, item := range items {
		if i > 0 {
			result += sep
		}
		result += item
	}
	return result
}
