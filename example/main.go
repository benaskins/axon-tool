// Example: defining and executing a simple tool with axon-tool.
package main

import (
	"context"
	"fmt"

	tool "github.com/benaskins/axon-tool"
)

func main() {
	// Define a tool with a name, description, parameter schema, and execute function.
	greet := tool.ToolDef{
		Name:        "greet",
		Description: "Greet someone by name.",
		Parameters: tool.ParameterSchema{
			Type:     "object",
			Required: []string{"name"},
			Properties: map[string]tool.PropertySchema{
				"name": {Type: "string", Description: "The person to greet."},
			},
		},
		Execute: func(ctx *tool.ToolContext, args map[string]any) tool.ToolResult {
			name, _ := args["name"].(string)
			return tool.ToolResult{Content: fmt.Sprintf("Hello, %s!", name)}
		},
	}

	// Execute the tool with a context and arguments.
	tc := &tool.ToolContext{Ctx: context.Background()}
	result := greet.Execute(tc, map[string]any{"name": "World"})
	fmt.Println(result.Content)
}
