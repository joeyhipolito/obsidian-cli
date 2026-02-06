package cmd

import (
	"fmt"

	"github.com/joeyhipolito/obsidian-cli/internal/output"
)

// IndexOutput represents the JSON output format for the index command.
type IndexOutput struct {
	NotesIndexed int `json:"notes_indexed"`
	NotesSkipped int `json:"notes_skipped"`
	Errors       int `json:"errors"`
}

// IndexCmd builds or updates the SQLite search index for the vault.
// Crawls vault, parses frontmatter/headings/wikilinks, builds FTS5 index,
// and generates Gemini vector embeddings for semantic search.
func IndexCmd(vaultPath string, jsonOutput bool) error {
	// TODO: implement — crawl vault, parse notes, build SQLite FTS5 + vector index
	if jsonOutput {
		return output.JSON(IndexOutput{
			NotesIndexed: 0,
			NotesSkipped: 0,
			Errors:       0,
		})
	}

	fmt.Printf("obsidian index: not yet implemented (vault: %s)\n", vaultPath)
	return nil
}
