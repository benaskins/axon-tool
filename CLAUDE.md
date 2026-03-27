@AGENTS.md

## Conventions
- Every `ToolDef` must have a non-empty `Description` field
- Parameters use JSON Schema encoding — define with `Parameters` map on `ToolDef`
- Built-in tool constructors live in `tools.go`; add new ones there
- `ToolContext` carries execution context; `ToolResult` is the return type

## Constraints
- Zero dependencies on axon or any other axon-* module — this is a leaf library
- Keep the interface surface minimal; tool implementors depend on these types
- Do not add HTTP, networking, or server concerns — those belong in axon
- External dependency (`go-readability`) is for page fetching only; avoid adding more

## Testing
- `go test ./...` — all tests run without network access
- `go vet ./...` for lint
