package main

import (
    "context"
    "fmt"
    "log"
    "time"
    "goflow/pkg/components"
    "goflow/pkg/llms/openai"
)



func basicFlow(client *openai.OpenAIClient, sysMessage string, uMessage string, fields []components.SchemaField) (interface{}, error){
    schemaFields := fields

    schema := &components.JSONSchemaBuilder{
        Fields: schemaFields,
    }

    //Create Parser From Schema
    parser := components.NewJSONParser(schemaFields)

    prompt := components.Prompt{
        SystemMessage: sysMessage,
        UserMessage:  uMessage,
        OutputFormat: components.OutputFormat{
            Type:        "json",
            Schema:      schema.Build(), // Call Build() to get the schema
            Description: "Return a JSON object with the specified fields.",
        },
    }
    // Format Prompt to include ouput schema 
    prompt.FormatPrompt()

    // 5. Create and run workflow
    workflow, err := components.NewWorkflow(
        "Test Assistant",
        components.WorkFlowDo,
        client,
        parser,  // Pass the parser directly, not a pointer to it
        components.WorkflowConfig{
            MaxRetries:  3,
            Timeout:    time.Second * 30,
        },
        prompt,
        nil,
        &components.Logger{LogFile: "workflow.log"},
    )
    if err != nil {
        log.Fatalf("Failed to create workflow: %v", err)
    }

    ctx := context.Background()
    result, err := workflow.Run(ctx)
    if err != nil {
        log.Fatalf("Workflow failed: %v", err)
    }

    // 6. Handle the result
    // Current problematic code
    if resultMap, ok := result.(map[string]interface{}); ok {
        return resultMap, nil
    } else {
        return nil, fmt.Errorf("Unexpected result type")
        // log.Fatalf("Unexpected result type: %T", result)
         // This will never execute
    }



}

