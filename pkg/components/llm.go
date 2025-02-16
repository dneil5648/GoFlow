package components

import (
    "context"
    "time"
)

type LLMClient interface {
    Generate(ctx context.Context, prompt Prompt) (string, error)
    GetModelInfo() ModelInfo
    ValidateResponse(response string) error
}

type ModelInfo struct {
    Provider     string
    Model        string
    MaxTokens    int64
    Capabilities map[string]bool
}

type ClientConfig struct {
    APIKey       string
    BaseURL      string
    Timeout      time.Duration
    MaxRetries   int
    Temperature  float64
    Model        string
	MaxTokens    int64 
}