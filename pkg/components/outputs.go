package components

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Output interfaces
type OutputParser interface {
	Parse(input string) (interface{}, error)
	ValidateSchema(schema interface{}) error
}

// SchemaField represents a single field in the JSON schema
type SchemaField struct {
	Field       string `json:"field"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
}

// JSONSchemaBuilder simplifies schema creation
type JSONSchemaBuilder struct {
	Fields []SchemaField
}

type JSONParser struct {
	fields []SchemaField
}

func NewJSONParser(fields []SchemaField) *JSONParser {
	return &JSONParser{
		fields: fields,
	}
}

func (p *JSONParser) Parse(input string) (interface{}, error) {
	if input == "" {
		return nil, errors.New("input must not be empty")
	}

	// Parse into a map first
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(input), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// Check if we're getting a schema instead of content
	if _, hasProperties := result["properties"]; hasProperties {
		return nil, fmt.Errorf("received schema definition instead of content. Please provide actual data")
	}

	// Validate all required fields are present
	for _, field := range p.fields {
		if field.Required {
			if _, ok := result[field.Field]; !ok {
				return nil, fmt.Errorf("missing required field: %s", field.Field)
			}
		}
	}

	return result, nil
}

func (p *JSONParser) ValidateSchema(schema interface{}) error {
	return nil
}

// Creates the actual JSON schema from the fields
func (b *JSONSchemaBuilder) Build() map[string]interface{} {
	properties := make(map[string]interface{})

	for _, field := range b.Fields {
		properties[field.Field] = map[string]interface{}{
			"type":        field.Type,
			"description": field.Description,
		}
	}

	return properties
}
