---
name: obsidian
description: Reads, writes, searches, and indexes Obsidian vault notes via obsidian CLI. Use when user asks about notes, vault, daily notes, or wants to find/create/update markdown files.
allowed-tools: Bash(obsidian:*)
---

# Obsidian - Vault Notes

Read, write, search, and index Obsidian vault notes via the standalone `obsidian` CLI tool.

## When to Use

- User mentions notes, vault, daily notes, or Obsidian
- User wants to find a note by content or topic
- User wants to read, create, or append to a note
- User asks "what notes do I have about X"
- User wants to update their daily note

## Commands

### Read & List

```bash
obsidian read "path/to/note"              # Read note with parsed frontmatter
obsidian read "path/to/note" --json       # JSON with frontmatter, headings, wikilinks
obsidian list                             # List all notes in vault
obsidian list "subfolder"                 # List notes in subdirectory
obsidian list --json                      # JSON output with metadata
```

### Write & Create

```bash
# Create a new note
obsidian create "path/to/note" --title "My Note"

# Append text to existing note
obsidian append "path/to/note" "Text to append"

# Append under a specific section
obsidian append "path/to/note" --section "## Capture" "New entry"

# Pipe content from stdin
echo "piped content" | obsidian append "path/to/note"
```

### Search

```bash
# Hybrid search (default — combines keyword + semantic)
obsidian search "query terms"

# Keyword only (FTS5, exact matches)
obsidian search "query" --mode keyword

# Semantic only (vector similarity, conceptual matches)
obsidian search "query" --mode semantic

# JSON output with scores
obsidian search "query" --json
```

### Index

```bash
# Build/update search index (incremental — only re-indexes changed files)
obsidian index

# Full re-index
obsidian index --force
```

### Setup

```bash
obsidian configure            # Interactive setup (vault path, API key)
obsidian configure show       # Show current config
obsidian doctor               # Health checks (config, vault, index, API)
```

## Configuration

Config file: `~/.obsidian/config`

```ini
vault_path=/path/to/your/vault
gemini_apikey=your-gemini-api-key
```

Or use environment variables: `OBSIDIAN_VAULT_PATH`, `GEMINI_API_KEY`

## Examples

**User**: "what notes do I have about routines"
**Action**: `obsidian search "routines"` or `obsidian search "household routines" --mode semantic`

**User**: "add to my daily note: finished thesis chapter 3"
**Action**: `obsidian append "40 Time/Daily/2026/02/2026-02-07" "- Finished thesis chapter 3"`

**User**: "create a new meeting note"
**Action**: `obsidian create "10 Work/Meetings/2026-02-07 — Team Sync" --title "Team Sync"`

**User**: "show me my hub files"
**Action**: `obsidian list "10 Active/Hubs"`

**User**: "read my personal context"
**Action**: `obsidian read "90 System/Organization Files/Personal Context — Joey Hipolito"`

All commands support `--json` for machine-readable output.
