package cmd

import (
	"testing"

	"github.com/joeyhipolito/obsidian-cli/internal/vault"
)

func TestCountOrphans(t *testing.T) {
	notes := []vault.NoteInfo{
		{Path: "Ideas/my-idea.md", Name: "my-idea"},
		{Path: "Notes/meeting.md", Name: "meeting"},
		{Path: "References/article.md", Name: "article"},
	}

	tests := []struct {
		name          string
		inboundLinks  map[string]int
		wantOrphans   int
	}{
		{
			name:         "all orphans when no links",
			inboundLinks: map[string]int{},
			wantOrphans:  3,
		},
		{
			name: "note linked by name reduces orphan count",
			inboundLinks: map[string]int{
				"my-idea": 2,
				"article": 1,
			},
			wantOrphans: 1, // only meeting is orphaned
		},
		{
			name: "note linked by path reduces orphan count",
			inboundLinks: map[string]int{
				"notes/meeting": 1,
			},
			wantOrphans: 2, // my-idea and article are orphaned
		},
		{
			name: "all notes linked",
			inboundLinks: map[string]int{
				"my-idea": 1,
				"meeting": 1,
				"article": 1,
			},
			wantOrphans: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countOrphans(notes, tt.inboundLinks)
			if got != tt.wantOrphans {
				t.Errorf("countOrphans() = %d, want %d", got, tt.wantOrphans)
			}
		})
	}
}

func TestNoteClassificationDist(t *testing.T) {
	tests := []struct {
		name  string
		notes []vault.NoteInfo
		want  map[string]int
	}{
		{
			name:  "empty vault",
			notes: []vault.NoteInfo{},
			want:  map[string]int{},
		},
		{
			name: "notes in multiple folders",
			notes: []vault.NoteInfo{
				{Path: "Ideas/idea-a.md"},
				{Path: "Ideas/idea-b.md"},
				{Path: "Notes/note-a.md"},
				{Path: "Inbox/capture.md"},
			},
			want: map[string]int{
				"Ideas": 2,
				"Notes": 1,
				"Inbox": 1,
			},
		},
		{
			name: "notes at vault root counted as Root",
			notes: []vault.NoteInfo{
				{Path: "readme.md"},
				{Path: "index.md"},
				{Path: "Notes/note.md"},
			},
			want: map[string]int{
				"Root":  2,
				"Notes": 1,
			},
		},
		{
			name: "deeply nested notes use top-level folder",
			notes: []vault.NoteInfo{
				{Path: "Daily/2026/03/2026-03-17.md"},
				{Path: "Daily/2026/03/2026-03-16.md"},
				{Path: "Projects/obsidian/health.md"},
			},
			want: map[string]int{
				"Daily":    2,
				"Projects": 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := noteClassificationDist(tt.notes)
			if len(got) != len(tt.want) {
				t.Errorf("noteClassificationDist() len = %d, want %d; got %v", len(got), len(tt.want), got)
				return
			}
			for folder, wantCount := range tt.want {
				if got[folder] != wantCount {
					t.Errorf("noteClassificationDist()[%q] = %d, want %d", folder, got[folder], wantCount)
				}
			}
		})
	}
}

func TestAvgWikilinkDensity(t *testing.T) {
	tests := []struct {
		name       string
		totalLinks int
		noteCount  int
		want       float64
	}{
		{
			name:       "zero notes returns zero",
			totalLinks: 0,
			noteCount:  0,
			want:       0,
		},
		{
			name:       "no links in vault",
			totalLinks: 0,
			noteCount:  10,
			want:       0,
		},
		{
			name:       "even distribution",
			totalLinks: 20,
			noteCount:  4,
			want:       5.0,
		},
		{
			name:       "fractional density",
			totalLinks: 7,
			noteCount:  3,
			want:       7.0 / 3.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := avgWikilinkDensity(tt.totalLinks, tt.noteCount)
			if got != tt.want {
				t.Errorf("avgWikilinkDensity(%d, %d) = %v, want %v", tt.totalLinks, tt.noteCount, got, tt.want)
			}
		})
	}
}
