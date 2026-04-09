// Package tool provides primitives for defining and executing tools
// that can be used by LLM-powered agents. It is provider-agnostic,
// with no dependency on any specific LLM backend.
//
// Core types:
//   - [ToolDef] describes a tool's name, description, parameters, and execute function
//   - [ToolResult] is the output of a tool execution
//   - [ToolContext] carries request-scoped state (user, agent, conversation)
//   - [ToolCall] is a parsed tool invocation from an LLM response
//
// The [NormalizeArguments] function coerces LLM-provided arguments
// to match JSON Schema types, handling common mismatches like numbers
// sent as strings.
package tool
