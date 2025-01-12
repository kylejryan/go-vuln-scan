package analyzer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/kylejryan/go-vuln-scan/internal/config"
	"github.com/kylejryan/go-vuln-scan/pkg/models"
)

type HFRequest struct {
	Inputs string `json:"inputs"`
}

func AnalyzeCodeWithHF(code string) (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(`You are a code security reviewer.
Identify potential security vulnerabilities or risky patterns in this code, and provide a short explanation along with a severity level:

%s
`, code)

	reqBody := HFRequest{
		Inputs: prompt,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", cfg.HuggingFaceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.HuggingFaceToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("Hugging Face API error. Status: %d, Body: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var hfResponses []models.HFResponse
	if err := json.Unmarshal(bodyBytes, &hfResponses); err != nil {
		// If not an array, try single object
		var singleResponse models.HFResponse
		if err := json.Unmarshal(bodyBytes, &singleResponse); err != nil {
			// Return raw text as fallback
			return string(bodyBytes), nil
		}
		return singleResponse.GeneratedText, nil
	}

	if len(hfResponses) > 0 {
		return hfResponses[0].GeneratedText, nil
	}

	return "", errors.New("no response from Hugging Face model")
}

func Analyze(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	codeSnippet := string(content)
	if len(codeSnippet) > 3000 {
		codeSnippet = codeSnippet[:3000] // Simple truncation
	}

	analysis, err := AnalyzeCodeWithHF(codeSnippet)
	if err != nil {
		return "", err
	}

	// Clean up analysis text
	analysis = strings.TrimSpace(analysis)
	return analysis, nil
}
