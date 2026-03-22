# Axosoft CLI / MCP Server

## Overview
MCP server that exposes Axosoft work item management as tools for Claude Code. Also provides a CLI `login` command for OAuth authentication.

## Build & Test
```bash
make          # build binary
make test     # run tests
go mod tidy   # after dependency changes
```

## Architecture
- **Based on**: `github.com/deweysasser/golang-program` template (Kong CLI, Zerolog, Makefile)
- **MCP framework**: `github.com/mark3labs/mcp-go`
- **Two modes**: `login` (OAuth browser flow) and `serve` (MCP stdio server)

## Key Packages
- `program/` — CLI options, subcommands (login, serve)
- `auth/` — OAuth flow and token caching (`~/.config/axosoft-cli/token.json`)
- `axosoft/` — HTTP API client for Axosoft v5 API
- `tools/` — MCP tool definitions and handlers

## MCP Tools
| Tool | Description |
|------|-------------|
| `list_items` | Search/filter work items |
| `get_item` | Get full item details |
| `list_workflow_steps` | Available workflow transitions |
| `transition_item` | Move item to new state |
| `list_projects` | List all projects |
| `list_releases` | List releases |

## Design Decisions
- Logging goes to **stderr** (stdout is MCP stdio transport)
- Token loaded once at serve startup; re-run login if expired
- Item type (feature/task/incident/defect) required for API calls — Axosoft routes differ by type
- `BaseURL` and `ClientID` via env vars (`AXOSOFT_URL`, `AXOSOFT_CLIENT_ID`) or flags

## Configuration
```json
{
  "axosoft": {
    "type": "stdio",
    "command": "/path/to/axosoft-cli",
    "args": ["serve"],
    "env": {
      "AXOSOFT_URL": "https://example.axosoft.com",
      "AXOSOFT_CLIENT_ID": "your-client-id"
    }
  }
}
```
