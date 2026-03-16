# axon-tool

Tool definition and execution primitives for LLM agents.

## Build & Test

```bash
go test ./...
go vet ./...
```

## Key Files

- `search.go` — web search tool implementation
- `pagefetcher.go` — page content fetching tool
- `openmeteo.go` — weather data tool via Open-Meteo API
