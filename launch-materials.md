# Launch Materials for notion-cli v0.2.0

## One-liner Pitch
Like `gh` for GitHub, but for Notion. 39 commands. One binary.

## Reddit r/Notion Post

**Title:** I built a full CLI for Notion — 39 commands, human-friendly filters, Markdown I/O, one binary

**Body:**
Hey r/Notion! I built **notion-cli** — a command-line tool that gives you complete access to Notion from your terminal.

**Why?** I'm an AI agent developer and needed programmatic access to Notion without the browser. The existing CLIs only cover a fraction of the API, so I built one that covers 100%.

**What it does:**
- 39 commands covering pages, databases, blocks, comments, users, and files
- Human-friendly filters: `--filter 'Status=Done'` instead of raw JSON
- Schema-aware: auto-detects property types when creating pages
- Markdown in/out: read blocks as Markdown, write Markdown to Notion
- Pipe-friendly: colored tables in terminal, clean JSON when piped
- Single binary, zero dependencies

**Install:**
```
brew install 4ier/tap/notion-cli    # macOS/Linux
scoop bucket add 4ier https://github.com/4ier/scoop-bucket && scoop install notion-cli  # Windows
npm install -g @4ier/notion-cli     # npm
go install github.com/4ier/notion-cli@latest  # Go
```

**Quick example:**
```
notion db query <db-id> --filter 'Status=Done' --sort 'Date:desc'
notion page create <db-id> --db "Name=Weekly Review" "Status=Todo"
notion block list <page-id> --md --depth 3
```

GitHub: https://github.com/4ier/notion-cli

Feedback welcome — what commands would you use most?

---

## Reddit r/commandline Post (Day+1)

**Title:** notion-cli: Full Notion API from your terminal — 39 commands, filters, Markdown I/O

**Body:**
Built a CLI that wraps 100% of the Notion API into 39 subcommands. Think `gh` for GitHub, but for Notion.

Highlights for the CLI crowd:
- Human-friendly filter syntax: `--filter 'Status=Done' --sort 'Date:desc'`
- JSON escape hatch for complex queries: `--filter-json '{...}'`
- Smart output: colored tables in TTY, clean JSON when piped
- Markdown I/O: `--md` flag, `--file notes.md`
- Shell completion (bash/zsh/fish/powershell)
- Single static binary (Go, CGO_ENABLED=0)

Particularly useful for scripting and AI agents that need to interact with Notion workspaces.

`brew install 4ier/tap/notion-cli`

GitHub: https://github.com/4ier/notion-cli

---

## Hacker News Post

**Title:** Show HN: Notion-CLI – Full Notion API from the terminal, 39 commands, one binary

**URL:** https://github.com/4ier/notion-cli

---

## X/Twitter Post

Just shipped notion-cli v0.2.0 🚀

Like `gh` for GitHub, but for @NotionHQ.
39 commands. One binary. Zero dependencies.

Human-friendly filters, Markdown I/O, schema-aware page creation, pipe-friendly JSON output.

Built for developers and AI agents.

github.com/4ier/notion-cli

---

## Product Hunt (Day+3)

**Tagline:** Like gh for GitHub, but for Notion. 39 commands. One binary.

**Description:**
notion-cli is a full-featured command-line interface for Notion. It covers 100% of the Notion API with 39 subcommands — manage pages, databases, blocks, comments, users, and files without leaving your terminal.

Key features:
• Human-friendly filter syntax (no raw JSON needed)
• Schema-aware page creation
• Markdown input/output
• Smart output (colored tables in terminal, JSON when piped)
• Single binary, available via Homebrew, Scoop, npm, Go, Docker, deb/rpm

Built for developers who live in the terminal and AI agents that need programmatic Notion access.
