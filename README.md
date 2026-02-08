# obsidian-cli

A Go CLI for managing and searching [Obsidian](https://obsidian.md/) vault notes from the terminal. Supports keyword, semantic, and hybrid search powered by Google Gemini embeddings.

## Features

- **Read/write notes** — read, create, and append to markdown notes with frontmatter parsing
- **Full-text search** — SQLite FTS5 keyword search with ranked results
- **Semantic search** — vector similarity search using Gemini embeddings (768-dim)
- **Hybrid search** — combines keyword + semantic with Reciprocal Rank Fusion (RRF)
- **Incremental indexing** — only re-indexes changed files
- **Frontmatter parsing** — extracts YAML metadata, headings, and wikilinks
- **Interactive configuration** — `obsidian configure` setup
- **Diagnostics** — built-in `doctor` command for troubleshooting
- **JSON output** — machine-readable format for scripting (`--json`)
- **Cross-platform** — macOS (arm64/amd64) and Linux (amd64/arm64)

## Installation

### Prerequisites

- Go 1.25 or later
- An Obsidian vault directory
- (Optional) Gemini API key for semantic search ([get one here](https://aistudio.google.com/api-keys))

### Build and Install

```bash
make install            # Build and symlink to ~/bin
obsidian configure      # Interactive setup
obsidian doctor         # Verify everything works
```

## Configuration

Config is stored in `~/.obsidian/config` (INI format, `chmod 600`).

```bash
obsidian configure          # Interactive setup (recommended)
obsidian configure show     # Show current config (API key masked)
```

### Config file keys

| Key | Description |
|-----|-------------|
| `gemini_apikey` | Gemini API key (required for semantic/hybrid search) |
| `vault_path` | Path to your Obsidian vault |

### Environment variables (fallback)

| Variable | Description |
|----------|-------------|
| `GEMINI_API_KEY` | Gemini API key |
| `OBSIDIAN_VAULT_PATH` | Vault directory path |

## Commands

### Reading notes

```bash
obsidian read "Daily/2024-01-15.md"         # Read note content
obsidian read "Projects/ideas.md" --json    # Structured output (frontmatter, headings, wikilinks)
```

### Creating and appending

```bash
obsidian create "Notes/meeting.md" --title "Sprint Planning"
obsidian append "Notes/meeting.md" "Action item: review PR #42"
echo "piped content" | obsidian append "Notes/log.md"
```

### Listing notes

```bash
obsidian list                   # All notes in vault
obsidian list "Projects/"       # Notes in a subdirectory
obsidian list --json            # JSON with path, name, size, modified time
```

### Searching

```bash
# Keyword search (FTS5)
obsidian search "golang error handling" --mode keyword

# Semantic search (requires Gemini API key)
obsidian search "how to handle errors in Go" --mode semantic

# Hybrid search (default — combines both with RRF)
obsidian search "error handling patterns"
```

### Building the search index

```bash
obsidian index      # Build/update index (incremental)
```

The index is stored at `<vault>/.obsidian/search.db` (SQLite). Incremental indexing skips unchanged files and removes deleted notes.

### Diagnostics

```bash
obsidian doctor     # Validate config, vault access, index status, API key
```

## Architecture

```
cmd/obsidian-cli/            # Entry point and command routing
internal/
├── cmd/                     # Command implementations
│   ├── read.go              # Read note with frontmatter parsing
│   ├── append.go            # Append text to notes
│   ├── create.go            # Create new notes
│   ├── list.go              # List vault files
│   ├── search.go            # Search (keyword/semantic/hybrid)
│   ├── index.go             # Build/update search index
│   ├── configure.go         # Configuration management
│   └── doctor.go            # Diagnostics
├── config/                  # Config file loading/saving
├── vault/                   # Note I/O and markdown parsing
│   ├── vault.go             # ReadNote, WriteNote, AppendToNote, ListNotes
│   └── parse.go             # YAML frontmatter, wikilinks, headings
├── index/                   # Search index
│   ├── store.go             # SQLite FTS5 + vector storage
│   └── embeddings.go        # Gemini embedding API client
└── output/                  # JSON output helpers
```

### Design decisions

- **Hybrid search by default** — keyword search for precision, semantic for meaning, RRF to combine
- **Pure-Go SQLite** — uses `modernc.org/sqlite` (no CGO required)
- **Custom YAML parser** — lightweight frontmatter parsing without external YAML library
- **Batch embeddings** — processes up to 100 texts per Gemini API request
- **Cosine similarity** — computed in-memory over float32 vectors (scales to hundreds of notes)

## Development

```bash
make build              # Build for current platform
make install            # Build and install to ~/bin
make test               # Run tests
make test-coverage      # Tests with HTML coverage report
make build-all          # Cross-compile for macOS/Linux
make fmt                # Format code
make vet                # Vet for issues
make lint               # fmt + vet
make clean              # Remove build artifacts
```

## License

MIT License. See [LICENSE](LICENSE) for details.

## Links

- [Obsidian](https://obsidian.md/)
- [Gemini API Keys](https://aistudio.google.com/api-keys)
