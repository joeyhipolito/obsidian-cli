---
name: obsidian
description: Reads, writes, searches, and indexes Obsidian vault notes via obsidian CLI. Use when user asks about notes, vault, daily notes, or wants to find/create/update markdown files.
allowed-tools: Bash(obsidian:*)
argument-hint: "[note-path or query]"
keywords: obsidian, vault, notes, markdown, search, embeddings, capture, triage, inbox
category: integration
version: "1.0.0"
---

# Obsidian - Vault Notes

Read, write, search, and index Obsidian vault notes via the standalone `obsidian` CLI tool.

## When to Use

- User mentions notes, vault, daily notes, or Obsidian
- User wants to capture a quick thought, link, or idea into the inbox
- User wants to process or triage inbox notes into the correct vault folder
- User wants to find a note by content or topic
- User wants to read, create, or append to a note
- User asks "what notes do I have about X"
- User wants to update their daily note
- User wants to check vault health or diagnose setup issues

## Capture → Triage → Connect Workflow

The standard workflow for creating notes:

1. **Capture** — dump thought into `Inbox/` instantly; no folder decision needed
2. **Triage** — classify and move inbox notes to the correct folder automatically
3. **Connect** — search and link related notes

## Commands

### Capture (Inbox)

Creates a `type: fleeting` note in `Inbox/` with a timestamp filename. No folder decision required at capture time.

```bash
# Capture a quick idea or thought
obsidian capture "rough idea about search indexing"

# Capture with a source URL or origin label
obsidian capture "link worth reading" --source https://example.com
obsidian capture "from the standup meeting" --source mission

# Capture from stdin (piped content)
echo "piped text" | obsidian capture

# JSON output — returns the created path
obsidian capture "my idea" --json
```

**Output frontmatter set by capture:**
```yaml
type: fleeting
created: 2026-02-21
source: <value of --source, if provided>
```

### Triage (Process Inbox)

Review and process notes in `Inbox/`. Default mode (no flags) lists pending notes.

```bash
# List pending inbox notes with age (default)
obsidian triage

# Filter to notes older than a duration (7d, 24h, 2w)
obsidian triage --older 7d

# Classify, enrich with wikilinks, and move each inbox note
obsidian triage --auto

# Preview what --auto would do without writing
obsidian triage --auto --dry-run

# Structured JSON output
obsidian triage --auto --json

# Cron-friendly: no output when inbox is clear; errors still surface
obsidian triage --auto --quiet
```

**Auto-triage type → destination mapping:**

| Classified type | Destination folder |
|---|---|
| `idea` | `Ideas/` |
| `task` | `Tasks/` |
| `reference` | `References/` |
| `note` | `Notes/` |

**Classification rules (in priority order):**
1. Existing non-`fleeting` type in frontmatter → preserved as-is
2. Checkbox items or `TODO:` / `action:` / `followup` in body → `task`
3. `source:` URL or inline `https://` URL → `reference`
4. Two or more headings → `note`
5. Default → `idea`

**Frontmatter set by `--auto`:**
```yaml
type: <classified-type>
status: processed
triaged: 2026-02-21
```

**Cron setup** (hourly, emails only on activity):
```
0 * * * * /usr/local/bin/obsidian triage --auto --quiet 2>&1
```

### Create

Create a new note with rich frontmatter. All flags are optional.

```bash
# Minimal: create a note (empty body)
obsidian create "path/to/note"

# With title (also adds an H1 heading to the body)
obsidian create "path/to/note" --title "My Note"

# Full frontmatter flags
obsidian create "Projects/my-idea" \
  --title "My Idea" \
  --type idea \
  --status draft \
  --summary "One-line description" \
  --tags "go,cli,search" \
  --context-set personal

# Clone body from an existing vault note as a template
obsidian create "path/to/note" --template "System/Templates/idea.md"
```

**Available `--create` flags:**

| Flag | Description |
|---|---|
| `--title <title>` | Sets `title:` in frontmatter and adds `# Title` H1 heading |
| `--type <type>` | Sets `type:` in frontmatter (e.g. `idea`, `note`, `task`, `reference`) |
| `--status <status>` | Sets `status:` in frontmatter (e.g. `draft`, `active`, `archived`) |
| `--summary <text>` | Sets `summary:` in frontmatter |
| `--tags <t1,t2>` | Sets `tags:` list in frontmatter (comma-separated) |
| `--context-set <name>` | Sets `context-set:` in frontmatter |
| `--template <path>` | Uses a vault note's body as the template for the new note |

### Append

Append text to an existing note, optionally inside a named section.

```bash
# Append to end of note
obsidian append "path/to/note" "Text to append"

# Append inside a specific section (inserts before the next heading)
obsidian append "path/to/note" --section "## Tasks" "- buy milk"
obsidian append "Daily/2026/02/2026-02-21" --section "## Capture" "New entry"

# Append from stdin
echo "piped content" | obsidian append "path/to/note"
```

### Read & List

```bash
# Read note with parsed body (frontmatter stripped, shown as header)
obsidian read "path/to/note"

# JSON with full frontmatter, headings, wikilinks, body
obsidian read "path/to/note" --json

# List all notes in vault
obsidian list

# List notes in a subdirectory
obsidian list "Inbox"
obsidian list "Projects"

# JSON output with path metadata
obsidian list --json
obsidian list "Projects" --json
```

### Search

```bash
# Hybrid search — combines keyword (FTS5) + semantic (vector) results (default)
obsidian search "query terms"

# Keyword-only search (FTS5, exact matches, fast)
obsidian search "query" --mode keyword

# Semantic-only search (vector similarity, conceptual matches)
obsidian search "query" --mode semantic

# JSON output with relevance scores
obsidian search "query" --json
```

### Resurface

Surface old notes that match a query (hybrid search) or are randomly selected. Useful for injecting relevant past knowledge into mission context.

```bash
# Find old notes relevant to a query (default: older than 7 days, top 5)
obsidian resurface "query terms"

# Narrow by age and result count
obsidian resurface "golang patterns" --older 14d --limit 3

# Random old note for serendipitous rediscovery
obsidian resurface --random

# Random with custom age threshold
obsidian resurface --random --older 30d --limit 10

# JSON output — suitable for injection into orchestrator mission context
obsidian resurface "query" --json
obsidian resurface --random --json
```

**Flags:**

| Flag | Default | Description |
|---|---|---|
| `--limit N` | 5 | Maximum number of results to return |
| `--older <duration>` | 7d | Only include notes older than this (e.g. `7d`, `14d`, `2w`, `24h`) |
| `--random` | — | Return randomly selected old notes instead of query-based results |
| `--json` | — | Machine-readable output with path, title, snippet, score, mod_time, age_days |

**JSON output format:**
```json
{
  "query": "golang patterns",
  "mode": "query",
  "older_than": "7d",
  "results": [
    {
      "path": "Notes/golang-concurrency.md",
      "title": "Go Concurrency Patterns",
      "snippet": "…goroutines and channels…",
      "score": 0.0162,
      "mod_time": 1740000000,
      "age_days": 42
    }
  ]
}
```

**Orchestrator integration:**
To inject relevant vault context at mission start, call:
```bash
obsidian resurface "<mission topic>" --json --limit 5
```
Parse the `results` array and prepend titles + snippets to the mission prompt as "Related past notes:".

### Index

Build or update the search index. Required before semantic search works.

```bash
# Incremental — only re-indexes notes changed since last run
obsidian index

# Full re-index — reprocesses all notes
obsidian index --force
```

### Doctor

Validate installation and configuration. Run this first when troubleshooting.

```bash
# Human-readable health check
obsidian doctor

# JSON output (for scripting)
obsidian doctor --json
```

**Doctor checks:**
- Binary path and version
- Config file presence and permissions (`~/.obsidian/config`, must be 600)
- Gemini API key presence
- Vault path existence
- Search index (note count, file size)
- Inbox backlog (pending triage count, age of oldest note)

## Configuration

Config file: `~/.obsidian/config`

```ini
vault_path=/path/to/your/vault
gemini_apikey=your-gemini-api-key
```

Or use environment variables: `OBSIDIAN_VAULT_PATH`, `GEMINI_API_KEY`

## Examples

**User**: "capture this idea: use RRF for hybrid search ranking"
**Action**: `obsidian capture "use RRF for hybrid search ranking"`

**User**: "capture this link for later reading" + URL
**Action**: `obsidian capture "interesting article on Go concurrency" --source https://example.com/article`

**User**: "process my inbox"
**Action**: `obsidian triage --auto`

**User**: "what's pending in my inbox?"
**Action**: `obsidian triage`

**User**: "add to my daily note: finished thesis chapter 3"
**Action**: `obsidian append "Daily/2026/02/2026-02-21" "- Finished thesis chapter 3"`

**User**: "add a task under the Tasks section of today's note"
**Action**: `obsidian append "Daily/2026/02/2026-02-21" --section "## Tasks" "- [ ] Review PR #42"`

**User**: "create a new meeting note"
**Action**: `obsidian create "Meetings/2026-02-21 — Team Sync" --title "Team Sync" --type note --status draft`

**User**: "what notes do I have about routines"
**Action**: `obsidian search "routines"` or `obsidian search "household routines" --mode semantic`

**User**: "show me my hub files"
**Action**: `obsidian list "Hubs"`

**User**: "read my personal context"
**Action**: `obsidian read "System/Personal Context — Joey Hipolito"`

**User**: "is my obsidian setup working?"
**Action**: `obsidian doctor`

**User**: "rebuild the search index"
**Action**: `obsidian index --force`

All commands support `--json` for machine-readable output.
