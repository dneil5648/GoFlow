package components 

import "fmt"

// State interface defines the contract for state management
type State interface {
    Add(item map[string]interface{}) error
    Get() ([]interface{}, error)
    Clear() error
    GetLast() (interface{}, error)
}

// FlowState implements the State interface
type FlowState struct {
    Memory []interface{}
}

// NewFlowState creates a new FlowState instance
func NewFlowState() *FlowState {
    return &FlowState{
        Memory: make([]interface{}, 0),
    }
}

// Add appends an item to the memory
func (flow *FlowState) Add(item map[string]interface{}) error {
    flow.Memory = append(flow.Memory, item)
    return nil
}

// Get returns all items in memory
func (flow *FlowState) Get() ([]interface{}, error) {
    if flow.Memory == nil {
        return nil, fmt.Errorf("memory is not initialized")
    }
    return flow.Memory, nil
}

// GetLast returns the most recent item in memory
func (flow *FlowState) GetLast() (interface{}, error) {
    if len(flow.Memory) == 0 {
        return nil, fmt.Errorf("memory is empty")
    }
    return flow.Memory[len(flow.Memory)-1], nil
}

// Clear removes all items from memory
func (flow *FlowState) Clear() error {
    flow.Memory = make([]interface{}, 0)
    return nil
}