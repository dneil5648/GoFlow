package openai

import (
    "context"
    "fmt"
    "github.com/openai/openai-go"
    gf "goflow/pkg/components"
)

var supportedModels = map[string]bool{
    "gpt-4o":              true,
    "gpt-4":              true,
    "gpt-4-1106-preview": true,
    "gpt-4-vision-preview": true,
    "gpt-3.5-turbo":     true,
}

type OpenAIClient struct {
    client    *openai.Client
    config    gf.ClientConfig
    modelInfo gf.ModelInfo
}

func NewOpenAIClient(config gf.ClientConfig) (*OpenAIClient, error) {
    if err := validateModel(config.Model); err != nil {
        return nil, err
    }

    client := openai.NewClient()
    
    return &OpenAIClient{
        client: client,
        config: config,
        modelInfo: gf.ModelInfo{
            Provider: "openai",
            Model:    config.Model,
            MaxTokens: getModelMaxTokens(config.Model),
            Capabilities: map[string]bool{
                "functions": isModelFunctionCapable(config.Model),
                "vision":    config.Model == "gpt-4-vision-preview",
            },
        },
    }, nil
}

func (c *OpenAIClient) Generate(ctx context.Context, prompt gf.Prompt) (string, error) {
    messages := []openai.ChatCompletionMessageParamUnion{
        openai.SystemMessage(prompt.SystemMessage),
        openai.UserMessage(prompt.UserMessage),
    }

    completion, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
        Messages: openai.F(messages),
        Model:    openai.F(c.modelInfo.Model),
        Temperature: openai.Float(c.config.Temperature),
        MaxTokens: openai.Int(int64(c.config.MaxTokens)),
    })
    if err != nil {
        return "", fmt.Errorf("openai generation failed: %w", err)
    }

    return completion.Choices[0].Message.Content, nil
}

// Helper functions remain the same
func validateModel(model string) error {
    if !supportedModels[model] {
        return fmt.Errorf("unsupported model: %s", model)
    }
    return nil
}

// Update getModelMaxTokens to return int
func getModelMaxTokens(model string) int64 {
    switch model {
    case "gpt-4o":
        return 8192
    case "gpt-4":
        return 8192
    case "gpt-4-1106-preview":
        return 128000
    case "gpt-4-vision-preview":
        return 128000
    case "gpt-3.5-turbo":
        return 4096
    default:
        return 4096
    }
}

func isModelFunctionCapable(model string) bool {
    return model == "gpt-4" || model == "gpt-4-1106-preview" || model == "gpt-3.5-turbo"
}

func (c *OpenAIClient) GetModelInfo() gf.ModelInfo {
    return c.modelInfo
}

func (c *OpenAIClient) ValidateResponse(response string) error {
    if response == "" {
        return fmt.Errorf("empty response from OpenAI")
    }
    return nil
}