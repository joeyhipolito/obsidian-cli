package cmd

import (
	"testing"

	"github.com/joeyhipolito/obsidian-cli/internal/vault"
)

func TestClassifyNoteType(t *testing.T) {
	tests := []struct {
		name     string
		fm       map[string]any
		body     string
		headings []vault.Heading
		want     string
	}{
		{
			name: "preserves existing non-fleeting type",
			fm:   map[string]any{"type": "reference"},
			body: "some content",
			want: "reference",
		},
		{
			name: "reclassifies fleeting with task signals",
			fm:   map[string]any{"type": "fleeting"},
			body: "- [ ] buy milk\n- [ ] call dentist\n",
			want: "task",
		},
		{
			name: "reclassifies fleeting with TODO marker",
			fm:   map[string]any{"type": "fleeting"},
			body: "todo: finish the report\n",
			want: "task",
		},
		{
			name: "reference from source URL",
			fm:   map[string]any{"type": "fleeting", "source": "https://example.com/article"},
			body: "interesting article",
			want: "reference",
		},
		{
			name: "reference from inline URL in body",
			fm:   map[string]any{"type": "fleeting"},
			body: "Check out https://golang.org for docs.",
			want: "reference",
		},
		{
			name: "note from multiple headings",
			fm:   map[string]any{"type": "fleeting"},
			body: "## Background\nSome text.\n## Analysis\nMore text.\n",
			headings: []vault.Heading{
				{Level: 2, Text: "Background"},
				{Level: 2, Text: "Analysis"},
			},
			want: "note",
		},
		{
			name: "defaults to idea",
			fm:   map[string]any{"type": "fleeting"},
			body: "A rough idea about improving the search UI.",
			want: "idea",
		},
		{
			name: "no type in frontmatter → idea",
			fm:   map[string]any{},
			body: "Just a quick thought.",
			want: "idea",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed := &vault.Note{
				Frontmatter: tt.fm,
				Body:        tt.body,
				Headings:    tt.headings,
			}
			got := classifyNoteType(parsed)
			if got != tt.want {
				t.Errorf("classifyNoteType() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTypeFolder(t *testing.T) {
	tests := []struct {
		noteType string
		want     string
	}{
		{"task", "Tasks"},
		{"reference", "References"},
		{"idea", "Ideas"},
		{"note", "Notes"},
		{"unknown", "Notes"},
		{"", "Notes"},
	}
	for _, tt := range tests {
		got := typeFolder(tt.noteType)
		if got != tt.want {
			t.Errorf("typeFolder(%q) = %q, want %q", tt.noteType, got, tt.want)
		}
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Hello World", "hello-world"},
		{"Golang Error Handling", "golang-error-handling"},
		{"  spaces   ", "spaces"},
		{"Special: Characters!", "special-characters"},
		{"", "untitled"},
		{"---", "untitled"},
		{"A  B", "a-b"},
	}
	for _, tt := range tests {
		got := slugify(tt.input)
		if got != tt.want {
			t.Errorf("slugify(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestNoteSlug(t *testing.T) {
	tests := []struct {
		name     string
		fromPath string
		fm       map[string]any
		headings []vault.Heading
		want     string
	}{
		{
			name:     "uses frontmatter title",
			fromPath: "Inbox/20260221-143022.md",
			fm:       map[string]any{"title": "My Great Idea"},
			want:     "my-great-idea",
		},
		{
			name:     "uses first H1 heading",
			fromPath: "Inbox/20260221-143022.md",
			fm:       map[string]any{},
			headings: []vault.Heading{{Level: 1, Text: "Article Title"}},
			want:     "article-title",
		},
		{
			name:     "falls back to filename",
			fromPath: "Inbox/20260221-143022.md",
			fm:       map[string]any{},
			want:     "20260221-143022",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed := &vault.Note{
				Frontmatter: tt.fm,
				Headings:    tt.headings,
			}
			got := noteSlug(tt.fromPath, parsed)
			if got != tt.want {
				t.Errorf("noteSlug() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatAgeLabel(t *testing.T) {
	tests := []struct {
		days int
		want string
	}{
		{0, "today"},
		{1, "1 day ago"},
		{5, "5 days ago"},
		{7, "1 week ago"},
		{10, "1 week ago"},
		{14, "2 weeks ago"},
		{30, "1 month ago"},
		{60, "2 months ago"},
	}
	for _, tt := range tests {
		got := formatAgeLabel(tt.days)
		if got != tt.want {
			t.Errorf("formatAgeLabel(%d) = %q, want %q", tt.days, got, tt.want)
		}
	}
}

func TestFrontmatterString(t *testing.T) {
	fm := map[string]any{
		"type":   "fleeting",
		"source": "https://example.com",
		"count":  42,
	}

	if got := frontmatterString(fm, "type"); got != "fleeting" {
		t.Errorf("expected %q, got %q", "fleeting", got)
	}
	if got := frontmatterString(fm, "missing"); got != "" {
		t.Errorf("expected empty string for missing key, got %q", got)
	}
	// Non-string value returns empty string (type assertion fails gracefully).
	if got := frontmatterString(fm, "count"); got != "" {
		t.Errorf("expected empty string for non-string value, got %q", got)
	}
}
