---
name: notion-cli
description: |
  Work with Notion from the terminal using the `notion-agent` CLI. Use when the user needs to read, create, update, query, or manage Notion pages, databases, blocks, comments, users, or files programmatically. Covers the entire Notion API with 47 commands. Triggers: Notion workspace automation, database queries, page creation, block manipulation, comment threads, file uploads, relation management, database export, multi-workspace management, or any Notion API interaction from the command line.
---

# Notion Agent CLI

`notion-agent` is a CLI for the Notion API. Single Go binary, full API coverage, dual output (pretty tables for humans, JSON for agents).

## Install

```bash
# Homebrew
brew install MaxMa04/tap/notion-agent-cli

# Go
go install github.com/MaxMa04/notion-agent-cli@latest

# npm
npm install -g @vibelabsio/notion-agent-cli

# Or download binary from GitHub Releases
# https://github.com/MaxMa04/notion-agent-cli/releases
```

## Auth

```bash
notion-agent auth login --with-token <<< "ntn_xxxxxxxxxxxxx"   # or interactive
notion-agent auth login --with-token --profile work <<< "ntn_xxx"  # save as named profile
export NOTION_TOKEN=ntn_xxxxxxxxxxxxx                     # env var alternative
notion-agent auth status                                        # show current profile
notion-agent auth switch                                        # interactive profile picker
notion-agent auth switch work                                   # switch to named profile
notion-agent auth doctor                                        # health check
```

## Command Reference

### Search
```bash
notion-agent search "query"                    # search everything
notion-agent search "query" --type page        # pages only
notion-agent search "query" --type database    # databases only
```

### Pages
```bash
notion-agent page view <id|url>                # render page content
notion-agent page list                         # list workspace pages
notion-agent page create <parent> --title "X" --body "content"
notion-agent page create <db-id> --db "Name=Review" "Status=Todo"  # database row
notion-agent page delete <id>                  # archive page
notion-agent page restore <id>                 # unarchive page
notion-agent page move <id> --to <parent>
notion-agent page open <id>                    # open in browser
notion-agent page edit <id|url>                # edit in $EDITOR (Markdown round-trip)
notion-agent page edit <id> --editor nano      # specify editor
notion-agent page set <id> Key=Value ...       # set properties (type-aware)
notion-agent page props <id>                   # show all properties
notion-agent page props <id> <prop-id>         # get specific property
notion-agent page link <id> --prop "Rel" --to <target-id>    # add relation
notion-agent page unlink <id> --prop "Rel" --from <target-id> # remove relation
```

### Databases
```bash
notion-agent db list                           # list databases
notion-agent db view <id>                      # show schema
notion-agent db query <id>                     # query all rows
notion-agent db query <id> -F 'Status=Done' -s 'Date:desc'  # filter + sort
notion-agent db query <id> --filter-json '{"or":[...]}'     # complex JSON filter
notion-agent db query <id> --all               # fetch all pages
notion-agent db create <parent> --title "X" --props "Status:select,Date:date"
notion-agent db update <id> --title "New Name" --add-prop "Priority:select"
notion-agent db add <id> "Name=Task" "Status=Todo" "Priority=High"
notion-agent db add-bulk <id> --file items.json              # bulk create from JSON
notion-agent db export <id>                    # export all rows as CSV (default)
notion-agent db export <id> --format json      # export as JSON
notion-agent db export <id> --format md -o report.md  # export as Markdown table to file
notion-agent db open <id>                      # open in browser
```

#### Filter operators
| Syntax | Meaning |
|--------|---------|
| `=` | equals |
| `!=` | not equals |
| `>` / `>=` | greater than (or equal) |
| `<` / `<=` | less than (or equal) |
| `~=` | contains |

Multiple `-F` flags combine with AND. Property types are auto-detected from schema.

#### Sort: `-s 'Date:desc'` or `-s 'Name:asc'`

#### Bulk add file format
```json
[{"Name": "Task A", "Status": "Todo"}, {"Name": "Task B", "Status": "Done"}]
```

### Blocks
```bash
notion-agent block list <parent-id>            # list child blocks
notion-agent block list <parent-id> --all      # paginate through all
notion-agent block list <parent-id> --depth 3  # recursive nested blocks
notion-agent block list <parent-id> --md       # output as Markdown
notion-agent block get <id>                    # get single block
notion-agent block append <parent> "text"      # append paragraph
notion-agent block append <parent> "text" -t bullet          # bullet point
notion-agent block append <parent> "text" -t code --lang go  # code block
notion-agent block append <parent> --file notes.md           # from file
notion-agent block insert <parent> "text" --after <block-id> # positional insert
notion-agent block update <id> --text "new"    # update content
notion-agent block delete <id1> [id2] [id3]    # delete one or more
notion-agent block move <id> --after <target>  # reposition after target block
notion-agent block move <id> --before <target> # reposition before target block
notion-agent block move <id> --parent <new-parent>  # move to different parent
```

Block types: `paragraph`/`p`, `h1`, `h2`, `h3`, `bullet`, `numbered`, `todo`, `quote`, `code`, `callout`, `divider`

### Comments
```bash
notion-agent comment list <page-id>
notion-agent comment add <page-id> "comment text"
notion-agent comment get <comment-id>
notion-agent comment reply <comment-id> "reply text"  # reply in same thread
```

### Users
```bash
notion-agent user me                           # current bot info
notion-agent user list                         # all workspace users
notion-agent user get <user-id>
```

### Files
```bash
notion-agent file list                         # list uploads
notion-agent file upload ./path/to/file        # upload (auto MIME detection)
```

### Raw API (escape hatch)
```bash
notion-agent api GET /v1/users/me
notion-agent api POST /v1/search '{"query":"test"}'
notion-agent api PATCH /v1/pages/<id> '{"archived":true}'
```

## Output Modes

- **Terminal (TTY)**: colored tables, readable formatting
- **Piped/scripted**: JSON automatically
- **Explicit**: `--format json` or `--format table`
- `--debug`: show HTTP request/response details

All output includes full Notion UUIDs. All commands accept Notion URLs or IDs.

## Tips

- `notion-agent db add` and `notion-agent page set` auto-detect property types from schema
- Multi-select: `Tags=tag1,tag2,tag3`
- Checkbox: `Done=true`
- Pipe to jq: `notion-agent db query <id> -F 'Status=Done' --format json | jq '.results[].id'`
