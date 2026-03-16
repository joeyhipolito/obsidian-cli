package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/joeyhipolito/obsidian-cli/internal/index"
	"github.com/joeyhipolito/obsidian-cli/internal/output"
	"github.com/joeyhipolito/obsidian-cli/internal/vault"
)

// TriageOptions holds flags for the triage command.
type TriageOptions struct {
	List       bool
	Auto       bool
	Older      string // duration string like "7d", "24h" — parsed by parseSinceDuration
	DryRun     bool
	JSONOutput bool
	Quiet      bool // suppress all output when nothing was processed (cron-friendly)
}

// PendingNote represents a note in the inbox awaiting triage.
type PendingNote struct {
	Path     string `json:"path"`
	Type     string `json:"type"`
	AgeDays  int    `json:"age_days"`
	AgeLabel string `json:"age_label"`
	Created  string `json:"created"`
}

// ProcessedNote represents a note that was triaged by --auto.
type ProcessedNote struct {
	FromPath   string   `json:"from_path"`
	ToPath     string   `json:"to_path"`
	NoteType   string   `json:"note_type"`
	LinksAdded []string `json:"links_added,omitempty"`
	DryRun     bool     `json:"dry_run,omitempty"`
}

// TriageSummary holds aggregate counts for the triage run.
type TriageSummary struct {
	Total     int `json:"total"`
	Processed int `json:"processed"`
	Skipped   int `json:"skipped"`
	Errors    int `json:"errors"`
}

// TriageOutput is the full JSON output for the triage command.
type TriageOutput struct {
	Pending   []PendingNote   `json:"pending"`
	Processed []ProcessedNote `json:"processed"`
	Errors    []string        `json:"errors"`
	Summary   TriageSummary   `json:"summary"`
}

// TriageCmd triages notes in the Inbox/ folder.
// --list shows pending notes with age; --auto classifies, enriches, and moves them.
// Default mode (neither flag) is equivalent to --list.
func TriageCmd(vaultPath string, opts TriageOptions) error {
	// Default to --list when neither mode is specified.
	if !opts.List && !opts.Auto {
		opts.List = true
	}

	// Parse --older duration (empty string → 0, meaning no filter).
	olderDuration, err := parseSinceDuration(opts.Older)
	if err != nil {
		return fmt.Errorf("invalid --older value %q: %w", opts.Older, err)
	}

	// Guard: Inbox must exist.
	inboxPath := filepath.Join(vaultPath, "Inbox")
	if _, err := os.Stat(inboxPath); os.IsNotExist(err) {
		if opts.JSONOutput {
			return output.JSON(TriageOutput{})
		}
		fmt.Println("No inbox found (Inbox/ does not exist in vault)")
		return nil
	}

	notes, err := vault.ListNotes(vaultPath, "Inbox")
	if err != nil {
		return fmt.Errorf("listing inbox: %w", err)
	}

	result := TriageOutput{}
	now := time.Now()

	// Collect pending notes (skip already-processed ones; apply --older filter).
	for _, info := range notes {
		fullPath := filepath.Join(vaultPath, info.Path)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", info.Path, err))
			continue
		}

		parsed := vault.ParseNote(string(data))
		status := frontmatterString(parsed.Frontmatter, "status")

		// Notes with status: processed have already been handled; skip them.
		if status == "processed" {
			continue
		}

		noteType := frontmatterString(parsed.Frontmatter, "type")
		created := frontmatterString(parsed.Frontmatter, "created")
		ageDays, createdLabel := computeAge(created, info.ModTime, now)

		// Apply --older filter: skip notes younger than the threshold.
		if olderDuration > 0 {
			thresholdDays := int(olderDuration.Hours() / 24)
			if ageDays < thresholdDays {
				continue
			}
		}

		result.Pending = append(result.Pending, PendingNote{
			Path:     info.Path,
			Type:     noteType,
			AgeDays:  ageDays,
			AgeLabel: formatAgeLabel(ageDays),
			Created:  createdLabel,
		})
	}

	result.Summary.Total = len(result.Pending)

	// --auto: classify, enrich, rewrite frontmatter, move each pending note.
	if opts.Auto {
		// Open index for wikilink enrichment (best-effort; skipped if not built).
		var store *index.Store
		dbPath := index.IndexDBPath(vaultPath)
		if _, err := os.Stat(dbPath); err == nil {
			if s, err := index.Open(dbPath); err == nil {
				store = s
				defer store.Close()
			}
		}

		for _, pending := range result.Pending {
			processed, err := triageNote(vaultPath, pending, store, opts.DryRun, now)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", pending.Path, err))
				result.Summary.Errors++
				continue
			}
			result.Processed = append(result.Processed, processed)
			result.Summary.Processed++
		}
		result.Summary.Skipped = result.Summary.Total - result.Summary.Processed - result.Summary.Errors
	}

	if opts.JSONOutput {
		return output.JSON(result)
	}

	// --quiet: skip all output when the run produced nothing worth reporting.
	// Errors still surface so cron email catches failures.
	if opts.Quiet && result.Summary.Processed == 0 && len(result.Errors) == 0 {
		return nil
	}

	if opts.Auto {
		printTriageAutoReport(result, opts.DryRun)
	} else {
		printTriageListReport(result)
	}
	return nil
}

// triageNote classifies, enriches, rewrites frontmatter, and moves a single note.
func triageNote(vaultPath string, pending PendingNote, store *index.Store, dryRun bool, now time.Time) (ProcessedNote, error) {
	fullPath := filepath.Join(vaultPath, pending.Path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return ProcessedNote{}, fmt.Errorf("reading note: %w", err)
	}

	parsed := vault.ParseNote(string(data))

	// Step 1: Classify note type from content.
	noteType := classifyNoteType(parsed)

	// Step 2: Find wikilink suggestions via index cosine similarity (optional).
	var linksAdded []string
	if store != nil {
		suggestions := enrichSingleNote(store, pending.Path)
		for _, s := range suggestions {
			toName := strings.TrimSuffix(filepath.Base(s.To), ".md")
			linksAdded = append(linksAdded, toName)
		}
	}

	// Step 3: Compute destination path.
	toPath := triageDestination(pending.Path, noteType, parsed)

	// Step 4: Build updated note content.
	newContent := buildTriagedContent(parsed, noteType, linksAdded, now)

	if dryRun {
		return ProcessedNote{
			FromPath:   pending.Path,
			ToPath:     toPath,
			NoteType:   noteType,
			LinksAdded: linksAdded,
			DryRun:     true,
		}, nil
	}

	// Step 5: Write to new path (deconflict if already exists).
	targetFull := filepath.Join(vaultPath, toPath)
	if _, err := os.Stat(targetFull); err == nil {
		toPath = deconflictPath(toPath)
		targetFull = filepath.Join(vaultPath, toPath)
	}

	if err := os.MkdirAll(filepath.Dir(targetFull), 0755); err != nil {
		return ProcessedNote{}, fmt.Errorf("creating target directory: %w", err)
	}
	if err := os.WriteFile(targetFull, []byte(newContent), 0644); err != nil {
		return ProcessedNote{}, fmt.Errorf("writing triaged note: %w", err)
	}

	// Step 6: Remove original.
	if err := os.Remove(fullPath); err != nil {
		return ProcessedNote{}, fmt.Errorf("removing original note: %w", err)
	}

	return ProcessedNote{
		FromPath:   pending.Path,
		ToPath:     toPath,
		NoteType:   noteType,
		LinksAdded: linksAdded,
	}, nil
}

// classifyNoteType determines the best note type for a parsed note.
// Priority order:
//  1. Existing non-fleeting type in frontmatter → preserve it.
//  2. Task signals (checkbox, TODO) → "task".
//  3. Source URL or inline URL → "reference".
//  4. Multiple headings (structured content) → "note".
//  5. Default → "idea".
func classifyNoteType(parsed *vault.Note) string {
	existing := frontmatterString(parsed.Frontmatter, "type")
	if existing != "" && existing != "fleeting" {
		return existing
	}

	body := strings.ToLower(parsed.Body)

	// Task signals: checkbox items or explicit TODO/action markers.
	taskRe := regexp.MustCompile(`(?m)^[-*]\s+\[[ xX]\]|^\s*todo:|\bfollowup\b|^action:`)
	if taskRe.MatchString(body) {
		return "task"
	}

	// Reference signals: URL in source frontmatter or body.
	source := frontmatterString(parsed.Frontmatter, "source")
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		return "reference"
	}
	urlRe := regexp.MustCompile(`https?://\S+`)
	if urlRe.MatchString(body) {
		return "reference"
	}

	// Note signals: structured content with multiple headings.
	if len(parsed.Headings) >= 2 {
		return "note"
	}

	return "idea"
}

// typeFolder maps a note type to its target vault folder.
func typeFolder(noteType string) string {
	switch noteType {
	case "task":
		return "Tasks"
	case "reference":
		return "References"
	case "idea":
		return "Ideas"
	case "note":
		return "Notes"
	default:
		return "Notes"
	}
}

// triageDestination computes the vault-relative target path for a note.
// Format: {Folder}/{slug}.md
func triageDestination(fromPath, noteType string, parsed *vault.Note) string {
	folder := typeFolder(noteType)
	slug := noteSlug(fromPath, parsed)
	return filepath.Join(folder, slug+".md")
}

// noteSlug derives a filename slug from the note's frontmatter title,
// first H1 heading, or — as a fallback — the original filename.
func noteSlug(fromPath string, parsed *vault.Note) string {
	if title := frontmatterString(parsed.Frontmatter, "title"); title != "" {
		return slugify(title)
	}
	for _, h := range parsed.Headings {
		if h.Level == 1 {
			return slugify(h.Text)
		}
	}
	// Fall back to the timestamp filename (e.g. "20260221-143022").
	return strings.TrimSuffix(filepath.Base(fromPath), ".md")
}

var (
	slugNonAlnum = regexp.MustCompile(`[^a-z0-9]+`)
)

// slugify converts a title to a lowercase, hyphen-separated slug.
func slugify(s string) string {
	s = strings.ToLower(s)
	s = slugNonAlnum.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return "untitled"
	}
	return s
}

// deconflictPath appends a short timestamp suffix to avoid path collisions.
func deconflictPath(path string) string {
	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)
	return fmt.Sprintf("%s-%d%s", base, time.Now().UnixMilli()%100000, ext)
}

// buildTriagedContent rewrites a note's frontmatter with the classified type,
// sets status: processed and triaged date, preserves all other existing fields,
// and appends a ## Related Notes section for any wikilink suggestions.
func buildTriagedContent(parsed *vault.Note, noteType string, linksAdded []string, now time.Time) string {
	var b strings.Builder
	b.WriteString("---\n")

	// Deterministic field order: title → created → type → status → triaged → source → tags → extras.
	if title := frontmatterString(parsed.Frontmatter, "title"); title != "" {
		fmt.Fprintf(&b, "title: %s\n", title)
	}
	if created := frontmatterString(parsed.Frontmatter, "created"); created != "" {
		fmt.Fprintf(&b, "created: %s\n", created)
	}
	fmt.Fprintf(&b, "type: %s\n", noteType)
	b.WriteString("status: processed\n")
	fmt.Fprintf(&b, "triaged: %s\n", now.Format("2006-01-02"))

	if source := frontmatterString(parsed.Frontmatter, "source"); source != "" {
		fmt.Fprintf(&b, "source: %s\n", source)
	}

	// Preserve tags list.
	if tags, ok := parsed.Frontmatter["tags"]; ok {
		switch v := tags.(type) {
		case []string:
			if len(v) > 0 {
				b.WriteString("tags:\n")
				for _, t := range v {
					fmt.Fprintf(&b, "  - %s\n", t)
				}
			}
		case string:
			if v != "" {
				fmt.Fprintf(&b, "tags: %s\n", v)
			}
		}
	}

	b.WriteString("---\n")

	// Append body (ensure leading newline between frontmatter and body).
	body := parsed.Body
	if body != "" && !strings.HasPrefix(body, "\n") {
		b.WriteByte('\n')
	}
	b.WriteString(body)

	// Append suggested wikilinks as a Related Notes section.
	if len(linksAdded) > 0 && !strings.Contains(body, "## Related Notes") {
		if !strings.HasSuffix(b.String(), "\n") {
			b.WriteByte('\n')
		}
		b.WriteString("\n## Related Notes\n")
		for _, link := range linksAdded {
			fmt.Fprintf(&b, "- [[%s]]\n", link)
		}
	}

	return b.String()
}

// enrichSingleNote finds wikilink suggestions for a single note from the index.
// Returns up to 5 suggestions above a 0.7 cosine similarity threshold, sorted
// by similarity descending (highest-similarity notes first).
func enrichSingleNote(store *index.Store, notePath string) []LinkSuggestion {
	notes, err := store.GetAllNoteRows()
	if err != nil {
		return nil
	}
	return enrichSingleNoteFromRows(notes, notePath)
}

// enrichSingleNoteFromRows is the testable core of enrichSingleNote.
func enrichSingleNoteFromRows(notes []index.NoteRow, notePath string) []LinkSuggestion {
	const threshold = 0.7
	const maxLinks = 5

	// Locate the target note's embedding and its existing wikilinks.
	var targetEmb []float32
	existingLinks := make(map[string]bool)
	for _, n := range notes {
		if n.Path != notePath {
			continue
		}
		targetEmb = n.Embedding
		if n.Wikilinks != "" {
			for _, link := range strings.Split(n.Wikilinks, ", ") {
				existingLinks[strings.ToLower(strings.TrimSpace(link))] = true
			}
		}
		break
	}
	if targetEmb == nil {
		return nil // note not yet indexed
	}

	targetName := strings.ToLower(strings.TrimSuffix(filepath.Base(notePath), ".md"))

	// Collect all candidates above the threshold first, then sort.
	var suggestions []LinkSuggestion
	for _, n := range notes {
		if n.Path == notePath || n.Embedding == nil {
			continue
		}

		sim := float64(index.CosineSimilarity(targetEmb, n.Embedding))
		if sim < threshold {
			continue
		}

		toName := strings.ToLower(strings.TrimSuffix(filepath.Base(n.Path), ".md"))
		if existingLinks[toName] || toName == targetName {
			continue
		}

		suggestions = append(suggestions, LinkSuggestion{
			From:       notePath,
			To:         n.Path,
			Similarity: sim,
		})
	}

	// Sort by similarity descending so the highest-scoring notes come first.
	for i := 1; i < len(suggestions); i++ {
		for j := i; j > 0 && suggestions[j].Similarity > suggestions[j-1].Similarity; j-- {
			suggestions[j], suggestions[j-1] = suggestions[j-1], suggestions[j]
		}
	}

	// Truncate to top maxLinks.
	if len(suggestions) > maxLinks {
		suggestions = suggestions[:maxLinks]
	}

	return suggestions
}

// computeAge returns (days old, created-label) from a frontmatter date string or file mod time.
func computeAge(created string, modTime int64, now time.Time) (int, string) {
	if created != "" {
		t, err := time.Parse("2006-01-02", created)
		if err == nil {
			days := int(now.Sub(t).Hours() / 24)
			if days < 0 {
				days = 0
			}
			return days, created
		}
	}
	t := time.Unix(modTime, 0)
	days := int(now.Sub(t).Hours() / 24)
	if days < 0 {
		days = 0
	}
	return days, t.Format("2006-01-02")
}

// formatAgeLabel converts a day count to a human-readable string.
func formatAgeLabel(days int) string {
	switch {
	case days == 0:
		return "today"
	case days == 1:
		return "1 day ago"
	case days < 7:
		return fmt.Sprintf("%d days ago", days)
	case days < 14:
		return "1 week ago"
	case days < 30:
		return fmt.Sprintf("%d weeks ago", days/7)
	case days < 60:
		return "1 month ago"
	default:
		return fmt.Sprintf("%d months ago", days/30)
	}
}

// frontmatterString extracts a string value from a parsed frontmatter map.
func frontmatterString(fm map[string]any, key string) string {
	v, ok := fm[key]
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}

func printTriageListReport(result TriageOutput) {
	header := fmt.Sprintf("Inbox Triage (pending: %d)", result.Summary.Total)
	fmt.Println(header)
	fmt.Println(strings.Repeat("=", len(header)))

	if len(result.Pending) == 0 {
		fmt.Println("\nNo pending notes in inbox.")
		return
	}

	fmt.Println()
	for _, p := range result.Pending {
		noteType := p.Type
		if noteType == "" {
			noteType = "untyped"
		}
		fmt.Printf("  %-40s  %-12s  %s\n",
			filepath.Base(p.Path),
			noteType,
			p.AgeLabel,
		)
	}

	if len(result.Errors) > 0 {
		fmt.Printf("\nErrors (%d):\n", len(result.Errors))
		for _, e := range result.Errors {
			fmt.Printf("  ! %s\n", e)
		}
	}
}

func printTriageAutoReport(result TriageOutput, dryRun bool) {
	header := "Triage"
	if dryRun {
		header += " (dry run)"
	}
	fmt.Println(header)
	fmt.Println(strings.Repeat("=", len(header)))

	if len(result.Processed) > 0 {
		fmt.Println()
		for _, p := range result.Processed {
			prefix := "✓"
			if p.DryRun {
				prefix = "→"
			}
			line := fmt.Sprintf("  %s %s → %s (type: %s",
				prefix,
				filepath.Base(p.FromPath),
				p.ToPath,
				p.NoteType,
			)
			if len(p.LinksAdded) > 0 {
				line += ", links: " + strings.Join(p.LinksAdded, ", ")
			}
			line += ")"
			fmt.Println(line)
		}
	}

	if len(result.Errors) > 0 {
		fmt.Printf("\nErrors (%d):\n", len(result.Errors))
		for _, e := range result.Errors {
			fmt.Printf("  ! %s\n", e)
		}
	}

	fmt.Printf("\nSummary: %d processed, %d errors (of %d total)\n",
		result.Summary.Processed, result.Summary.Errors, result.Summary.Total)
}
