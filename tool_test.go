package tool_test

import (
	"context"
	"testing"

	tool "github.com/benaskins/axon-tool"
)

func TestToolDefExecute(t *testing.T) {
	def := tool.ToolDef{
		Name:        "greet",
		Description: "Says hello",
		Execute: func(ctx *tool.ToolContext, args map[string]any) tool.ToolResult {
			name, _ := args["name"].(string)
			return tool.ToolResult{Content: "Hello, " + name + "!"}
		},
	}

	tc := &tool.ToolContext{
		Ctx:    context.Background(),
		UserID: "user-1",
	}

	result := def.Execute(tc, map[string]any{"name": "World"})

	if result.Content != "Hello, World!" {
		t.Errorf("got %q, want %q", result.Content, "Hello, World!")
	}
}

func TestToolContextCarriesMetadata(t *testing.T) {
	tc := &tool.ToolContext{
		Ctx:            context.Background(),
		UserID:         "user-1",
		Username:       "alice",
		AgentSlug:      "helper",
		ConversationID: "conv-1",
		SystemPrompt:   "You are helpful.",
	}

	if tc.UserID != "user-1" {
		t.Errorf("UserID = %q, want %q", tc.UserID, "user-1")
	}
	if tc.AgentSlug != "helper" {
		t.Errorf("AgentSlug = %q, want %q", tc.AgentSlug, "helper")
	}
}

func TestHasRequiredParam(t *testing.T) {
	def := tool.ToolDef{
		Name:        "search",
		Description: "Search",
		Parameters: tool.ParameterSchema{
			Type:     "object",
			Required: []string{"query", "limit"},
		},
	}

	if !def.HasRequiredParam("query") {
		t.Error("expected query to be required")
	}
	if !def.HasRequiredParam("limit") {
		t.Error("expected limit to be required")
	}
	if def.HasRequiredParam("offset") {
		t.Error("expected offset to not be required")
	}
}

func TestToolDefWithParameters(t *testing.T) {
	params := tool.ParameterSchema{
		Type:     "object",
		Required: []string{"query"},
		Properties: map[string]tool.PropertySchema{
			"query": {
				Type:        "string",
				Description: "The search query",
			},
		},
	}

	def := tool.ToolDef{
		Name:        "search",
		Description: "Search the web",
		Parameters:  params,
		Execute: func(ctx *tool.ToolContext, args map[string]any) tool.ToolResult {
			return tool.ToolResult{Content: "results"}
		},
	}

	if def.Parameters.Required[0] != "query" {
		t.Errorf("expected 'query' in required, got %v", def.Parameters.Required)
	}
	if def.Parameters.Properties["query"].Type != "string" {
		t.Errorf("expected string type for query property")
	}
}
