package tool_test

import (
	"context"
	"fmt"

	tool "github.com/benaskins/axon-tool"
)

func ExampleToolDef() {
	greet := tool.ToolDef{
		Name:        "greet",
		Description: "Greet a user by name",
		Parameters: tool.ParameterSchema{
			Type:     "object",
			Required: []string{"name"},
			Properties: map[string]tool.PropertySchema{
				"name": {
					Type:        "string",
					Description: "The name of the person to greet",
				},
			},
		},
		Execute: func(ctx *tool.ToolContext, args map[string]any) tool.ToolResult {
			name, _ := args["name"].(string)
			return tool.ToolResult{Content: "Hello, " + name + "!"}
		},
	}

	fmt.Println(greet.Name)
	fmt.Println(greet.Description)
	// Output:
	// greet
	// Greet a user by name
}

func ExampleToolDef_Execute() {
	add := tool.ToolDef{
		Name:        "add",
		Description: "Add two numbers",
		Parameters: tool.ParameterSchema{
			Type:     "object",
			Required: []string{"a", "b"},
			Properties: map[string]tool.PropertySchema{
				"a": {Type: "number", Description: "First number"},
				"b": {Type: "number", Description: "Second number"},
			},
		},
		Execute: func(ctx *tool.ToolContext, args map[string]any) tool.ToolResult {
			a, _ := args["a"].(float64)
			b, _ := args["b"].(float64)
			return tool.ToolResult{Content: fmt.Sprintf("%.0f", a+b)}
		},
	}

	tc := &tool.ToolContext{
		Ctx:    context.Background(),
		UserID: "user-1",
	}

	result := add.Execute(tc, map[string]any{"a": 3.0, "b": 4.0})
	fmt.Println(result.Content)
	// Output: 7
}
