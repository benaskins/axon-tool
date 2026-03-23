# axon-tool

Tool definition and execution primitives for LLM agents.

## Build & Test

```bash
go test ./...
go vet ./...
```

## Dependencies

- `github.com/go-shiori/go-readability` — HTML content extraction for page fetching

## Key Files

- `tool.go` — core types: ToolDef, ToolResult, ToolContext, schemas
- `tools.go` — built-in tool constructors (time, search, weather, page fetch)
- `search.go` — web search abstractions and utilities
- `searxng.go` — SearXNG client implementation
- `searchqualifier.go` — LLM-powered search query refinement
- `pagefetcher.go` — page content fetching and extraction
- `weather.go` — weather data abstractions
- `openmeteo.go` — Open-Meteo weather API client
- `doc.go` — package documentation