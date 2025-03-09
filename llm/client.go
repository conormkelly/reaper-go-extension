package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-reaper/pkg/logger"
	"io"
	"net/http"
	"time"
)

// Client defines the interface for LLM services
type Client interface {
	// SendPrompt sends a system prompt and user prompt to the LLM service
	// and returns the response text or an error
	SendPrompt(systemPrompt, userPrompt string) (string, error)
}

// Constants for OpenAI API
const (
	OpenAICompletionURL = "https://api.openai.com/v1/chat/completions"
	DefaultModel        = "gpt-3.5-turbo"
	DefaultMaxTokens    = 1024
	DefaultTemp         = 0.7
	DefaultTimeoutSec   = 30
)

// Message represents a single message in a conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIClient implements the Client interface for OpenAI
type OpenAIClient struct {
	APIKey     string
	Model      string
	MaxTokens  int
	Temp       float64
	HTTPClient *http.Client
}

// NewOpenAIClient creates a new OpenAI client with default settings
func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		APIKey:    apiKey,
		Model:     DefaultModel,
		MaxTokens: DefaultMaxTokens,
		Temp:      DefaultTemp,
		HTTPClient: &http.Client{
			Timeout: time.Duration(DefaultTimeoutSec) * time.Second,
		},
	}
}

// SendPrompt implements the Client interface
func (c *OpenAIClient) SendPrompt(systemPrompt, userPrompt string) (string, error) {
	// Log the start of the API call
	logger.Debug("Starting OpenAI API call...")

	// Set up the request body
	type RequestBody struct {
		Model       string    `json:"model"`
		Messages    []Message `json:"messages"`
		MaxTokens   int       `json:"max_tokens"`
		Temperature float64   `json:"temperature"`
	}

	reqBody := RequestBody{
		Model: c.Model,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		MaxTokens:   c.MaxTokens,
		Temperature: c.Temp,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request body: %v", err)
	}

	logger.Debug("Request body prepared, creating HTTP request...")

	// Create the request
	req, err := http.NewRequest("POST", OpenAICompletionURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	logger.Debug("Sending HTTP request to OpenAI API...")

	// Send the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	// Guard against nil response
	if resp == nil {
		return "", fmt.Errorf("nil response received from HTTP client")
	}
	defer resp.Body.Close()

	logger.Debug("Received response with status: %d", resp.StatusCode)

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		logger.Error("API error response: %s", string(body))
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	logger.Debug("Successfully read response body (%d bytes)", len(body))

	// Parse the response
	var openAIResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &openAIResp); err != nil {
		logger.Error("Error parsing response: %v", err)
		logger.Error("Response content: %s", string(body))
		return "", fmt.Errorf("error parsing API response: %v", err)
	}

	// Check for API error
	if openAIResp.Error != nil && openAIResp.Error.Message != "" {
		return "", fmt.Errorf("API error: %s", openAIResp.Error.Message)
	}

	// Check for valid choices
	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned from API")
	}

	// Check for nil or empty content
	content := openAIResp.Choices[0].Message.Content
	if content == "" {
		return "", fmt.Errorf("empty content in API response")
	}

	logger.Debug("Successfully parsed OpenAI response")
	return content, nil
}
