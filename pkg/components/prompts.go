package components


import (
	"fmt"
	"strings"
	"encoding/json"
)

type Prompt struct {
    SystemMessage string
    UserMessage   string
    Variables     map[string]interface{}
    Tools         *ToolList
    OutputFormat  OutputFormat    // Add explicit output format
}

type OutputFormat struct {
    Type        string         // e.g., "json", "text"
    Schema      interface{}    // for JSON schema definition
    Description string         // human readable description
}

// FormatPrompt creates the final prompt with variables replaced and output requirements
func (p *Prompt) FormatPrompt() (string, string) {
    // Format system message with output requirements
    systemMsg := p.SystemMessage
    if p.OutputFormat.Type == "json" {
        systemMsg += fmt.Sprintf("\nYou must return a JSON object in the following format. %s", p.OutputFormat.Description)
        if schema, ok := p.OutputFormat.Schema.(map[string]interface{}); ok {
            schemaStr, _ := json.MarshalIndent(schema, "", "  ")
            systemMsg += fmt.Sprintf("\nUse this JSON schema: %s", string(schemaStr))
        }
    }

    // Replace variables in user message
    userMsg := p.UserMessage
    for key, value := range p.Variables {
        placeholder := fmt.Sprintf("{{%s}}", key)
        userMsg = strings.ReplaceAll(userMsg, placeholder, fmt.Sprintf("%v", value))
    }
	p.SystemMessage = systemMsg
	p.UserMessage = userMsg
    return systemMsg, userMsg
}


func (p *Prompt) AddTools() error {
    if p.Tools == nil {
        return nil
    }

    toolsDescription := "\nAvailable tools:\n"
    for name, tool := range p.Tools.Tools {
        toolsDescription += fmt.Sprintf("- %s: %s\n", name, tool.Description)
        // Add tool inputs schema if available
        if tool.Inputs != nil {
            inputsJSON, err := json.MarshalIndent(tool.Inputs, "  ", "  ")
            if err != nil {
                return fmt.Errorf("failed to marshal tool inputs: %w", err)
            }
            toolsDescription += fmt.Sprintf("  Inputs schema: %s\n", string(inputsJSON))
        }
    }

    p.SystemMessage += toolsDescription
    return nil
}

func (p Prompt) GetTools() *ToolList {
    return p.Tools
}