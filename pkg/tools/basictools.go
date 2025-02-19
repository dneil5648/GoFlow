package tools

import (
	"goflow/pkg/components"
)

func CreateTools(toolMap map[string]components.Tool) *components.ToolList {
    return &components.ToolList{
        Tools: toolMap,
    }
}









