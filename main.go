package main

import (
	"fmt"
	"log"

	"goflow/pkg/components"
	"goflow/pkg/flows"
	"goflow/pkg/llms/openai"
	"goflow/pkg/tools"
)

func main() {
	// 1. Define final output schema fields
	finalSchema := []components.SchemaField{
		{
			Field:       "domain",
			Description: "The domain being analyzed",
			Type:        "string",
			Required:    true,
		},
		{
			Field:       "analysis",
			Description: "Complete analysis of the domain",
			Type:        "string",
			Required:    true,
		},
		{
			Field:       "security_posture",
			Description: "Overall security assessment",
			Type:        "string",
			Required:    true,
		},
		{
			Field:       "recommendations",
			Description: "Recommended actions",
			Type:        "string",
			Required:    true,
		},
	}

	// 2. Set up tools
	customTools := map[string]components.Tool{
		"whois": tools.CreateWhoisTool(),
		// Add other tools as needed
	}
	customToolList := tools.CreateTools(customTools)

	// 3. Client setup
	clientConfig := components.ClientConfig{
		Model:       "gpt-4",
		Temperature: 0.1,
		MaxTokens:   1000,
	}

	client, err := openai.NewOpenAIClient(clientConfig)
	if err != nil {
		log.Fatalf("Failed to create OpenAI client: %v", err)
	}

	// 4. Create the Prompts
	systemMessage := `You are an AI analyst specialized in domain analysis. Your role is to:
                1. First, examine the tools available to you. These will be your only source of data.
                2. Plan your analysis based ONLY on the tools you have access to.
                3. For each tool:
                - State what information you plan to gather
                - Use the tool to collect the data
                - Analyze the results
                4. Once you have exhausted all available tools:
                - Summarize all gathered data
                - Provide analysis based ONLY on the information collected from these tools
                - Do not make assumptions about data you cannot verify with your tools
                - Clearly state if there are important security aspects you cannot assess due to tool limitations

                Important: 
                - Do not attempt to access external resources or tools not explicitly provided
                - If you need information but don't have the appropriate tool, note this in your analysis
                - Structure your findings based solely on verifiable data from your available tools

                Return your analysis in the specified JSON format once you have completed your investigation or can no longer gather more data with your current tools.`

	userMessage := `Please analyze the security posture of this domain: {{domain}}.Start by gathering basic information and then dig deeper based on what you find.`

	inputContext := map[string]interface{}{
		"domain": "recordedfuture.com",
	}

	// 5. Run the CoT Workflow
	result, err := flows.CoTWorkFlow(
		client,
		systemMessage,
		userMessage,
		finalSchema,
		inputContext,
		customToolList,
	)
	if err != nil {
		log.Fatalf("Error Running Flow: %v", err)
	}

	// 6. Process results
	data, ok := result.(map[string]interface{})
	if !ok {
		log.Fatalf("Invalid result type")
	}

	// Print step-by-step analysis
	steps, ok := data["steps"].([]interface{})
	if ok {
		fmt.Println("\nAnalysis Steps:")
		for i, step := range steps {
			stepData, ok := step.(map[string]interface{})
			if ok {
				fmt.Printf("\nStep %d: %s\n", i+1, stepData["workflowName"])
				if toolOutput, exists := stepData["tool_output"]; exists {
					fmt.Printf("Tool Output: %v\n", toolOutput)
				}
			}
		}
	}

	// Print final analysis
	if finalOutput, ok := data["final_output"].(map[string]interface{}); ok {
		fmt.Println("\nFinal Analysis:")
		fmt.Printf("Domain: %s\n", finalOutput["domain"])
		fmt.Printf("Analysis: %s\n", finalOutput["analysis"])
		fmt.Printf("Security Posture: %s\n", finalOutput["security_posture"])
		fmt.Printf("Recommendations: %s\n", finalOutput["recommendations"])
	}
}
