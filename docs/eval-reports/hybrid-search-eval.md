# Hybrid Search Evaluation Report

**Date:** 2026-03-18
**Vault:** 117 notes (pro-vault)
**Embedding model:** gemini-embedding-001 (768 dimensions)
**Test queries:** 23 queries × 3 modes = 69 query executions

---

## 1. Executive Summary

The hybrid search system works well for simple single-keyword queries but has **three critical bugs** and several quality issues that degrade the experience for real-world queries:

1. **FTS5 query injection** — Special characters (`/`, `#`, `[`, `NOT`, etc.) crash keyword search and cascade into hybrid mode, returning zero results even when semantic search would succeed.
2. **No semantic score threshold** — Nonsense queries like "xyzzy12345" return 20 results with high confidence scores (0.62), creating false positives.
3. **Hybrid ranking regression** — For conceptual queries, RRF fusion buries the correct semantic result under noisy keyword matches.

---

## 2. Test Query Results

### 2.1 Exact Keyword Matches

| Query | Mode | Results | Top Score | Top Result | Verdict |
|-------|------|---------|-----------|------------|---------|
| kubernetes | keyword | 1 | 4.94 | People/Frank Li.md | OK |
| kubernetes | semantic | 20 | 0.56 | Inbox/20260308-212511.md | OK |
| kubernetes | hybrid | 20 | 0.016 | Inbox/20260308-212511.md | OK |
| obsidian | keyword | 10 | 4.78 | Projects/Via/Obsidian CLI.md | OK |
| obsidian | semantic | 20 | 0.62 | Projects/Via/Obsidian CLI.md | OK |
| obsidian | hybrid | 20 | 0.033 | Projects/Via/Obsidian CLI.md | OK — both modes agree |
| docker | keyword | 4 | 6.62 | Blog/early-docker-adoption.md | OK |
| docker | semantic | 20 | 0.59 | Blog/early-docker-adoption.md | OK |
| docker | hybrid | 20 | 0.033 | Blog/early-docker-adoption.md | OK — both modes agree |

**Finding:** Simple keyword queries work well. When both modes agree, hybrid correctly ranks the result highest.

### 2.2 Semantic/Conceptual Queries

| Query | Mode | Results | Top Score | Top Result | Verdict |
|-------|------|---------|-----------|------------|---------|
| how to organize my notes effectively | keyword | 0 | — | — | Expected: FTS5 can't match this |
| how to organize my notes effectively | semantic | 20 | 0.63 | Note Writing Rules & Conventions.md | **Good** |
| how to organize my notes effectively | hybrid | 20 | 0.016 | Note Writing Rules & Conventions.md | OK — semantic fills the gap |
| daily note template | keyword | 4 | 10.11 | Vision & Design Document.md | **Wrong** — should be Daily Note.md |
| daily note template | semantic | 20 | 0.73 | Templates/Daily Note.md | **Correct** |
| daily note template | hybrid | 20 | 0.030 | Vision & Design Document.md | **WRONG** — keyword noise buries correct answer to #5 |
| mental model for decision making | keyword | 0 | — | — | Expected |
| mental model for decision making | semantic | 20 | 0.62 | Templates/Decision Record.md | Good |

**Finding:** Semantic search handles conceptual queries well. But hybrid fails when FTS5 returns confident-but-wrong results that outrank the correct semantic match.

### 2.3 Should-Return-Nothing Queries

| Query | Mode | Results | Top Score | Verdict |
|-------|------|---------|-----------|---------|
| quantum chromodynamics hadron collider | keyword | 0 | — | **Correct** |
| quantum chromodynamics hadron collider | semantic | 20 | 0.50 | **False positive** — no relevant notes exist |
| underwater basket weaving | semantic | 20 | 0.52 | **False positive** |
| xyzzy12345 nonexistent term | semantic | 20 | 0.62 | **False positive** — score HIGHER than valid queries |

**Finding:** Semantic search has no ability to say "nothing matches." Cosine similarity with 768-dimensional Gemini vectors clusters around 0.45-0.55 for unrelated content, so the `> 0` threshold catches nothing.

### 2.4 Edge Cases / Syntax Errors

| Query | Mode | Error? | Verdict |
|-------|------|--------|---------|
| `#tag` | keyword | **YES** — `fts5: syntax error near "#"` | **BUG** |
| `#tag` | hybrid | **YES** — error propagates | **BUG** — should fall back to semantic |
| `[[wikilink]]` | keyword | **YES** — FTS5 syntax error | **BUG** |
| `[[wikilink]]` | hybrid | **YES** — error propagates | **BUG** |
| `CI/CD pipeline` | keyword | **YES** — `syntax error near "/"` | **BUG** |
| `CI/CD pipeline` | hybrid | **YES** — error propagates | **BUG** |
| `NOT something` | keyword | **YES** — FTS5 parses as boolean | **BUG** |
| `the` | keyword | 20 | 0.25 | OK but low utility |

**Finding:** Any query containing FTS5 special characters (`#`, `/`, `[`, `]`, `NOT`, `AND`, `OR`, `NEAR`, `*`, `^`) crashes both keyword and hybrid modes.

---

## 3. Failure Modes Identified

### 3.1 [CRITICAL] FTS5 Query Injection

**Location:** `store.go:198-223` — `SearchKeyword()`

The query string is passed directly to FTS5 `MATCH` without sanitization. FTS5 has its own query syntax where characters like `#`, `/`, `[`, `]` and keywords like `NOT`, `AND`, `OR`, `NEAR` are operators.

**Impact:** Any user query containing common punctuation crashes the search. This is especially bad for Obsidian users who naturally search for `#tags` and `[[wikilinks]]`.

**Fix:** Quote the query as an FTS5 phrase using `"` wrapping, or escape special characters. Example:
```go
escaped := `"` + strings.ReplaceAll(query, `"`, `""`) + `"`
```

### 3.2 [CRITICAL] Hybrid Error Propagation

**Location:** `store.go:269-275` — `SearchHybrid()`

```go
keywordResults, err := s.SearchKeyword(query, limit*2)
if err != nil {
    return nil, err  // ← kills hybrid, never tries semantic
}
```

When FTS5 fails (syntax error), hybrid mode returns an error instead of falling back to semantic-only results. The `SearchCmd` handler in `search.go:82-84` also fails to catch this — it only handles embedding failures, not keyword failures.

**Fix:** Catch keyword errors gracefully and continue with semantic-only results:
```go
keywordResults, kwErr := s.SearchKeyword(query, limit*2)
if kwErr != nil {
    keywordResults = nil // degrade gracefully
}
```

### 3.3 [HIGH] No Semantic Score Threshold

**Location:** `store.go:247` — `if score > 0`

The threshold of `> 0` is effectively no threshold. With 768-dimensional Gemini embeddings, cosine similarity for unrelated content clusters around 0.45-0.55. The query "xyzzy12345 nonexistent term" gets a top score of 0.62, which is *higher* than legitimate queries like "kubernetes" (0.55).

**Score distribution observed:**
- Legitimate queries: top scores 0.55-0.73, bottom scores 0.50-0.60
- Nonsense queries: top scores 0.50-0.62, bottom scores 0.47-0.52
- Spread within a query: typically 0.05-0.12

**Fix options (ranked by complexity):**
1. **Static threshold ~0.55** — Simple, eliminates most false positives but may cut valid results for broad queries
2. **Adaptive threshold** — Use mean + 1 standard deviation of all scores; only return results above this
3. **Top-k with gap detection** — Return results only until there's a large drop in score

### 3.4 [HIGH] RRF Ranking Regression for Conceptual Queries

**Example:** Query "daily note template"
- Semantic #1: `Templates/Daily Note.md` (score 0.73) — **correct**
- Keyword #1: `Vision & Design Document.md` (score 10.11) — **wrong** (contains words separately)
- Hybrid #1: `Vision & Design Document.md` — **wrong**, correct answer pushed to #5

RRF treats keyword and semantic rankings equally (each contributes `1/(60+rank)`). When FTS5 returns 4 results and the correct answer isn't among them, those 4 keyword results all outrank every semantic-only result because they get scores from both lists while semantic-only results get scores from one.

**Fix options:**
1. **Weighted RRF** — Give semantic results higher weight (e.g., 0.6 semantic, 0.4 keyword)
2. **Boost results appearing in both lists** — Multiplicative bonus for overlap
3. **Cascade strategy** — If keyword results ≤ 5, pad with top semantic results instead of pure RRF

### 3.5 [MEDIUM] FTS5 Score Normalization Is Misleading

**Location:** `store.go:219` — `r.Score = -r.Score`

The comment says "normalize to 0-1 range" but negating the rank does not normalize to [0,1]. Observed keyword scores range from 0.25 to 11.85. This doesn't affect hybrid search (RRF is rank-based) but makes standalone keyword scores confusing for users.

### 3.6 [MEDIUM] Missing Snippets in Semantic Results

Semantic search results have empty snippet fields. In hybrid mode, only results that also appeared in keyword search get snippets. This degrades the user experience for semantic-only matches.

**Fix:** After hybrid ranking, fetch snippets for the top results by querying FTS5 for snippet generation, or generate snippets from the body text directly.

### 3.7 [LOW] Brute-Force Semantic Search

**Location:** `store.go:227` — `SearchSemantic()` loads ALL embeddings into memory.

For 117 notes this is fast. At 1,000+ notes it will become noticeable. At 10,000+ notes it will be problematic.

**Future fix:** Use an approximate nearest neighbor index (e.g., sqlite-vss, or in-memory HNSW).

---

## 4. Score Calibration Analysis

### FTS5 (keyword) scores
- Range: 0.25 to 11.85 (unbounded positive)
- Higher = more term matches and frequency
- Not comparable across queries

### Cosine similarity (semantic) scores
- Range: 0.45 to 0.73 (bounded -1 to 1, but observed floor ~0.45)
- Spread within single query: 0.05-0.12
- Not enough dynamic range for confident thresholding

### RRF (hybrid) scores
- Range: 0.0152 to 0.0328
- Tiny absolute values, unintuitive for users
- Score only reflects rank position, not match quality

**Recommendation:** For user-facing scores, consider normalizing to a 0-100 scale with semantic meaning (e.g., >80 = strong match, 50-80 = related, <50 = weak).

---

## 5. Prioritized Recommendations

### P0 — Must Fix (bugs)

1. **FTS5 query escaping** — Quote user queries to prevent syntax errors. Wrap in double quotes and escape embedded quotes. ~10 lines of code.
2. **Hybrid graceful degradation** — When keyword search errors, fall back to semantic-only results instead of returning an error. ~5 lines of code.

### P1 — High Impact

3. **Semantic score threshold** — Add a minimum cosine similarity cutoff (start with 0.55, tune empirically). Prevents false positives for nonsense queries. ~3 lines of code.
4. **Weighted RRF or boosted overlap** — Give semantic ranking more weight than keyword ranking in hybrid mode, or apply a multiplier for results appearing in both lists. ~15 lines of code.

### P2 — Quality of Life

5. **Snippet generation for semantic results** — Extract relevant body text around query terms for semantic-only results. ~30 lines of code.
6. **Score normalization for display** — Normalize all scores to a consistent 0-100 scale for user-facing output. ~20 lines of code.
7. **BM25 parameter tuning** — FTS5 supports `bm25()` ranking function with configurable k1 and b parameters. Default (k1=1.2, b=0.75) may overweight long documents. Worth experimenting with k1=1.5, b=0.5.

### P3 — Future Scalability

8. **Approximate nearest neighbor index** — Replace brute-force cosine scan with sqlite-vss or in-memory HNSW for vaults with 1,000+ notes.
9. **Query expansion** — For keyword search, add stemming or synonym expansion to bridge the term mismatch gap.
10. **Cross-encoder re-ranking** — For the final top-20, use an LLM or cross-encoder to re-rank by true relevance. High latency cost but maximum quality.

---

## 6. BM25 Parameter Analysis

FTS5 uses the `rank` column which applies BM25 internally with default parameters:
- **k1 = 1.2** (term frequency saturation)
- **b = 0.75** (document length normalization)

The `b=0.75` default means longer documents are penalized, but FTS5 searches across multiple columns (path, title, tags, headings, body). Long notes with many term matches in the body may still dominate.

**Experiment suggestion:** Use explicit `bm25()` ranking function with per-column weights:
```sql
SELECT *, bm25(notes_fts, 0, 10.0, 5.0, 3.0, 1.0) as score
FROM notes_fts WHERE notes_fts MATCH ?
ORDER BY score
```
This would weight: path(0), title(10), tags(5), headings(3), body(1) — making title matches much more important than body matches.

---

## 7. Test Methodology

- All queries run against a live 117-note Obsidian vault
- Each query tested in keyword, semantic, and hybrid modes (69 total executions)
- Results captured as JSON and analyzed programmatically
- Correctness judged by whether the most intuitively relevant note appears in top-3
- Tests run on 2026-03-18 after fresh re-index (9 notes updated, 0 errors)

## 8. Appendix: Raw Score Distributions

### Semantic scores by query (top / bottom / spread)

```
Query                                          Top    Bot    Spread
kubernetes                                     0.5554 0.5159 0.0396
obsidian                                       0.6228 0.5330 0.0898
golang                                         0.6022 0.5065 0.0957
docker                                         0.5877 0.5262 0.0615
how to organize my notes effectively           0.6255 0.5482 0.0773
daily note template                            0.7295 0.6047 0.1248
quantum chromodynamics hadron collider         0.5028 0.4733 0.0295
underwater basket weaving                      0.5161 0.4575 0.0586
xyzzy12345 nonexistent term                    0.6216 0.5199 0.1017
```
