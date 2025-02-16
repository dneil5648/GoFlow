package components

import "fmt"

type ToolType interface {
    Run() (interface{}, error)
}

type HandlerFunc func(inputs interface{}) (interface{}, error)

type Tool struct {
    Name         string      `json:"name"`
    Description  string      `json:"description"`
    Inputs       interface{} `json:"inputs"`
    HandlerFunc  HandlerFunc `json:"-"`
}

type ToolSelectionOutput struct {
    ToolName   string
    ToolInputs interface{}
}

type ToolList struct {
    Tools map[string]Tool  // Changed from 'tools' to 'Tools' for consistency
}

func (t Tool) Run() (interface{}, error) {
    output, err := t.HandlerFunc(t.Inputs)
    if err != nil {
        return nil, fmt.Errorf("error running tool: %v", err)
    }
    return output, nil
}