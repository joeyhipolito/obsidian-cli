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
- User wants to sync website content into the vault
- User wants to find connections or suggested links between notes
- User wants a vault health check or to fix broken links/missing frontmatter
- User wants to import scout intel or orchestrator learnings into the vault

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

### Sync

Sync website MDX content metadata into the Obsidian vault as note stubs under `20 Projects/Website/`. Compares modification times and only updates changed files.

```bash
# Sync website content into vault (incremental)
obsidian sync

# Preview what would change without writing
obsidian sync --dry-run

# Force overwrite unchanged notes and include unpublished content
obsidian sync --force

# JSON output with created/updated/unchanged/skipped lists
obsidian sync --json
```

Requires `website_path` in `~/.obsidian/config` or `OBSIDIAN_WEBSITE_PATH` env var.

### Enrich

Analyze the vault index to suggest wikilinks between semantically similar notes, recommend tags via consensus from neighbors, and detect orphan notes with no incoming links.

```bash
# Show link suggestions, tag suggestions, and orphan notes
obsidian enrich

# Apply suggested wikilinks to notes (appends to Related Notes section)
obsidian enrich --apply

# JSON output with similarity scores
obsidian enrich --json
```

Requires a built index (`obsidian index` first).

### Maintain

Run vault health checks: stale notes, broken wikilinks, empty notes, large notes (>10KB), missing frontmatter, and index coverage. Outputs a 0-100 health score.

```bash
# Full vault health report
obsidian maintain

# Custom staleness threshold (default: 30 days)
obsidian maintain --stale-days 60

# Auto-fix: add empty frontmatter to notes missing it
obsidian maintain --fix

# JSON output with all health data
obsidian maintain --json
```

### Ingest

Import data from external sources into the vault as structured notes. Deduplicates using `~/.obsidian/ingest-state.json`.

**Scout intel** (from `~/.scout/intel/{topic}/{date}_{source}.json`):
- Creates notes at `vault/Intel/{topic}/{slug}.md`
- Frontmatter includes: type, source, topic, url, date, score, tags

**Learnings** (from `~/.via/learnings.db`):
- Creates notes at `vault/Learnings/{domain}/{type}-{id}.md`
- Frontmatter includes: domain, learning-type, agent-type, created, seen-count, used-count

```bash
obsidian ingest --source scout                          # Ingest all scout intel
obsidian ingest --source scout --topic "ai-models"     # Filter by topic
obsidian ingest --source scout --since 7d              # Last 7 days only
obsidian ingest --source learnings                     # Ingest all learnings
obsidian ingest --source learnings --domain dev        # Filter by domain
obsidian ingest --source learnings --since 30d         # Last 30 days only
obsidian ingest --source scout --dry-run               # Preview without writing
obsidian ingest --source scout --json                  # JSON output
```

Supported `--since` units: `h` (hours), `d` (days), `w` (weeks). Example: `7d`, `24h`, `2w`.

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

**User**: "sync my website content to the vault"
**Action**: `obsidian sync` or `obsidian sync --dry-run` to preview first

**User**: "find connections between my notes"
**Action**: `obsidian enrich` then `obsidian enrich --apply` to write the links

**User**: "how healthy is my vault"
**Action**: `obsidian maintain`

**User**: "import my scout intel into my vault"
**Action**: `obsidian ingest --source scout --dry-run` then `obsidian ingest --source scout`

**User**: "add my recent learnings to Obsidian"
**Action**: `obsidian ingest --source learnings --since 7d`

**User**: "ingest the last week of ai-models scout intel"
**Action**: `obsidian ingest --source scout --topic "ai-models" --since 7d`

All commands support `--json` for machine-readable output.
