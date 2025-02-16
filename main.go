package main

import (
   
    "fmt"
    "log"
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
	var systemMessage string = `You are a helpful ai that can answer any question that is passed.`
	var userMessage string = `Please tell me well know facts about the berge khalifa `

    // Run your flow
    result, err := basicFlow(client,systemMessage, userMessage, schemaFields )
	if err != nil {
		fmt.Printf("Error Running Flow: %v", err)
	}

	fmt.Println(result)
}