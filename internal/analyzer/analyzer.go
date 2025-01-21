package analyzer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kylejryan/go-vuln-scan/internal/config"
)

// ChatMessage represents a single message in a chat.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest is the payload we send to the Hugging Face endpoint.
type ChatCompletionRequest struct {
	Model     string        `json:"model"`
	Messages  []ChatMessage `json:"messages"`
	MaxTokens int           `json:"max_tokens"`
	Stream    bool          `json:"stream"`
}

// Choice and ChatCompletionResponse help parse the returned JSON structure.
type Choice struct {
	Index   int `json:"index"`
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	FinishReason string `json:"finish_reason"`
}

type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Choices []Choice `json:"choices"`
}

// Example function that asks "What is the capital of France?"
func Analyze(code string) (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}

	// Build our request payload
	payload := ChatCompletionRequest{
		Model: "meta-llama/Llama-3.3-70B-Instruct",
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: "What is the capital of France?",
			},
		},
		MaxTokens: 500,
		Stream:    false,
	}

	// Convert to JSON
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Create HTTP POST request
	req, err := http.NewRequest("POST",
		"https://api-inference.huggingface.co/models/meta-llama/Llama-3.3-70B-Instruct/v1/chat/completions",
		bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.HuggingFaceToken))

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("hugging face API error. Status: %d, Body: %s",
			resp.StatusCode, string(bodyBytes))
	}

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse JSON into our ChatCompletionResponse struct
	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(bodyBytes, &chatResp); err != nil {
		return "", err
	}

	// Return first message from the assistant, if available
	if len(chatResp.Choices) > 0 {
		return strings.TrimSpace(chatResp.Choices[0].Message.Content), nil
	}

	return "", nil
}
