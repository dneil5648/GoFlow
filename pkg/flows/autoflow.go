package flows

import (
	"context"
	"encoding/json"
	"fmt"
	"goflow/pkg/components"
	"goflow/pkg/llms/openai"
	"time"
	// "goflow/pkg/tools"
)

func CoTWorkFlow(client *openai.OpenAIClient, sysMessage string, uMessage string, fields []components.SchemaField, variables map[string]interface{}, tools *components.ToolList) (interface{}, error) {
	state := components.NewFlowState()
	currentMessage := uMessage
	maxSteps := 50
	workflowName := "Entry Workflow"

	for steps := 0; steps < maxSteps; steps++ {
		schemaFields := []components.SchemaField{
			{
				Field:       "tool_name",
				Description: "Name of the tool to use",
				Type:        "string",
				Required:    true,
			},
			{
				Field:       "tool_input",
				Description: "Input for the selected tool",
				Type:        "object",
				Required:    true,
			},
			{
				Field:       "isComplete",
				Description: "Whether the task is complete, must be 'true' or 'false'",
				Type:        "boolean",
				Required:    true,
			},
			{
				Field:       "nextStep",
				Description: "The next step to take if the task is not complete",
				Type:        "string",
				Required:    false,
			}, {
				Field:       "thought",
				Description: "This section is to capture your thoughts on the current task that can carry over to the next",
				Type:        "string",
				Required:    false,
			},
			{
				Field:       "workflowName",
				Description: "Name for the next workflow",
				Type:        "string",
				Required:    true,
			},
		}

		schema := &components.JSONSchemaBuilder{
			Fields: schemaFields,
		}

		result, err := runSingleStep(workflowName, client, sysMessage, currentMessage, schema, variables, tools)
		if err != nil {
			return nil, err
		}

		if err := state.Add(result); err != nil {
			return nil, fmt.Errorf("failed to add to state: %v", err)
		}

		isComplete, ok := result["isComplete"].(bool)
		if ok && isComplete {
			break
		}

		lastResult, err := state.GetLast()
		if err != nil {
			return nil, fmt.Errorf("failed to get last state: %v", err)
		}

		allSteps, err := state.Get()
		if err != nil {
			return nil, fmt.Errorf("failed to get state history: %v", err)
		}

		variables["previous_result"] = lastResult
		variables["workflow_history"] = allSteps

		nextQuestion, ok := result["nextQuestion"].(string)
		if !ok || nextQuestion == "" {
			break
		}

		nextWorkflowName, ok := result["workflowName"].(string)
		if !ok || nextWorkflowName == "" {
			break
		}

		currentMessage = nextQuestion
		workflowName = nextWorkflowName
	}

	allResults, err := state.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get final state: %v", err)
	}

	finalSysMessage := `You are a helpful analysis AI that can take all of the data gathered and provide accurate responses.\n Ensure that you are not returning schema definiton.`
	finalUserMessage := `Please finish your analysis and respond with the properly formatted JSON object from the provided schema.
Please use this context to complete the analysis: {{context}}`
	finalStepResult, err := runSingleStep(
		"Exit Workflow",
		client,
		finalSysMessage,
		finalUserMessage,
		&components.JSONSchemaBuilder{Fields: fields},
		map[string]interface{}{
			"context": allResults,
		},
	)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"steps":        allResults,
		"final_output": finalStepResult,
		"step_count":   len(allResults),
	}, nil
}

func runSingleStep(workflowName string, client *openai.OpenAIClient, sysMessage string, uMessage string, schema *components.JSONSchemaBuilder, variables map[string]interface{}, tools ...*components.ToolList) (map[string]interface{}, error) {
	parser := components.NewJSONParser(schema.Fields)

	var toolList *components.ToolList

	// If tools were provided, use the first one
	if len(tools) > 0 {
		toolList = tools[0]
	}

	prompt := components.Prompt{
		SystemMessage: sysMessage,
		UserMessage:   uMessage,
		Variables:     variables,
		Tools:         toolList,
		OutputFormat: components.OutputFormat{
			Type:        "json",
			Schema:      schema.Build(),
			Description: "Return a JSON object with the specified fields.",
		},
	}

	if len(tools) > 0 {
		prompt.Tools = tools[0]
	}

	prompt.FormatPrompt()
	if prompt.Tools != nil {
		prompt.AddTools()
	}

	// fmt.Printf("System Message: %v\n" ,prompt.SystemMessage)
	// fmt.Printf("User Message: %v\n" ,prompt.UserMessage)

	workflow, err := components.NewWorkflow(
		workflowName,
		components.WorkFlowDo,
		client,
		parser,
		components.WorkflowConfig{
			MaxRetries: 3,
			Timeout:    time.Second * 30,
		},
		prompt,
		nil,
		&components.Logger{LogFile: "workflow.log"},
	)

	if err != nil {
		return nil, fmt.Errorf("workflow creation failed: %v", err)
	}

	result, err := workflow.Run(context.Background())
	if err != nil {
		return nil, fmt.Errorf("workflow execution failed: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected result type")
	}
	if tools != nil {
		if toolName, ok := resultMap["tool_name"].(string); ok {
			if tool, exists := toolList.Tools[toolName]; exists {

				toolInput, err := json.Marshal(resultMap["tool_input"])
				if err != nil {
					return nil, fmt.Errorf("failed to marshal tool input: %v", err)
				}

				tool.Inputs = string(toolInput)
				fmt.Printf("Tool Input: %v\n", tool.Inputs)
				toolResult, err := tool.Run()
				if err != nil {
					return nil, fmt.Errorf("tool execution failed: %w", err)
				}
				resultMap["tool_output"] = toolResult
			}
		}
	}

	return resultMap, nil
}
