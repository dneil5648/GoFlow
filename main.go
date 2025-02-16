package main

import (
    "fmt"
    "log"
    "os"
    "goflow/pkg/components"
    "goflow/pkg/llms/openai"
)

func main() {
    // 1. Define schema fields
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
        {
            Field:       "keyFact",
            Description: "The Key Facts of the article",
            Type:        "string",
            Required:    true,
        },
        {
            Field:       "nextStep",
            Description: "The next step you would like to take to complete the task.",
            Type:        "string",
            Required:    true,
        },
        {
            Field: "nextSystemPrompt",
            Description: "The next prompt system Prompt for the next workflow",
            Type:        "string",
            Required:    true,
        },
        {
            Field: "nextUserPrompt",
            Description: "The next prompt user Prompt for the next workflow",
            Type:        "string",
            Required:    true,
        },
    }
    
   

    //Set up the client Config and Client Struct
    clientConfig := components.ClientConfig{ 
        Model:       "gpt-4",
        Temperature: 0.1,
        MaxTokens:   1000,
    }

    client, err := openai.NewOpenAIClient(clientConfig)
    if err != nil {
        log.Fatalf("Failed to create OpenAI client: %v", err)
    }


    // Create the Prompts
	var systemMessage string = `You are a helpful ai who can analyze data and provide a response. To provide `
	var userMessage string = `Please provide me a summary of this document and pull out key points.\n Please use the following context:\n {{context}}`

    content, err := os.ReadFile("testfile.md")
    if err != nil {
        fmt.Printf("Error loading file: %v\n", err)
        return
    }

    fmt.Println(string(content))
    docContext := map[string]interface{}{
        "context": string(content),
    }
    result, err := contextFlow(client,systemMessage, userMessage, schemaFields, docContext)
    if err != nil{
        fmt.Printf("Error Running Flow: %v", err)
    }

    // // Run your flow
    // result, err := basicFlow(client,systemMessage, userMessage, schemaFields )
	// if err != nil {
	// 	fmt.Printf("Error Running Flow: %v", err)
	// }

	fmt.Println(result)
}