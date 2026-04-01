package tool

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/benaskins/axon-tape"
)

const toolCallMaxAccumulate = 2048

// ToolCall is a provider-agnostic tool call representation.
type ToolCall struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

// ToolCallAction is returned when a tool call is detected in the stream.
type ToolCallAction struct{ Calls []ToolCall }

func (ToolCallAction) IsFilterAction() {}

// ToolCallMatcher detects tool call JSON emitted as text content.
type ToolCallMatcher struct{}

func NewToolCallMatcher() *ToolCallMatcher { return &ToolCallMatcher{} }
func (m *ToolCallMatcher) Name() string    { return "tool_call" }

func (m *ToolCallMatcher) Scan(buf []byte, _ string) tape.MatchResult {
	trimmed := bytes.TrimSpace(buf)
	if len(trimmed) == 0 {
		return tape.NoMatch
	}

	trimmed = stripCodeFencePrefix(trimmed)

	if trimmed[0] != '{' && trimmed[0] != '[' {
		return tape.NoMatch
	}

	if len(trimmed) > toolCallMaxAccumulate {
		return tape.NoMatch
	}

	if tryParseToolCallJSON(trimmed) {
		return tape.FullMatch
	}

	if looksLikeToolCallStart(trimmed) {
		return tape.PartialMatch
	}

	return tape.NoMatch
}

func (m *ToolCallMatcher) Extract(buf []byte) tape.FilterAction {
	trimmed := bytes.TrimSpace(buf)
	trimmed = stripCodeFencePrefix(trimmed)
	trimmed = stripCodeFenceSuffix(trimmed)
	trimmed = bytes.TrimSpace(trimmed)

	calls := parseToolCallJSON(trimmed)
	if len(calls) > 0 {
		return ToolCallAction{Calls: calls}
	}
	return tape.ContinueAction{}
}

func stripCodeFencePrefix(b []byte) []byte {
	if !bytes.HasPrefix(b, []byte("```")) {
		return b
	}
	rest := b[3:]
	if idx := bytes.IndexByte(rest, '\n'); idx >= 0 {
		return bytes.TrimSpace(rest[idx+1:])
	}
	return bytes.TrimSpace(rest)
}

func stripCodeFenceSuffix(b []byte) []byte {
	if bytes.HasSuffix(b, []byte("```")) {
		return bytes.TrimSpace(b[:len(b)-3])
	}
	return b
}

func looksLikeToolCallStart(b []byte) bool {
	s := string(b)
	if strings.HasPrefix(s, "{") {
		if len(s) < 10 {
			return true
		}
		return strings.Contains(s, `"name"`)
	}
	if strings.HasPrefix(s, "[") {
		if len(s) < 12 {
			return true
		}
		return strings.Contains(s, `"name"`)
	}
	return false
}

func tryParseToolCallJSON(b []byte) bool {
	b = stripCodeFenceSuffix(b)
	b = bytes.TrimSpace(b)

	var single struct {
		Name string `json:"name"`
	}
	if json.Unmarshal(b, &single) == nil && single.Name != "" {
		return true
	}

	var arr []struct {
		Name string `json:"name"`
	}
	if json.Unmarshal(b, &arr) == nil && len(arr) > 0 && arr[0].Name != "" {
		return true
	}

	return false
}

func parseToolCallJSON(b []byte) []ToolCall {
	var single struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}
	if json.Unmarshal(b, &single) == nil && single.Name != "" {
		return []ToolCall{{
			Name:      single.Name,
			Arguments: single.Arguments,
		}}
	}

	var arr []struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}
	if json.Unmarshal(b, &arr) == nil && len(arr) > 0 {
		var calls []ToolCall
		for _, item := range arr {
			if item.Name == "" {
				continue
			}
			calls = append(calls, ToolCall{
				Name:      item.Name,
				Arguments: item.Arguments,
			})
		}
		return calls
	}

	return nil
}
