# axosoft-cli

An [MCP](https://modelcontextprotocol.io/) server and CLI tool for [Axosoft](https://www.axosoft.com/) work item management. Exposes Axosoft's v5 API as MCP tools for use with Claude Code, and provides direct CLI commands for testing and scripting.

## Installation

```bash
go install github.com/deweysasser/axosoft-cli@latest
```

Or build from source:

```bash
git clone https://github.com/deweysasser/axosoft-cli.git
cd axosoft-cli
make
```

## Authentication

Before using any commands, you must authenticate with your Axosoft instance. Tokens are cached to `~/.config/axosoft-cli/token.json`.

### Browser OAuth Flow

Opens a browser for interactive login. Requires an OAuth client ID (create one in Axosoft under Tools > System Options > Axosoft API Settings).

```bash
export AXOSOFT_URL=https://example.axosoft.com
export AXOSOFT_CLIENT_ID=your-client-id

axosoft-cli login
```

### Password Grant Flow

Non-interactive login using username and password. Requires an OAuth client ID and secret.

```bash
axosoft-cli login \
  --base-url=https://example.axosoft.com \
  --client-id=your-client-id \
  --client-secret=your-client-secret \
  --username=your-username \
  --password=your-password
```

All flags can also be set via environment variables:

| Flag | Environment Variable |
|------|---------------------|
| `--base-url` | `AXOSOFT_URL` |
| `--client-id` | `AXOSOFT_CLIENT_ID` |
| `--client-secret` | `AXOSOFT_CLIENT_SECRET` |
| `--username` | `AXOSOFT_USERNAME` |
| `--password` | `AXOSOFT_PASSWORD` |

## CLI Commands

All read commands output indented JSON to stdout. Logging goes to stderr.

### List Projects

```bash
axosoft-cli projects
```

### List Releases

```bash
axosoft-cli releases
axosoft-cli releases --project-id=3
```

### List Work Items

```bash
axosoft-cli items
axosoft-cli items --search="login bug"
axosoft-cli items --project-id=3 --type=defect
axosoft-cli items --release-id=10 --type=task
```

| Flag | Description | Default |
|------|-------------|---------|
| `--project-id` | Filter by project ID | (none) |
| `--release-id` | Filter by release ID | (none) |
| `--search` | Search string | (none) |
| `--type` | Item type: `feature`, `task`, `incident`, `defect` | `feature` |

### Get a Single Work Item

```bash
axosoft-cli item 12345
axosoft-cli item 12345 --type=defect
```

## MCP Server

Start the MCP server for use with Claude Code or other MCP clients:

```bash
axosoft-cli serve
```

The server communicates over stdio using the MCP protocol.

### Available MCP Tools

| Tool | Description |
|------|-------------|
| `list_projects` | List all projects |
| `list_releases` | List releases, optionally filtered by project |
| `list_items` | Search and filter work items |
| `get_item` | Get full details of a single work item |
| `list_workflow_steps` | List available workflow transitions for an item |
| `transition_item` | Move an item to a new workflow state |

### Claude Code Configuration

Add to your Claude Code MCP settings (`.claude/settings.json` or `~/.claude/settings.json`):

```json
{
  "mcpServers": {
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
}
```

You must run `axosoft-cli login` before starting the MCP server to cache a token.

## Global Flags

| Flag | Description |
|------|-------------|
| `--base-url` / `AXOSOFT_URL` | Axosoft instance URL (required) |
| `--client-id` / `AXOSOFT_CLIENT_ID` | OAuth client ID |
| `--debug` | Show debug logging |
| `--quiet` | Suppress info logging |
| `--output-format` | Log format: `auto`, `terminal`, `jsonl` (default: `auto`) |
| `--version` | Show program version |

## Axosoft Item Types

Axosoft uses separate API endpoints for each item type. The `--type` flag (default: `feature`) controls which endpoint is queried:

- `feature` — Feature requests
- `task` — Tasks
- `incident` — Incidents
- `defect` — Defects

## Development

```bash
make          # build binary
make test     # run tests
go mod tidy   # after dependency changes
```
