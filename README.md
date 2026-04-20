# axon-tool

> Primitives · Part of the [lamina](https://github.com/benaskins/lamina-mono) workspace

Tool definition and execution primitives for LLM agents. axon-tool provides a provider-agnostic way to declare tools with typed parameter schemas and execute them within a request-scoped context. The core value is the `ToolDef` / `ToolResult` contract that any agent framework can build on, plus helpers for parsing tool calls from streaming model output.

## Getting started

```
go get github.com/benaskins/axon-tool@latest
```

Requires Go 1.26+.

```go
package main

import (
	"context"
	"fmt"

	tool "github.com/benaskins/axon-tool"
)

func main() {
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

	tc := &tool.ToolContext{Ctx: context.Background()}
	result := greet.Execute(tc, map[string]any{"name": "World"})
	fmt.Println(result.Content)
}
```

See [`example/main.go`](example/main.go) for the runnable version.

## Key types

- **`ToolDef`**: tool definition with name, description, parameter schema, and execute function
- **`ToolResult`**: execution result containing text content
- **`ToolContext`**: request-scoped context carrying user ID, agent slug, conversation ID
- **`ParameterSchema`** / **`PropertySchema`**: JSON Schema for tool parameters
- **`TextGenerator`**: function type for sending a prompt to an LLM and getting text back
- **`ToolCall`**: provider-agnostic representation of a tool invocation parsed from model output
- **`ToolCallMatcher`**: [axon-tape](https://github.com/benaskins/axon-tape) matcher that detects tool-call JSON (bare or in a fenced code block) in streaming text

## License

MIT. See [LICENSE](LICENSE).
