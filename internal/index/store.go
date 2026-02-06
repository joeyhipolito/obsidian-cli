// Package index provides SQLite-backed full-text and vector search indexing
// for Obsidian vault notes. Uses FTS5 for keyword search and vector embeddings
// stored as blobs for semantic search.
package index

import (
	"database/sql"
	"encoding/binary"
	"fmt"
	"math"
	"strings"

	_ "modernc.org/sqlite"
)

// EmbeddingDimensions is the size of Gemini text-embedding-004 vectors.
const EmbeddingDimensions = 768

// Store manages the SQLite search index for an Obsidian vault.
type Store struct {
	db *sql.DB
}

// NoteRow represents a row in the notes table.
type NoteRow struct {
	Path      string
	Title     string
	Tags      string // comma-separated
	Headings  string // newline-separated
	Wikilinks string // comma-separated
	Body      string
	ModTime   int64
	Embedding []float32
}

// Open opens or creates the SQLite index database at the given path.
// Creates the schema if it doesn't exist.
func Open(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open index database: %w", err)
	}

	// Enable WAL mode for better concurrent read performance
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
	}

	s := &Store{db: db}
	if err := s.createSchema(); err != nil {
		db.Close()
		return nil, err
	}

	return s, nil
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// createSchema creates the tables if they don't exist.
func (s *Store) createSchema() error {
	// Main notes table with metadata and vector embedding blob
	_, err := s.db.Exec(`
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
	`)
	if err != nil {
		return fmt.Errorf("failed to create notes table: %w", err)
	}

	// FTS5 virtual table for keyword search over title, tags, headings, body
	_, err = s.db.Exec(`
		CREATE VIRTUAL TABLE IF NOT EXISTS notes_fts USING fts5(
			path,
			title,
			tags,
			headings,
			body,
			content='notes',
			content_rowid='rowid'
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create FTS5 table: %w", err)
	}

	// Triggers to keep FTS5 in sync with the notes table
	triggers := []string{
		`CREATE TRIGGER IF NOT EXISTS notes_ai AFTER INSERT ON notes BEGIN
			INSERT INTO notes_fts(rowid, path, title, tags, headings, body)
			VALUES (new.rowid, new.path, new.title, new.tags, new.headings, new.body);
		END`,
		`CREATE TRIGGER IF NOT EXISTS notes_ad AFTER DELETE ON notes BEGIN
			INSERT INTO notes_fts(notes_fts, rowid, path, title, tags, headings, body)
			VALUES ('delete', old.rowid, old.path, old.title, old.tags, old.headings, old.body);
		END`,
		`CREATE TRIGGER IF NOT EXISTS notes_au AFTER UPDATE ON notes BEGIN
			INSERT INTO notes_fts(notes_fts, rowid, path, title, tags, headings, body)
			VALUES ('delete', old.rowid, old.path, old.title, old.tags, old.headings, old.body);
			INSERT INTO notes_fts(rowid, path, title, tags, headings, body)
			VALUES (new.rowid, new.path, new.title, new.tags, new.headings, new.body);
		END`,
	}
	for _, t := range triggers {
		if _, err := s.db.Exec(t); err != nil {
			return fmt.Errorf("failed to create trigger: %w", err)
		}
	}

	return nil
}

// GetModTime returns the stored mod_time for a note path, or 0 if not indexed.
func (s *Store) GetModTime(path string) (int64, error) {
	var modTime int64
	err := s.db.QueryRow("SELECT mod_time FROM notes WHERE path = ?", path).Scan(&modTime)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return modTime, err
}

// GetAllPaths returns all indexed note paths.
func (s *Store) GetAllPaths() (map[string]bool, error) {
	rows, err := s.db.Query("SELECT path FROM notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	paths := make(map[string]bool)
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, err
		}
		paths[path] = true
	}
	return paths, rows.Err()
}

// UpsertNote inserts or updates a note in the index.
func (s *Store) UpsertNote(note *NoteRow) error {
	var embBlob []byte
	if note.Embedding != nil {
		embBlob = encodeEmbedding(note.Embedding)
	}

	_, err := s.db.Exec(`
		INSERT INTO notes (path, title, tags, headings, wikilinks, body, mod_time, embedding)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(path) DO UPDATE SET
			title     = excluded.title,
			tags      = excluded.tags,
			headings  = excluded.headings,
			wikilinks = excluded.wikilinks,
			body      = excluded.body,
			mod_time  = excluded.mod_time,
			embedding = excluded.embedding
	`, note.Path, note.Title, note.Tags, note.Headings, note.Wikilinks, note.Body, note.ModTime, embBlob)
	return err
}

// DeleteNote removes a note from the index.
func (s *Store) DeleteNote(path string) error {
	_, err := s.db.Exec("DELETE FROM notes WHERE path = ?", path)
	return err
}

// NoteCount returns the total number of indexed notes.
func (s *Store) NoteCount() (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM notes").Scan(&count)
	return count, err
}

// SearchResult holds a single search match.
type SearchResult struct {
	Path    string  `json:"path"`
	Title   string  `json:"title"`
	Score   float64 `json:"score"`
	Snippet string  `json:"snippet"`
}

// SearchKeyword performs an FTS5 keyword search.
func (s *Store) SearchKeyword(query string, limit int) ([]SearchResult, error) {
	rows, err := s.db.Query(`
		SELECT n.path, n.title, rank, snippet(notes_fts, 4, '»', '«', '…', 32)
		FROM notes_fts
		JOIN notes n ON notes_fts.path = n.path
		WHERE notes_fts MATCH ?
		ORDER BY rank
		LIMIT ?
	`, query, limit)
	if err != nil {
		return nil, fmt.Errorf("FTS5 search failed: %w", err)
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var r SearchResult
		if err := rows.Scan(&r.Path, &r.Title, &r.Score, &r.Snippet); err != nil {
			return nil, err
		}
		// FTS5 rank is negative (lower = better), normalize to 0-1 range
		r.Score = -r.Score
		results = append(results, r)
	}
	return results, rows.Err()
}

// SearchSemantic performs vector similarity search using cosine similarity.
func (s *Store) SearchSemantic(queryEmbedding []float32, limit int) ([]SearchResult, error) {
	rows, err := s.db.Query("SELECT path, title, embedding FROM notes WHERE embedding IS NOT NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var path, title string
		var embBlob []byte
		if err := rows.Scan(&path, &title, &embBlob); err != nil {
			return nil, err
		}

		emb := decodeEmbedding(embBlob)
		if emb == nil {
			continue
		}

		score := cosineSimilarity(queryEmbedding, emb)
		if score > 0 {
			results = append(results, SearchResult{
				Path:  path,
				Title: title,
				Score: float64(score),
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Sort by score descending
	sortResults(results)

	if len(results) > limit {
		results = results[:limit]
	}
	return results, nil
}

// SearchHybrid combines FTS5 keyword and semantic vector search with RRF ranking.
func (s *Store) SearchHybrid(query string, queryEmbedding []float32, limit int) ([]SearchResult, error) {
	// Get both result sets
	keywordResults, err := s.SearchKeyword(query, limit*2)
	if err != nil {
		return nil, err
	}

	semanticResults, err := s.SearchSemantic(queryEmbedding, limit*2)
	if err != nil {
		return nil, err
	}

	// Reciprocal Rank Fusion (RRF) with k=60
	const k = 60.0
	scores := make(map[string]float64)
	titles := make(map[string]string)
	snippets := make(map[string]string)

	for i, r := range keywordResults {
		scores[r.Path] += 1.0 / (k + float64(i+1))
		titles[r.Path] = r.Title
		snippets[r.Path] = r.Snippet
	}
	for i, r := range semanticResults {
		scores[r.Path] += 1.0 / (k + float64(i+1))
		if titles[r.Path] == "" {
			titles[r.Path] = r.Title
		}
	}

	// Build combined results
	var results []SearchResult
	for path, score := range scores {
		results = append(results, SearchResult{
			Path:    path,
			Title:   titles[path],
			Score:   score,
			Snippet: snippets[path],
		})
	}

	sortResults(results)

	if len(results) > limit {
		results = results[:limit]
	}
	return results, nil
}

// sortResults sorts search results by score descending.
func sortResults(results []SearchResult) {
	// Simple insertion sort — result sets are small
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Score > results[j-1].Score; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}
}

// encodeEmbedding converts a float32 slice to a byte slice (little-endian).
func encodeEmbedding(v []float32) []byte {
	buf := make([]byte, len(v)*4)
	for i, f := range v {
		binary.LittleEndian.PutUint32(buf[i*4:], math.Float32bits(f))
	}
	return buf
}

// decodeEmbedding converts a byte slice back to a float32 slice.
func decodeEmbedding(b []byte) []float32 {
	if len(b) == 0 || len(b)%4 != 0 {
		return nil
	}
	v := make([]float32, len(b)/4)
	for i := range v {
		v[i] = math.Float32frombits(binary.LittleEndian.Uint32(b[i*4:]))
	}
	return v
}

// cosineSimilarity calculates the cosine similarity between two vectors.
// Returns a value between -1 and 1, where 1 means identical.
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dotProduct, normA, normB float32
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (sqrt32(normA) * sqrt32(normB))
}

// sqrt32 computes float32 square root using Newton's method.
func sqrt32(x float32) float32 {
	if x <= 0 {
		return 0
	}
	z := x / 2
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return z
}

// IndexDBPath returns the path to the index database for a given vault.
func IndexDBPath(vaultPath string) string {
	return vaultPath + "/.obsidian/search.db"
}

// BuildSearchText creates a combined text for embedding from note fields.
func BuildSearchText(title, tags, headings, body string) string {
	var parts []string
	if title != "" {
		parts = append(parts, title)
	}
	if tags != "" {
		parts = append(parts, tags)
	}
	if headings != "" {
		parts = append(parts, headings)
	}
	if body != "" {
		// Truncate body to ~8000 chars for embedding API limits
		if len(body) > 8000 {
			body = body[:8000]
		}
		parts = append(parts, body)
	}
	return strings.Join(parts, "\n")
}
