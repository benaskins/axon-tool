package tool

import (
	"strings"
	"testing"

	"github.com/benaskins/axon-tape"
)

func TestToolCallMatcher_SingleObject(t *testing.T) {
	m := NewToolCallMatcher()
	buf := []byte(`{"name": "web_search", "arguments": {"query": "golang"}}`)

	if r := m.Scan(buf, ""); r != tape.FullMatch {
		t.Fatalf("expected FullMatch, got %v", r)
	}

	action := m.Extract(buf)
	tc, ok := action.(ToolCallAction)
	if !ok {
		t.Fatalf("expected ToolCallAction, got %T", action)
	}
	if len(tc.Calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(tc.Calls))
	}
	if tc.Calls[0].Name != "web_search" {
		t.Errorf("expected 'web_search', got %q", tc.Calls[0].Name)
	}
	if tc.Calls[0].Arguments["query"] != "golang" {
		t.Errorf("expected query 'golang', got %v", tc.Calls[0].Arguments["query"])
	}
}

func TestToolCallMatcher_Array(t *testing.T) {
	m := NewToolCallMatcher()
	buf := []byte(`[{"name": "current_time", "arguments": {}}]`)

	if r := m.Scan(buf, ""); r != tape.FullMatch {
		t.Fatalf("expected FullMatch, got %v", r)
	}
	action := m.Extract(buf)
	tc, ok := action.(ToolCallAction)
	if !ok {
		t.Fatalf("expected ToolCallAction, got %T", action)
	}
	if len(tc.Calls) != 1 || tc.Calls[0].Name != "current_time" {
		t.Errorf("unexpected calls: %+v", tc.Calls)
	}
}

func TestToolCallMatcher_MarkdownFenced(t *testing.T) {
	m := NewToolCallMatcher()
	buf := []byte("```json\n{\"name\": \"fetch_page\", \"arguments\": {\"url\": \"https://example.com\", \"question\": \"test\"}}\n```")

	if r := m.Scan(buf, ""); r != tape.FullMatch {
		t.Fatalf("expected FullMatch, got %v", r)
	}
	action := m.Extract(buf)
	tc, ok := action.(ToolCallAction)
	if !ok {
		t.Fatalf("expected ToolCallAction, got %T", action)
	}
	if tc.Calls[0].Name != "fetch_page" {
		t.Errorf("expected 'fetch_page', got %q", tc.Calls[0].Name)
	}
}

func TestToolCallMatcher_MultipleToolsInArray(t *testing.T) {
	m := NewToolCallMatcher()
	buf := []byte(`[{"name": "web_search", "arguments": {"query": "go"}}, {"name": "current_time", "arguments": {}}]`)

	if r := m.Scan(buf, ""); r != tape.FullMatch {
		t.Fatalf("expected FullMatch, got %v", r)
	}

	action := m.Extract(buf)
	tc, ok := action.(ToolCallAction)
	if !ok {
		t.Fatalf("expected ToolCallAction, got %T", action)
	}
	if len(tc.Calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(tc.Calls))
	}
	if tc.Calls[0].Name != "web_search" {
		t.Errorf("first call = %q, want web_search", tc.Calls[0].Name)
	}
	if tc.Calls[1].Name != "current_time" {
		t.Errorf("second call = %q, want current_time", tc.Calls[1].Name)
	}
}

func TestToolCallMatcher_PartialJSON(t *testing.T) {
	m := NewToolCallMatcher()
	buf := []byte(`{"name": "web_search", "arguments": {"query": "go`)

	if r := m.Scan(buf, ""); r != tape.PartialMatch {
		t.Errorf("expected PartialMatch for incomplete JSON, got %v", r)
	}
}

func TestToolCallMatcher_Name(t *testing.T) {
	m := NewToolCallMatcher()
	if m.Name() != "tool_call" {
		t.Errorf("expected 'tool_call', got %q", m.Name())
	}
}

func TestToolCallMatcher_RegularText(t *testing.T) {
	m := NewToolCallMatcher()
	buf := []byte("This is just regular text about searching the web.")

	if r := m.Scan(buf, ""); r != tape.NoMatch {
		t.Errorf("expected NoMatch for regular text, got %v", r)
	}
}

func TestToolCallMatcher_BraceInProse(t *testing.T) {
	m := NewToolCallMatcher()

	buf := []byte(`{`)
	if r := m.Scan(buf, ""); r != tape.PartialMatch {
		t.Errorf("expected PartialMatch for lone brace, got %v", r)
	}

	buf = []byte(`{some random text that is definitely not JSON and is long enough}`)
	if r := m.Scan(buf, ""); r != tape.NoMatch {
		t.Errorf("expected NoMatch for non-JSON brace text, got %v", r)
	}
}

func TestToolCallMatcher_OversizedNonJSON(t *testing.T) {
	m := NewToolCallMatcher()
	buf := []byte(`{"name": "web_search", "arguments": {"query": "` + strings.Repeat("x", toolCallMaxAccumulate) + `"}}`)

	if r := m.Scan(buf, ""); r != tape.NoMatch {
		t.Errorf("expected NoMatch for oversized buffer, got %v", r)
	}
}

func TestToolCallMatcher_Empty(t *testing.T) {
	m := NewToolCallMatcher()
	if r := m.Scan([]byte(""), ""); r != tape.NoMatch {
		t.Errorf("expected NoMatch for empty, got %v", r)
	}
}

// --- Integration: StreamFilter + ToolCallMatcher ---

func TestStreamFilter_ToolCallDetection(t *testing.T) {
	c := &testCollector{}
	f := tape.NewStreamFilter(c.emit, []tape.Matcher{NewToolCallMatcher()}, 200)

	json := `{"name": "web_search", "arguments": {"query": "golang"}}`
	for _, ch := range json {
		action := f.Write(string(ch))
		if tc, ok := action.(ToolCallAction); ok {
			if tc.Calls[0].Arguments["query"] != "golang" {
				t.Errorf("expected query 'golang', got %v", tc.Calls[0].Arguments["query"])
			}
			if c.all() != "" {
				t.Errorf("expected no emission for tool call JSON, got %q", c.all())
			}
			return
		}
	}

	action := f.Flush()
	if tc, ok := action.(ToolCallAction); ok {
		if len(tc.Calls) != 1 || tc.Calls[0].Name != "web_search" {
			t.Errorf("unexpected tool call: %+v", tc.Calls)
		}
	} else {
		t.Fatalf("expected ToolCallAction, got %T", action)
	}
}

func TestStreamFilter_NormalTextThenToolCall(t *testing.T) {
	c := &testCollector{}
	f := tape.NewStreamFilter(c.emit, []tape.Matcher{NewToolCallMatcher()}, 20)

	f.Write("Let me search for that information now. ")

	if c.all() == "" {
		t.Fatal("expected normal text to be emitted before tool call")
	}

	toolJSON := `{"name": "web_search", "arguments": {"query": "test"}}`
	var gotToolCall bool
	for _, ch := range toolJSON {
		action := f.Write(string(ch))
		if _, ok := action.(ToolCallAction); ok {
			gotToolCall = true
			break
		}
	}
	if !gotToolCall {
		action := f.Flush()
		if _, ok := action.(ToolCallAction); ok {
			gotToolCall = true
		}
	}

	if !gotToolCall {
		t.Fatal("expected tool call to be detected")
	}

	emitted := c.all()
	if !strings.Contains(emitted, "Let me search") {
		t.Errorf("expected normal text in emitted output, got %q", emitted)
	}
	if strings.Contains(emitted, `"name"`) {
		t.Error("tool call JSON should not appear in emitted text")
	}
}

func TestStreamFilter_FalsePositiveBrace(t *testing.T) {
	c := &testCollector{}
	f := tape.NewStreamFilter(c.emit, []tape.Matcher{NewToolCallMatcher()}, 30)

	f.Write("I found {some results} in the data")
	f.Flush()

	emitted := c.all()
	if !strings.Contains(emitted, "{some results}") {
		t.Errorf("expected brace text to pass through, got %q", emitted)
	}
}

type testCollector struct {
	chunks []string
}

func (c *testCollector) emit(s string) {
	c.chunks = append(c.chunks, s)
}

func (c *testCollector) all() string {
	return strings.Join(c.chunks, "")
}
