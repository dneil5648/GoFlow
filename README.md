# GoFlow

GoFlow is a powerful, flexible Go framework for building AI-powered workflows using Large Language Models (LLMs). It provides a structured way to create, manage, and execute LLM-based tasks with features like state management, tool integration, and structured output parsing.

## Features

- 🔄 **Workflow Management**: Create and execute structured workflows with LLMs
- 🎯 **Schema Validation**: Enforce structured outputs with JSON schema validation
- 🛠️ **Tool Integration**: Integrate custom tools and functions into your LLM workflows
- 📝 **State Management**: Track and manage workflow state across executions
- 🔌 **Modular Design**: Easily extend with new LLM providers and tools
- 📊 **Logging**: Built-in logging capabilities for workflow monitoring
- 🔄 **Context Management**: Handle context and variables in prompts

## Installation

```bash
go get github.com/yourusername/goflow
```

## Quick Start

```go
package main

import (
    "fmt"
    "goflow/pkg/components"
    "goflow/pkg/llms/openai"
)

func main() {
    // 1. Define your schema
    schemaFields := []components.SchemaField{
        {
            Field:       "question",
            Description: "The question that was asked.",
            Type:        "string",
            Required:    true,
        },
        {
            Field:       "answer",
            Description: "The answer to the question",
            Type:        "string",
            Required:    true,
        },
    }

    // 2. Configure your LLM client
    clientConfig := components.ClientConfig{ 
        Model:       "gpt-4",
        Temperature: 0.7,
        MaxTokens:   1000,
    }

    // 3. Create the client
    client, err := openai.NewOpenAIClient(clientConfig)
    if err != nil {
        log.Fatal(err)
    }

    // 4. Create and run your workflow
    result, err := basicFlow(client, 
        "You are a helpful assistant",
        "What is the capital of France?",
        schemaFields,
    )
}
```

## Core Components

### Workflows

Workflows are the central concept in GoFlow. They orchestrate the execution of LLM operations with:
- Input/output schema validation
- Error handling
- Logging
- Tool integration

### Prompts

The `Prompt` structure manages:
- System messages
- User messages
- Variable substitution
- Output format requirements

```go
prompt := components.Prompt{
    SystemMessage: "You are a helpful assistant",
    UserMessage:  "Answer this: {{question}}",
    Variables: map[string]interface{}{
        "question": "What is GoFlow?",
    },
    OutputFormat: components.OutputFormat{
        Type: "json",
        Schema: schema.Build(),
    },
}
```

### State Management

Track workflow state with the built-in state management system:

```go
state := components.NewFlowState()
state.Add(map[string]interface{}{
    "step": 1,
    "data": "initial data",
})
```

### Tools

Integrate custom tools into your workflows:

```go
tool := components.Tool{
    Name:        "calculator",
    Description: "Performs basic calculations",
    HandlerFunc: func(inputs interface{}) (interface{}, error) {
        // Tool implementation
        return result, nil
    },
}
```

## Project Structure

goflow/
├── pkg/
│   ├── components/        # Core components
│   │   ├── llm.go        # LLM interface definitions
│   │   ├── logging.go    # Logging functionality
│   │   ├── outputs.go    # Output parsing and schemas
│   │   ├── prompts.go    # Prompt management
│   │   ├── state.go      # State management
│   │   ├── tools.go      # Tool definitions
│   │   └── workflow.go   # Workflow orchestration
│   ├── flows/            # Shared workflow implementations
│   ├── llms/
│   │   └── openai/       # OpenAI implementation
│   └── prompts/          # Shared prompt templates
└── main.go               # Example usage

## Advanced Usage

### Context-Aware Workflows

Use context in your workflows:

```go
docContext := map[string]interface{}{
    "context": "Your context here",
}
result, err := contextFlow(client, systemMessage, userMessage, schemaFields, docContext)
```

### Custom Output Schemas

Define custom output schemas:

```go
schemaFields := []components.SchemaField{
    {
        Field:       "analysis",
        Description: "Detailed analysis of the input",
        Type:        "string",
        Required:    true,
    },
    {
        Field:       "nextSteps",
        Description: "Recommended next steps",
        Type:        "string",
        Required:    true,
    },
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT

## Acknowledgments

- OpenAI for GPT API
