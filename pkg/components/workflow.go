package components

import (
    "context"
    "fmt"
    "time"
    "encoding/json"
)

type WorkFlowType int

const (
    WorkFlowDo WorkFlowType = iota
    WorkFlowChoose
)

type WorkFlow struct {
    Name         string
    Type         WorkFlowType
    Client       LLMClient
    OutputParser OutputParser
    Config       WorkflowConfig
    Prompt       Prompt
    Tools        *ToolList
    Logger       *Logger
}

type WorkflowConfig struct {
    MaxRetries   int
    Timeout      time.Duration
    Temperature  float64
}



func (wf *WorkFlow) Run(ctx context.Context) (interface{}, error) {
    // Log start of workflow
    wf.Logger.LogItem(wf.Name, "Starting workflow execution")
    
    // Generate LLM response
    response, err := wf.Client.Generate(ctx, wf.Prompt)
    if err != nil {
        wf.Logger.LogItem(wf.Name, fmt.Sprintf("Error generating response: %v", err))
        return nil, fmt.Errorf("LLM generation failed: %w", err)
    }
    
    // Parse the response based on workflow type
    switch wf.Type {
    case WorkFlowDo:
        result, err := wf.OutputParser.Parse(response)
        if err != nil {
            wf.Logger.LogItem(wf.Name, fmt.Sprintf("Error parsing response: %v", err))
            return nil, fmt.Errorf("output parsing failed: %w", err)
        }
        return result, nil
        
    case WorkFlowChoose:
        var toolSelection ToolSelectionOutput
        if err := json.Unmarshal([]byte(response), &toolSelection); err != nil {
            wf.Logger.LogItem(wf.Name, fmt.Sprintf("Error parsing tool selection: %v", err))
            return nil, fmt.Errorf("tool selection parsing failed: %w", err)
        }
        
        tool, exists := wf.Tools.Tools[toolSelection.ToolName]
        if !exists {
            return nil, fmt.Errorf("selected tool %s not found", toolSelection.ToolName)
        }
        
        tool.Inputs = toolSelection.ToolInputs
        result, err := tool.Run()
        if err != nil {
            wf.Logger.LogItem(wf.Name, fmt.Sprintf("Error running tool: %v", err))
            return nil, fmt.Errorf("tool execution failed: %w", err)
        }
        return result, nil
    }
    
    return nil, fmt.Errorf("invalid workflow type")
}

func NewWorkflow(
    name string,
    workflowType WorkFlowType,
    client LLMClient,
    outputParser OutputParser,
    config WorkflowConfig,
    prompt Prompt,
    tools *ToolList,
    logger *Logger,
) (*WorkFlow, error) {
    if client == nil {
        return nil, fmt.Errorf("client cannot be nil")
    }
    
    return &WorkFlow{
        Name:         name,
        Type:         workflowType,
        Client:       client,
        OutputParser: outputParser,
        Config:       config,
        Prompt:       prompt,
        Tools:        tools,
        Logger:       logger,
    }, nil
}