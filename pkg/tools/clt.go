package tools

import (

    "fmt"
    "os/exec"
	"goflow/pkg/components"
	"encoding/json"
)

type WhoisInput struct {
    Domain string`json:"domain"`
}

func CreateWhoisTool() components.Tool {
    return components.Tool{
        Name:        "whois",
        Description: "Performs a whois lookup for a domain. This tool will only provide new information when looking up a unique input.",
        Inputs:      WhoisInput{},
        HandlerFunc: handleWhois,
    }
}

func handleWhois(inputs interface{}) (interface{}, error) {
    
	var inputValue WhoisInput
    json.Unmarshal([]byte(inputs.(string)), &inputValue)

    cmd := exec.Command("whois", inputValue.Domain)
    output, err := cmd.CombinedOutput()
	
    if err != nil {
        return nil, fmt.Errorf("whois command failed: %v", err)
    }

    return string(output), nil
}