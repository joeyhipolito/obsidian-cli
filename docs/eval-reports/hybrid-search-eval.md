# Hybrid Search — Quality Evaluation Report

**Date:** 2026-03-18
**Vault:** 117 notes (pro-vault)
**Embedding model:** gemini-embedding-001 (768 dimensions)
**Test queries:** 27 queries × 3 modes = 81 executions
**Raw data:** [`search-query-evaluation-raw.json`](search-query-evaluation-raw.json)
**See also:** [`docs/hybrid-search-implementation.md`](../hybrid-search-implementation.md)

---

## Executive Summary

| Metric | Keyword | Semantic | Hybrid | Overall |
|--------|---------|----------|--------|---------|
| PASS (expected in top-1) | 13 | 17 | 17 | 47/81 |
| PARTIAL (expected in top-3) | 3 | 5 | 4 | 12/81 |
| FAIL | 10 | 5 | 5 | 20/81 |
| ERROR | 1 | 0 | 1 | 2/81 |
| **Accuracy (PASS + PARTIAL)** | **59% (16/27)** | **81% (22/27)** | **78% (21/27)** | **72.8% (59/81)** |

**Key finding:** Semantic mode is the strongest performer (81%). Hybrid mode trails semantic by 3pp because equal-weight RRF lets keyword noise bury correct semantic results. Keyword mode has a hard floor at ~59% due to vocabulary mismatch and crashes on queries with special characters.

**Immediate impact:** Two P0 bugs cause complete crashes (ERRORs on `multi-agent orchestration` in both keyword and hybrid modes). Four P1 issues measurably degrade result quality. All bugs are small, targeted fixes.

---

## 1. Failure Mode Analysis

### 1.1 FTS5 Query Injection (P0 — Crash)

**Affected modes:** keyword, hybrid
**Impact:** 2 ERRORs; cascades to hybrid failure for any special-character query

**Root cause:** User queries are passed directly into FTS5 `MATCH` without quoting or escaping. FTS5 treats unquoted tokens as column names, operators, or invalid syntax.

**Crash query:** `multi-agent orchestration`
The hyphen is interpreted as boolean NOT: `multi AND NOT agent`. Then `orchestration` is parsed as a column name, yielding `SQL logic error: no such column: orchestration`.

**Other at-risk patterns:**

| Pattern | FTS5 Interpretation | Effect |
|---------|---------------------|--------|
| `multi-agent` | `multi NOT agent` | Wrong result set silently |
| `#tag` | Syntax error | ERROR |
| `[[wikilink]]` | Syntax error | ERROR |
| `CI/CD pipeline` | Syntax error on `/` | ERROR |
| `NOT something` | Boolean NOT on all | ERROR |
| `keyword OR` | Dangling operator | ERROR |

FTS5 special characters include: `#  /  [  ]  -  NOT  AND  OR  NEAR  *  ^`

**Fix:** Wrap the user query in double-quotes and escape any internal double-quotes before passing to MATCH:
```go
escaped := `"` + strings.ReplaceAll(query, `"`, `""`) + `"`
```
This forces FTS5 phrase-match mode, which is the correct interpretation of a natural-language search. It removes the ability to use FTS5 boolean operators intentionally — acceptable given the CLI's phrase-first UX contract.

---

### 1.2 Hybrid Error Propagation (P0 — Crash)

**Affected mode:** hybrid
**Impact:** 1 ERROR (same `multi-agent orchestration` query that crashes keyword)

**Root cause:** `SearchHybrid` calls `SearchKeyword` first and returns the error if it fails, killing the entire hybrid call before semantic is ever tried.

```go
// store.go:268 (approximate)
kwResults, err := s.SearchKeyword(query, limit*2)
if err != nil {
    return nil, err  // kills hybrid on any keyword failure
}
```

**Fix:** Capture the keyword error, continue with semantic-only results, and surface it as a warning field rather than a fatal return. This mirrors the existing fallback in `SearchCmd` for missing API keys.

```go
kwResults, kwErr := s.SearchKeyword(query, limit*2)
if kwErr != nil {
    kwResults = nil  // degrade gracefully; semantic still runs
}
```

---

### 1.3 Semantic Score Threshold Too Low (P1 — False Positives)

**Affected modes:** semantic, hybrid
**Impact:** 4 FAILs — all four "should return nothing" queries return 20 false-positive results each

The current threshold is `score > 0`. With 768-dim Gemini embeddings and a 117-note vault whose notes share overlapping vocabulary (technology, productivity, software), cosine similarity never meaningfully approaches 0. All notes score above 0.45 for any query.

**Negative query false positive data:**

| Query | Semantic results | Top score | Threshold 0.55 filters? |
|-------|-----------------|-----------|--------------------------|
| `quantum chromodynamics hadron collider` | 20 | 0.5028 | Yes — removes all 20 |
| `underwater basket weaving certification` | 20 | 0.5269 | Yes — removes all 20 |
| `ancient roman aqueduct engineering` | 20 | 0.5240 | Yes — removes all 20 |
| `cryptocurrency mining profitability 2024` | 20 | 0.5056 | Yes — removes all 20 |

**Score distribution across all 27 queries:**

| Category | Top-1 score range | Notes |
|----------|------------------|-------|
| Correct top-1 results | 0.55–0.73 | All PASS cases |
| Negative query top results | 0.50–0.53 | All 4 false-positive clusters |
| Separation gap | ~0.02 points | Narrow but consistent in eval set |

**Recommended threshold: 0.55.** Eliminates 100% of false positives in the eval set while retaining all 47 PASS results. Threshold should be recalibrated if the vault grows into a new domain or diversifies significantly.

**Note on hybrid:** RRF uses ranks, not scores, so the semantic threshold doesn't automatically protect hybrid. An additional guard is needed: if `semanticResults` is empty after threshold filtering, skip the semantic RRF contribution entirely (rather than polluting the ranking with low-confidence noise).

---

### 1.4 RRF Ranking Regression (P1 — Wrong Top Result)

**Affected mode:** hybrid
**Impact:** At least 1 confirmed case (`daily note template`); likely contributes to hybrid trailing semantic by 3pp

**Mechanism:** With k=60 and equal weights, a note's RRF contribution is:
```
score = 1/(60 + rank_keyword) + 1/(60 + rank_semantic)
```
When keyword mode ranks irrelevant notes with high confidence (incidental term matches against a large body), those notes accumulate more RRF score than a note that ranks #1 in semantic only.

**Documented regression — query `daily note template`:**

| Rank | Keyword result | Keyword score | Semantic result | Semantic score | Hybrid result |
|------|---------------|---------------|-----------------|----------------|---------------|
| 1 | System/Vision & Design Document.md | 10.11 | **System/Templates/Daily Note.md** ✓ | 0.73 | System/Vision & Design Document.md |
| 2 | System/plugin-ideas.md | 9.91 | … | … | System/plugin-ideas.md |
| 5 | — | — | — | — | **System/Templates/Daily Note.md** ✓ |

`Daily Note.md` is the correct answer. It ranks #1 in semantic (score 0.73) but #5 in hybrid because three irrelevant keyword results each contribute `1/61 ≈ 0.0164` to their pool, while `Daily Note.md` has no keyword contribution and its semantic contribution `1/61 ≈ 0.0164` is exactly equal to a single keyword-rank-1 contribution.

**Fix:** Lower k for semantic to increase its per-rank contribution (see §3.3, Option A).

---

### 1.5 Missing Snippets for Semantic-Only Hybrid Results (P2)

**Affected mode:** hybrid
**Impact:** Semantic-only results in hybrid output have an empty `Snippet` field — no body context shown in the result card

FTS5's `snippet()` function generates excerpts for keyword results. Semantic results have no equivalent. The snippet map in `SearchHybrid` is populated only from keyword results:

```go
// Only keyword paths get snippets; semantic-only results get ""
for _, r := range keywordResults {
    snippets[r.Path] = r.Snippet
}
```

**Fix:** For each top-N hybrid result that has no snippet, extract the first N non-empty, non-heading lines of the note's `body` column from the `notes` table. No FTS5 involvement needed for this fallback.

---

### 1.6 FTS5 Score Not Normalized (P2 — Misleading)

**Affected mode:** keyword (cosmetic; no correctness impact on hybrid)

The comment at `store.go:218` says "normalize to 0-1 range." The implementation only negates the FTS5 rank:
```go
r.Score = -r.Score  // FTS5 rank is negative; negation makes it positive
// Actual observed range: 0.25–11.85 (unbounded)
```

BM25 scores (0.25–11.85) and cosine similarity scores (0.45–0.73) are incomparable. This doesn't affect hybrid (which uses ranks via RRF) but would matter for any future attempt to combine scores directly.

---

### 1.7 Brute-Force Cosine Scan (P3 — Scalability)

**Affected modes:** semantic, hybrid
**Impact:** None at 117 notes; latency will degrade at ~1,000+ notes

`SearchSemantic` loads all embeddings into memory and runs a sequential cosine scan. At 117 notes and 768 dimensions this is effectively instant. At 1,000 notes: ~3MB loaded per query, ~768K multiply-accumulate operations. At 10,000 notes: ~30MB, ~7.7M operations.

Mitigation: sqlite-vss or an in-memory HNSW index. Not urgent until the vault exceeds ~2,000 notes.

---

## 2. Improvement Analysis

### 2.1 BM25 Parameter Tuning

FTS5 uses BM25 with hardcoded defaults: k1=1.2, b=0.75. These are reasonable general-purpose values but sub-optimal for a personal knowledge vault.

**Current default behavior:**
- **k1=1.2** — moderate term-frequency saturation. After ~3 occurrences of a term, additional hits contribute diminishing returns. Appropriate for long notes.
- **b=0.75** — strong document-length normalization. Short notes are penalized compared to long notes when both contain the query term equally. This hurts short inbox captures (`Inbox/*.md`) that contain dense, relevant content.

**Option A — Lower b to reduce short-note penalty:**
```sql
-- At index creation time (requires full rebuild)
CREATE VIRTUAL TABLE notes_fts USING fts5(
    ...,
    bm25(k1=1.2, b=0.5)  -- was b=0.75
)
```
Expected effect: Short inbox notes (100–400 chars) rank closer to long project notes when they contain the query term. Predicted gain: 1–2 additional PASS results for queries targeting recent captures.

**Caveat:** FTS5 BM25 parameters are set at index creation time and require a full FTS5 rebuild (`DELETE FROM notes_fts` followed by `INSERT INTO notes_fts SELECT ...`). Measure accuracy before and after on the eval set.

**Option B — Explicit column weights in the rank expression:**
```go
// store.go — add column weights to ORDER BY
ORDER BY notes_fts.rank(1.0, 2.0, 1.5, 1.5, 0.5)
//                     path  title tags  hdgs  body
```
Higher title weight means a note whose title matches the query beats a note that only matches in the body. This would fix several FAIL cases where the correct note's title contains the query term but the body is short.

**Recommendation:** Start with Option B (column weights). It requires no index rebuild, is a one-line change, and directly addresses observable failures where title match is the right signal.

---

### 2.2 Query Expansion

Keyword mode fails on vocabulary mismatch — when the user's words don't appear verbatim in the relevant note.

**Observed failures by category:**

| Query | Expected note | Keyword verdict | Why it fails |
|-------|--------------|-----------------|--------------|
| `how to organize notes` | System/Vision & Design Document.md | 0 results | "organize" not in note; concept handled by semantic |
| `knowledge management` | System/Vision & Design Document.md | Wrong result | Term present but not in the right context |
| `daily note template` | System/Templates/Daily Note.md | Wrong #1 | Correct note has fewer total matches than a longer doc |

**Option A — Zero-result keyword fallback in hybrid (recommended, low effort):**
When `SearchKeyword` returns 0 results, skip keyword's RRF contribution in `SearchHybrid` rather than letting it contribute nothing (which already happens) but also rather than polluting the ranking with weak keyword results. This is a clarification of intent more than a behavior change — the real win is pairing it with higher semantic weight when keyword has no signal.

**Option B — Domain-specific synonym map:**
Maintain a small map (e.g., `"organize" → ["structure", "system"]`) and append synonyms to the FTS5 query as OR terms before dispatch. Low API cost, but requires manual curation and doesn't generalize beyond known synonym pairs.

**Option C — Term expansion from the query embedding:**
For the query embedding, find the top-K vocabulary terms in the indexed notes most similar to the query embedding, then add them to the FTS5 query as OR terms. Requires a term-to-embedding index that doesn't currently exist. High build cost for marginal gain over just using semantic mode directly.

**Recommendation:** Option A paired with weighted RRF (§3.3) covers the gap with minimal code. Options B and C are premature given that semantic mode already handles these cases well.

---

### 2.3 Re-Ranking Strategies

**Problem:** Equal-weight RRF (k=60 for both modes) allows keyword noise to bury correct semantic results (§1.4).

**Option A — Weighted RRF via asymmetric k (recommended):**
Lower k increases contribution from rank 1. Setting k_semantic < k_keyword gives semantic results higher influence:

```go
const kSemantic = 30.0  // was 60 — increases semantic rank-1 contribution
const kKeyword  = 60.0  // unchanged

for i, r := range keywordResults {
    scores[r.Path] += 1.0 / (kKeyword + float64(i+1))
}
for i, r := range semanticResults {
    scores[r.Path] += 1.0 / (kSemantic + float64(i+1))
}
```

At rank 1: semantic contributes `1/31 ≈ 0.0323` vs keyword `1/61 ≈ 0.0164`. The semantic top-1 result now outscores any single-list keyword result, while notes appearing in both lists still accumulate bonus.

**Impact on documented regression (`daily note template`):** `Daily Note.md` (semantic rank 1) scores `0.0323`. Vision & Design Document (keyword rank 1 only) scores `0.0164`. Correct result wins.

**Estimated accuracy gain:** +1–2 PASS results across the 27-query eval set.

**Option B — Score-weighted contribution for semantic:**
Replace semantic's rank-based term with its raw cosine score:
```go
scores[r.Path] += float64(r.Score)  // cosine similarity [0,1]
```
More expressive, but mixes incomparable scales (cosine 0.45–0.73 vs keyword rank contribution 0.015–0.032). Score normalization (§2.4) is a prerequisite. Higher complexity for unclear additional gain over Option A.

**Option C — Cross-encoder re-ranking:**
Run top-K candidates through a Gemini API call to compare query against each candidate document. High accuracy ceiling but adds ~200ms latency and per-query API cost. Not appropriate for interactive CLI use.

**Recommendation:** Option A (asymmetric k: semantic=30, keyword=60) is a 3-line change, no new dependencies, and directly addresses the documented regression.

---

### 2.4 Score Normalization

**Current state:**

| Source | Scale | Range observed | Comparable? |
|--------|-------|----------------|-------------|
| BM25 (keyword) | Unbounded positive | 0.25–11.85 | No |
| Cosine similarity (semantic) | [-1, 1] | 0.45–0.73 | No |
| RRF (hybrid) | Rank-derived | 0.0152–0.0328 | No |

The compressed RRF range (0.016-wide) makes scores unintuitive to end users. More critically, the incomparable scales would block any future score-based fusion (Option B in §2.3) or cross-query ranking.

**Recommended normalization (prerequisite for score-based fusion):**

For BM25, use per-query min-max normalization over the result set:
```
norm_score = (score - min_score) / (max_score - min_score + ε)
```
Maps each result set to [0, 1] where 1 = best match in this query. Applied post-retrieval, not stored.

For semantic, shift the baseline to account for the vault-specific floor:
```
effective_score = score - 0.45  // vault floor from eval data
```
Maps [0.45, 0.73] → [0.00, 0.28]. The 0.55 threshold becomes `effective_score > 0.10` — a cleaner signal for "above the noise floor."

**User-facing score:** After normalization, map to a 0–100 scale with human-readable buckets:
- 80–100: strong match
- 50–79: related
- <50: weak signal

**Priority:** Low for correctness; medium for future-proofing and user trust. Implement after P0/P1 fixes.

---

## 3. Prioritized Recommendations

### Tier 1 — Fix Before Shipping (P0)

| # | Issue | Fix | Effort | Impact |
|---|-------|-----|--------|--------|
| 1 | FTS5 query injection crashes search | Quote user query: `` ` `"` + strings.ReplaceAll(query, `"`, `""`) + `"` `` | 2 lines | Eliminates 2 ERRORs; fixes all special-char queries |
| 2 | Hybrid propagates keyword crash | Catch keyword error in `SearchHybrid`, continue with semantic-only | 5 lines | Eliminates 1 ERROR; hybrid degrades gracefully |

### Tier 2 — Quality Wins (P1)

| # | Issue | Fix | Effort | Impact |
|---|-------|-----|--------|--------|
| 3 | No semantic threshold — nonsense queries return 20 results | `score > 0.55` in `SearchSemantic` | 1 line | +4 correct-empty results (100% of negative test cases) |
| 4 | RRF buries correct semantic result | Lower k_semantic from 60 → 30 in `SearchHybrid` | 3 lines | Fixes `daily note template` regression; est. +1–2 PASS |
| 5 | Hybrid false positives on negative queries | Skip semantic RRF contribution when all semantic results filtered by threshold | 5 lines | Prevents hybrid inheriting semantic false positives |

### Tier 3 — UX and Correctness Polish (P2)

| # | Issue | Fix | Effort | Impact |
|---|-------|-----|--------|--------|
| 6 | Semantic-only hybrid results have empty snippet | Extract first N lines of `notes.body` for results without a snippet | ~20 lines | Non-empty result cards for semantic finds |
| 7 | BM25 column weighting | Add `rank(1.0, 2.0, 1.5, 1.5, 0.5)` column weights to ORDER BY | 1 line | +1–2 PASS for keyword mode on title-match queries |
| 8 | Zero-result keyword fallback | When keyword returns 0 results, omit keyword RRF term | 3 lines | Prevents empty hybrid when semantic has results |

### Tier 4 — Future Work (P3)

| # | Issue | Fix | Notes |
|---|-------|-----|-------|
| 9 | Brute-force cosine scan | ANN index (sqlite-vss or HNSW) | Needed at ~1,000+ notes |
| 10 | BM25 b-parameter tuning | Rebuild FTS5 index with b=0.5 | Requires eval re-run after rebuild |
| 11 | Score normalization for display | Per-query min-max BM25; shift semantic baseline | Prerequisite for score-based RRF (§2.3, Option B) |

---

## 4. Projected Accuracy After Tier 1 + 2 Fixes

**Baseline:** 72.8% (59/81)

| Fix | Current result | After fix | Delta |
|-----|---------------|-----------|-------|
| FTS5 escape (P0) | 2 ERROR → FAIL | 2 PASS (multi-agent orchestration now finds correct result) | +2 |
| Hybrid error propagation (P0) | 1 ERROR | 1 PASS | +1 |
| Semantic threshold 0.55 (P1) | 4 FAIL (neg. queries all modes) | 4 PASS | +4 |
| Weighted RRF k_sem=30 (P1) | 1–2 FAIL | 1–2 PASS | +1–2 |

**Projected accuracy after Tier 1+2:** ~(67–68)/81 = **83–84%**

This estimate is conservative. Fixing FTS5 injection may also rescue several additional FAIL/PARTIAL cases in keyword mode where queries currently fail silently due to syntax misinterpretation rather than genuine vocabulary mismatch.

---

## 5. Score Calibration Reference

### Keyword (BM25 via FTS5 `rank`)

- Range: 0.25–11.85 (unbounded; not normalized despite code comment at `store.go:218`)
- Higher = more term matches weighted by frequency and document length
- Not comparable across queries
- k1=1.2, b=0.75 (FTS5 defaults)

### Semantic (cosine similarity)

- Range: 0.45–0.73 (observed in this eval; theoretical bounds: -1 to 1)
- Effective floor ~0.45 for a coherent, single-domain vault — all notes share vocabulary
- Score spread within a single query: 0.05–0.12
- Dynamic range too narrow for confident thresholding at `> 0`; requires `> 0.55` to separate signal from noise

```
Query                                       Top    Bot    Spread
kubernetes                                  0.555  0.516  0.040
obsidian                                    0.623  0.533  0.090
golang                                      0.602  0.507  0.095
docker                                      0.588  0.526  0.062
daily note template                         0.730  0.605  0.125
how to organize notes effectively           0.626  0.548  0.077
quantum chromodynamics hadron collider      0.503  0.473  0.030  ← false positive cluster
underwater basket weaving certification     0.527  0.458  0.069  ← false positive cluster
ancient roman aqueduct engineering          0.524  0.461  0.063  ← false positive cluster
cryptocurrency mining profitability 2024    0.506  0.466  0.040  ← false positive cluster
```

### Hybrid (RRF)

- Range: 0.0152–0.0328 (all results compressed into a 0.018-wide band)
- Score reflects rank position only, not match quality
- Unintuitive for users; normalization to 0–100 recommended (§2.4)

---

## 6. Test Methodology

- All 27 queries run against a live 117-note Obsidian vault in three modes (81 total executions)
- Results captured as JSON in [`search-query-evaluation-raw.json`](search-query-evaluation-raw.json)
- Verdict criteria:
  - **PASS** — expected note in top-1
  - **PARTIAL** — expected note in top-2 or top-3
  - **FAIL** — expected note not in top-3, or expected empty result but results returned
  - **ERROR** — search crashed
- Tests run 2026-03-18 after fresh re-index (117 notes, 117 embeddings, 0 index errors)

**Query categories:**
- Exact keyword match (6 queries): single concrete terms that should FTS-match
- Semantic/conceptual (8 queries): natural-language questions and concepts
- Multi-word phrase (9 queries): two-to-four word technical phrases
- Negative (4 queries): out-of-domain queries that should return empty
