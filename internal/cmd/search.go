package cmd

import (
	"fmt"

	"github.com/joeyhipolito/obsidian-cli/internal/output"
)

// SearchResult represents a single search hit.
type SearchResult struct {
	Path    string  `json:"path"`
	Score   float64 `json:"score"`
	Snippet string  `json:"snippet"`
}

// SearchOutput represents the JSON output format for the search command.
type SearchOutput struct {
	Query   string         `json:"query"`
	Results []SearchResult `json:"results"`
}

// SearchCmd searches notes using keyword (FTS5) and semantic (vector) search.
func SearchCmd(vaultPath, query string, jsonOutput bool) error {
	// TODO: implement — hybrid search: FTS5 keyword + Gemini vector similarity
	if jsonOutput {
		return output.JSON(SearchOutput{
			Query:   query,
			Results: []SearchResult{},
		})
	}

	fmt.Printf("obsidian search: not yet implemented (query: %s)\n", query)
	return nil
}
