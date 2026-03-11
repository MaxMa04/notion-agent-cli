<h1 align="center">
  notion-agent
</h1>

<p align="center">
  <b>Like <code>gh</code> for GitHub, but for Notion. 47 commands. One binary.</b>
</p>

<p align="center">
  <img src="https://raw.githubusercontent.com/MaxMa04/notion-agent-cli/main/demo.gif" alt="demo" width="640">
</p>

<p align="center">
  <a href="https://github.com/MaxMa04/notion-agent-cli/releases"><img src="https://img.shields.io/github/v/release/MaxMa04/notion-agent-cli?style=flat-square" alt="Release"></a>
  <a href="https://github.com/MaxMa04/notion-agent-cli/actions"><img src="https://img.shields.io/github/actions/workflow/status/MaxMa04/notion-agent-cli/test.yml?style=flat-square&label=tests" alt="Tests"></a>
  <a href="https://github.com/MaxMa04/notion-agent-cli/blob/main/LICENSE"><img src="https://img.shields.io/github/license/MaxMa04/notion-agent-cli?style=flat-square" alt="License"></a>
  <a href="https://goreportcard.com/report/github.com/MaxMa04/notion-agent-cli"><img src="https://goreportcard.com/badge/github.com/MaxMa04/notion-agent-cli?style=flat-square" alt="Go Report Card"></a>
</p>

---

A full-featured command-line interface for [Notion](https://notion.so). Manage pages, databases, blocks, comments, users, and files — all from your terminal. Built for developers and AI agents who need programmatic access without the browser.

## Install

### npm (cross-platform)
```sh
npm install -g notion-agent-cli
```

### Go
```sh
go install github.com/MaxMa04/notion-agent-cli@latest
```

### Shell Script (Linux / macOS)
```sh
curl -fsSL https://raw.githubusercontent.com/MaxMa04/notion-agent-cli/main/install.sh | sh
```

### Binary

Download from [GitHub Releases](https://github.com/MaxMa04/notion-agent-cli/releases) — available for Linux, macOS, and Windows (amd64/arm64).

## Quick Start

```sh
# Authenticate
notion-agent auth login
# Or pipe token directly
echo "ntn_xxxxx" | notion-agent auth login --with-token

# Search your workspace
notion-agent search "meeting notes"

# Query a database with filters
notion-agent db query <db-id> --filter 'Status=Done' --sort 'Date:desc'

# Create a page in a database
notion-agent page create <db-id> --db "Name=Weekly Review" "Status=Todo"

# Apply a template to a new page
notion-agent page apply-template <page-id> <template-page-id>

# Read page content as Markdown
notion-agent block list <page-id> --depth 3 --md

# Append blocks from a Markdown file
notion-agent block append <page-id> --file notes.md

# Create a table
notion-agent block table <page-id> "Name,Status,Date" "Task 1,Done,2026-03-10"

# Raw API escape hatch
notion-agent api GET /v1/users/me
```

## Commands

| Group | Commands | Description |
|-------|----------|-------------|
| **auth** | `login` `logout` `status` `doctor` `switch` | Authentication, diagnostics & profile management |
| **search** | `search` | Search pages and databases |
| **page** | `view` `list` `create` `delete` `restore` `move` `open` `set` `props` `edit` `link` `unlink` `apply-template` | Full page lifecycle + templates |
| **db** | `list` `view` `query` `create` `update` `add` `add-bulk` `export` `open` | Database CRUD, query & export |
| **block** | `list` `get` `append` `insert` `update` `delete` `move` `table` `table-add` | Content blocks + native tables |
| **comment** | `list` `add` `get` `reply` | Discussion threads |
| **user** | `me` `list` `get` | Workspace members |
| **file** | `list` `upload` | File management |
| **api** | `<METHOD> <path>` | Raw API escape hatch |

**47 subcommands** covering 100% of the Notion API.

## Features

### Multi-Profile Support

Manage multiple Notion workspaces with named profiles:

```sh
# Login with a named profile
notion-agent auth login --profile work
notion-agent auth login --profile personal

# Switch active profile
notion-agent auth switch work

# Use a profile for a single command
notion-agent search "project" --profile work

# Or set via environment variable
export NOTION_PROFILE=work
```

Priority: `NOTION_TOKEN` env > `--profile` flag > `NOTION_PROFILE` env > current profile in config.

### Human-Friendly Filters
No JSON needed for 90% of queries:
```sh
notion-agent db query <id> --filter 'Status=Done' --filter 'Priority=High' --sort 'Date:desc'
```

For complex queries (OR, nesting), use the JSON escape hatch:
```sh
notion-agent db query <id> --filter-json '{"or":[{"property":"Status","status":{"equals":"Done"}},{"property":"Status","status":{"equals":"Cancelled"}}]}'
```

### Schema-Aware Properties
Property types are auto-detected from the database schema:
```sh
notion-agent page create <db-id> --db "Name=Sprint Review" "Date=2026-03-01" "Points=8" "Done=true"
```

### Template Support

Copy the content structure from any page to another — useful for applying task templates, meeting note formats, etc.:

```sh
# Create a task and apply a template
notion-agent db add <db-id> "Name=New Task" "Status=Backlog"
notion-agent page apply-template <new-page-id> <template-page-id>
```

Recursively copies all blocks including nested children (tables, toggle content, etc.).

### Native Tables

Create and extend tables directly:

```sh
# Create a table (first arg = headers, rest = rows)
notion-agent block table <page-id> "Name,Status,Due" "Task 1,Done,2026-03-10" "Task 2,Open,2026-03-15"

# Add rows to an existing table
notion-agent block table-add <table-block-id> "Task 3,In Progress,2026-03-20"
```

### Smart Output
- **Terminal**: Colored tables, formatted text
- **Pipe/Script**: Clean JSON for `jq`, scripts, and AI agents
```sh
# Pretty table in terminal
notion-agent db query <id>

# JSON when piped
notion-agent db query <id> | jq '.results[].properties.Name'
```

### Markdown I/O
```sh
# Read blocks as Markdown
notion-agent block list <page-id> --md --depth 3

# Write Markdown to Notion
notion-agent block append <page-id> --file document.md
```
Supports headings, bullets, numbered lists, todos, quotes, code blocks, tables, and dividers.

### Page Editing

Open a page in your text editor, edit as Markdown, save to sync back:
```sh
notion-agent page edit <page-id>
notion-agent page edit <page-id> --editor code
```

### Database Export
```sh
notion-agent db export <db-id> --format csv > tasks.csv
notion-agent db export <db-id> --format json > tasks.json
notion-agent db export <db-id> --format md > tasks.md
```

### Recursive Block Reading
```sh
notion-agent block list <page-id> --depth 5 --all
```

### URL or ID — Your Choice
```sh
# Both work
notion-agent page view abc123def
notion-agent page view https://notion.so/My-Page-abc123def456
```

### Actionable Error Messages
```
object_not_found: Could not find page with ID abc123
  -> Check the ID is correct and the page/database is shared with your integration
```

## For AI Agents

This CLI is designed to be agent-friendly:
- **JSON output** when piped — no parsing needed
- **Schema-aware** — agents don't need to know property types
- **URL resolution** — paste Notion URLs directly
- **Template support** — apply consistent page structures automatically
- **Multi-profile** — agents can target specific workspaces via `--profile` or `NOTION_PROFILE`
- **Single binary** — no runtime dependencies
- **Exit codes** — 0 for success, non-zero for errors

## Configuration

```sh
# Interactive login
notion-agent auth login

# Pipe token (non-interactive, CI-friendly)
echo "ntn_xxxxx" | notion-agent auth login --with-token

# Or use environment variable (no config needed)
export NOTION_TOKEN=ntn_xxxxx

# Config stored in ~/.config/notion-agent-cli/config.json (mode 0600)
# Legacy path ~/.config/notion-cli/ is auto-detected

# Check authentication
notion-agent auth status
notion-agent auth doctor

# Manage profiles
notion-agent auth login --profile work
notion-agent auth switch work
notion-agent auth logout work
```

## Acknowledgements

This project is based on [notion-cli](https://github.com/4ier/notion-cli) by [4ier](https://github.com/4ier).

## Contributing

Issues and PRs welcome at [github.com/MaxMa04/notion-agent-cli](https://github.com/MaxMa04/notion-agent-cli).

## License

[MIT](LICENSE)
