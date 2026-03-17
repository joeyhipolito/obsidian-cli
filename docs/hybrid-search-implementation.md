# Hybrid Search — Implementation Reference

**Last updated:** 2026-03-18
**See also:** [`docs/eval-reports/hybrid-search-eval.md`](eval-reports/hybrid-search-eval.md) — quality evaluation and known bugs

---

## Overview

Search is implemented in `internal/index/` as a three-mode system:

| Mode | Method | Source |
|------|--------|--------|
| `keyword` | FTS5 + BM25 | `store.go:SearchKeyword` |
| `semantic` | Gemini embeddings + cosine similarity | `store.go:SearchSemantic` |
| `hybrid` | Reciprocal Rank Fusion of both | `store.go:SearchHybrid` |

The user-facing entry point is `internal/cmd/search.go:SearchCmd`, which resolves the API key, embeds the query, and dispatches to the appropriate store method.

---

## 1. FTS5 Configuration

**File:** `internal/index/store.go:65–123`

### Schema

```sql
-- Primary data table
CREATE TABLE IF NOT EXISTS notes (
    path      TEXT PRIMARY KEY,
    title     TEXT NOT NULL DEFAULT '',
    tags      TEXT NOT NULL DEFAULT '',
    headings  TEXT NOT NULL DEFAULT '',
    wikilinks TEXT NOT NULL DEFAULT '',
    body      TEXT NOT NULL DEFAULT '',
    mod_time  INTEGER NOT NULL DEFAULT 0,
    embedding BLOB
)

-- FTS5 virtual table (external content mode)
CREATE VIRTUAL TABLE IF NOT EXISTS notes_fts USING fts5(
    path,
    title,
    tags,
    headings,
    body,
    content='notes',
    content_rowid='rowid'
)
```

**Indexed columns:** `path`, `title`, `tags`, `headings`, `body` (5 columns; `wikilinks` is stored but not indexed).

**Tokenizer:** Default (`unicode61`) — no custom tokenizer configured.

**Ranking function:** FTS5's built-in `rank` column, which implements BM25 with default parameters (k1=1.2, b=0.75).

### Triggers

Three AFTER triggers keep the FTS5 index in sync with the `notes` table:

```sql
-- INSERT
CREATE TRIGGER IF NOT EXISTS notes_ai AFTER INSERT ON notes BEGIN
    INSERT INTO notes_fts(rowid, path, title, tags, headings, body)
    VALUES (new.rowid, new.path, new.title, new.tags, new.headings, new.body);
END

-- DELETE
CREATE TRIGGER IF NOT EXISTS notes_ad AFTER DELETE ON notes BEGIN
    INSERT INTO notes_fts(notes_fts, rowid, path, title, tags, headings, body)
    VALUES ('delete', old.rowid, old.path, old.title, old.tags, old.headings, old.body);
END

-- UPDATE
CREATE TRIGGER IF NOT EXISTS notes_au AFTER UPDATE ON notes BEGIN
    INSERT INTO notes_fts(notes_fts, rowid, path, title, tags, headings, body)
    VALUES ('delete', old.rowid, old.path, old.title, old.tags, old.headings, old.body);
    INSERT INTO notes_fts(rowid, path, title, tags, headings, body)
    VALUES (new.rowid, new.path, new.title, new.tags, new.headings, new.body);
END
```

### Database Pragmas

WAL mode is enabled at open time (`store.go:44`):

```go
db.Exec("PRAGMA journal_mode=WAL")
```

---

## 2. Gemini Embedding Integration

**Files:** `internal/index/embeddings.go`, `internal/cmd/index.go`

### Client

```go
// embeddings.go:50–59
type EmbeddingClient struct {
    apiKey     string
    model      string       // "gemini-embedding-001"
    httpClient *http.Client // 30s timeout
}
```

**Model:** `gemini-embedding-001`
**Dimensions:** 768 (set explicitly in each request via `outputDimensionality`)
**Authentication:** API key passed as a query parameter (`?key=...`), not as a Bearer token header.

### Single-item embed

`embeddings.go:67–120` — `EmbeddingClient.Embed(ctx, text) ([]float32, error)`

POST to:
```
https://generativelanguage.googleapis.com/v1beta/models/gemini-embedding-001:embedContent?key=<KEY>
```

Request body:
```json
{
  "model": "models/gemini-embedding-001",
  "content": { "parts": [{ "text": "<text>" }] },
  "outputDimensionality": 768
}
```

### Batch embed

`embeddings.go:124–196` — `EmbeddingClient.EmbedBatch(ctx, texts) ([][]float32, error)`

POST to:
```
https://generativelanguage.googleapis.com/v1beta/models/gemini-embedding-001:batchEmbedContents?key=<KEY>
```

**Batch size:** 100 notes per API call (`internal/cmd/index.go:16`).

### Text preparation for embedding

`store.go:415–435` — `BuildSearchText(title, tags, headings, body) string`

Concatenates `title`, `tags`, `headings`, and `body` (body truncated to 8,000 characters) separated by newlines. This is the same text representation used for both indexing and query embedding.

### Vector storage

Embeddings are stored in the `embedding BLOB` column of the `notes` table as little-endian IEEE 754 float32 bytes:

```go
// store.go:328–347
func encodeEmbedding(v []float32) []byte {
    buf := make([]byte, len(v)*4)
    for i, f := range v {
        binary.LittleEndian.PutUint32(buf[i*4:], math.Float32bits(f))
    }
    return buf
}
```

768 dimensions × 4 bytes = **3,072 bytes per note**.

---

## 3. Scoring Weights and Ranking Algorithm

### Keyword scoring (BM25 via FTS5 `rank`)

`store.go:197–223` — `SearchKeyword`

```sql
SELECT n.path, n.title, rank, snippet(notes_fts, 4, '»', '«', '…', 32)
FROM notes_fts
JOIN notes n ON notes_fts.path = n.path
WHERE notes_fts MATCH ?
ORDER BY rank
LIMIT ?
```

- FTS5 `rank` is negative (more negative = better match). The code negates it: `r.Score = -r.Score`.
- **Observed range:** 0.25–11.85 (unbounded; the comment "normalize to 0-1" is incorrect — see eval report §3.5).
- Snippet uses column index 4 (`body`), with markers `»`/`«`, ellipsis `…`, and a 32-token window.

### Semantic scoring (cosine similarity)

`store.go:225–266` — `SearchSemantic`

All embeddings are loaded into memory and cosine similarity is computed against the query vector:

```go
// store.go:349–368
func CosineSimilarity(a, b []float32) float32 {
    var dotProduct, normA, normB float32
    for i := range a {
        dotProduct += a[i] * b[i]
        normA += a[i] * a[i]
        normB += b[i] * b[i]
    }
    return dotProduct / (sqrt32(normA) * sqrt32(normB))
}
```

Square root uses Newton's method (10 iterations).

- **Threshold:** `score > 0` — effectively no threshold.
- **Observed range:** 0.45–0.73 (calibrated against a 117-note vault).
- Results are sorted descending by score.
- **Note:** Semantic results have no snippet text.

### Hybrid ranking (Reciprocal Rank Fusion)

`store.go:268–316` — `SearchHybrid`

```
RRF score = Σ [ 1.0 / (k + rank_i) ]
```

where **k = 60** and `rank_i` is the 1-based position in each result list.

Each mode is queried for `limit*2` results. For each note, contributions from keyword and semantic rank lists are summed:

```go
const k = 60.0
scores := make(map[string]float64)

for i, r := range keywordResults {
    scores[r.Path] += 1.0 / (k + float64(i+1))
}
for i, r := range semanticResults {
    scores[r.Path] += 1.0 / (k + float64(i+1))
}
```

**Weights:** Equal — both keyword and semantic rankings contribute identically. There is no per-mode weight coefficient.

**Score range observed:** 0.0152–0.0328.

**Snippet propagation:** Only results that appeared in keyword results receive snippets; semantic-only results in hybrid mode have empty snippet fields.

### Sort

`store.go:318–326` — insertion sort descending by score. Appropriate for the small result sets (≤ 20 items) in practice.

---

## 4. Search API Surface

### `internal/index` package

#### Types

```go
type Store struct { db *sql.DB }

type SearchResult struct {
    Path    string  `json:"path"`
    Title   string  `json:"title"`
    Score   float64 `json:"score"`
    Snippet string  `json:"snippet"`
}

type NoteRow struct {
    Path      string
    Title     string
    Tags      string
    Headings  string
    Wikilinks string
    Body      string
    ModTime   int64
    Embedding []float32
}
```

#### Store methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `Open` | `(dbPath string) (*Store, error)` | Opens/creates the SQLite index; creates schema; enables WAL |
| `Close` | `() error` | Closes the database |
| `UpsertNote` | `(note *NoteRow) error` | Insert or update a note and its embedding |
| `DeleteNote` | `(path string) error` | Remove a note from the index |
| `GetModTime` | `(path string) (int64, error)` | Last-indexed modification time for incremental updates |
| `GetAllPaths` | `() (map[string]bool, error)` | All indexed note paths |
| `GetAllNoteRows` | `() ([]NoteRow, error)` | All notes with full metadata and embeddings |
| `NoteCount` | `() (int, error)` | Total notes indexed |
| `EmbeddingCount` | `() (int, error)` | Notes that have an embedding |
| `SearchKeyword` | `(query string, limit int) ([]SearchResult, error)` | FTS5 keyword search |
| `SearchSemantic` | `(queryEmbedding []float32, limit int) ([]SearchResult, error)` | Cosine similarity search |
| `SearchHybrid` | `(query string, queryEmbedding []float32, limit int) ([]SearchResult, error)` | RRF hybrid search |

#### Helper functions

```go
func BuildSearchText(title, tags, headings, body string) string
func CosineSimilarity(a, b []float32) float32
func encodeEmbedding(v []float32) []byte   // unexported
func decodeEmbedding(b []byte) []float32   // unexported
```

#### `EmbeddingClient` methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `NewEmbeddingClient` | `(apiKey string) *EmbeddingClient` | Constructor |
| `IsAvailable` | `() bool` | Returns true if API key is non-empty |
| `Embed` | `(ctx, text) ([]float32, error)` | Embed a single string |
| `EmbedBatch` | `(ctx, texts) ([][]float32, error)` | Embed a batch of strings (100 per API call) |

### `internal/cmd` package — user-facing command

```go
// search.go
type SearchOutput struct {
    Query   string               `json:"query"`
    Mode    string               `json:"mode"`
    Results []index.SearchResult `json:"results"`
}

func SearchCmd(vaultPath, query, mode string, jsonOutput bool) error
```

**`mode` values:** `"keyword"`, `"semantic"`, `"hybrid"` (default).

**Fallback behaviour in `SearchCmd`:**
1. If `mode == "hybrid"` and embedding fails → falls back to keyword with a warning.
2. If `mode == "hybrid"` and no API key → falls back to keyword with a warning.
3. If `mode == "semantic"` and no API key → returns an error (no fallback).

**Result limit:** Hardcoded to 20 (`search.go:43`).

---

## 5. Key Constants

| Constant | Value | Location | Purpose |
|----------|-------|----------|---------|
| `EmbeddingDimensions` | `768` | `store.go:17` | Gemini output size |
| `batchSize` | `100` | `cmd/index.go:16` | Notes per batch embed API call |
| `limit` | `20` | `cmd/search.go:43` | Max results returned |
| `k` (RRF) | `60.0` | `store.go:282` | RRF smoothing parameter |

---

## 6. Known Issues

See [`docs/eval-reports/hybrid-search-eval.md`](eval-reports/hybrid-search-eval.md) for the full evaluation. Summary:

| Severity | Issue | Fix |
|----------|-------|-----|
| P0 | FTS5 query injection — `#`, `/`, `[`, `NOT`, etc. crash search | Quote user query: `"` + escape `"` → `""` |
| P0 | Hybrid error propagation — keyword failure kills hybrid instead of degrading | Catch `kwErr`, continue with semantic-only |
| P1 | No semantic score threshold — nonsense queries return 20 results | Add `score > 0.55` cutoff |
| P1 | RRF ranking regression — keyword noise buries correct semantic result | Weighted RRF or semantic-boost |
| P2 | Missing snippets for semantic-only hybrid results | Post-process top results to generate snippets |
| P2 | FTS5 score not normalized to [0,1] despite code comment | Cosmetic; doesn't affect hybrid ranking |
| P3 | Brute-force cosine scan — loads all embeddings into memory | ANN index (sqlite-vss / HNSW) at 1,000+ notes |
