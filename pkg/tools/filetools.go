package tools

import (
	"goflow/pkg/components"
	"fmt"
	"os"
	"encoding/json"
)


// File tool input structures
type ReadFileInput struct {
    Path string `json:"path"`
}

type WriteFileInput struct {
    Path string `json:"path"`
    Data string `json:"data"`
}


func FileReadTool() components.Tool {
	return components.Tool{
		Name:        "readFile",
		Description: "Reads content from a file at the specified path",
		Inputs:      ReadFileInput{},
		HandlerFunc: handleReadFile,
	}
}


func FileWriteTool() components.Tool {
	return components.Tool{
		Name:        "writeFile",
		Description: "Writes content to a file at the specified path",
		Inputs:      WriteFileInput{},
		HandlerFunc: handleWriteFile,
	}
}





func handleReadFile(inputs interface{}) (interface{}, error) {
    var input ReadFileInput
    inputJSON, err := json.Marshal(inputs)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal input: %v", err)
    }
    
    if err := json.Unmarshal(inputJSON, &input); err != nil {
        return nil, fmt.Errorf("failed to parse input: %v", err)
    }

    content, err := os.ReadFile(input.Path)
    if err != nil {
        return nil, fmt.Errorf("failed to read file: %v", err)
    }

    return string(content), nil
}

func handleWriteFile(inputs interface{}) (interface{}, error) {
    var input WriteFileInput
    inputJSON, err := json.Marshal(inputs)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal input: %v", err)
    }
    
    if err := json.Unmarshal(inputJSON, &input); err != nil {
        return nil, fmt.Errorf("failed to parse input: %v", err)
    }

    if err := os.WriteFile(input.Path, []byte(input.Data), 0644); err != nil {
        return nil, fmt.Errorf("failed to write file: %v", err)
    }

    return map[string]string{"status": "success", "path": input.Path}, nil
}