# CLI Boundary Evaluation: Custom obsidian-cli vs Official Obsidian CLI

**Date:** 2026-03-17
**Scope:** Evaluate overlap between the custom `obsidian-cli` (Go, standalone) and the official Obsidian CLI (v1.12.4+, Feb 2026) to determine what to delegate, wrap, retire, or keep.

---

## Architectural Context

The two CLIs operate under fundamentally different models:

| Dimension | Custom CLI | Official CLI (v1.12.4+) |
|---|---|---|
| **Runtime** | Standalone headless binary | Remote control for running Obsidian app |
| **Dependency** | None (pure Go + SQLite) | Requires Obsidian GUI running |
| **Index** | Own SQLite FTS5 + Gemini vectors | Obsidian's internal index (always current) |
| **Link integrity** | No link updating on move/rename | Automatic wikilink updates |
| **Plugin access** | None | Full (Dataview, Templates, Omnisearch, etc.) |
| **Cron / CI / server** | Yes — designed for headless automation | No — auto-launches GUI on first command |
| **Output** | JSON (`--json`) | JSON, CSV, YAML, Markdown, paths, tree, TSV |
| **Platform** | macOS, Linux (arm64/amd64) | Desktop only (macOS, Windows, Linux) |

**Key constraint:** Any capability that must run headless (cron jobs, CI pipelines, server-side automation) **cannot** delegate to the official CLI. This is the primary factor driving keep/wrap decisions below.

---

## Overlap Matrix

### Legend

- **Full** — Official CLI covers the capability completely
- **Partial** — Official CLI covers the base operation but lacks specific features
- **None** — No upstream equivalent

| # | Custom Capability | Official CLI Equivalent | Coverage | Headless Required? |
|---|---|---|---|---|
| 1 | `read` — note with frontmatter, headings, wikilinks | `obsidian read file=` | **Full** | No |
| 2 | `append` — to note end or named section | `obsidian append file= content=` | **Partial** — no section targeting | No |
| 3 | `create` — note with rich custom frontmatter | `obsidian create name= template=` | **Partial** — template-based, no arbitrary frontmatter flags | No |
| 4 | `list` — notes in vault/subdirectory | `obsidian files folder=` | **Full** | No |
| 5 | `capture` — fleeting note to Inbox/ with source | None | **None** | Yes (cron capture) |
| 6 | `search --mode keyword` — FTS5 | `obsidian search query=` | **Full** — native index, richer query syntax (tag/property filters) | No |
| 7 | `search --mode semantic` — Gemini vectors | None | **None** | Yes |
| 8 | `search --mode hybrid` — RRF fusion | None | **None** | Yes |
| 9 | `index` — FTS5 + embedding build | None (upstream index is always live) | **N/A** — different model | Yes |
| 10 | `triage` — classify inbox, move, enrich | None | **None** | Yes (cron triage) |
| 11 | `enrich` — suggest wikilinks + tags via similarity | Partial — `obsidian orphans`, `obsidian backlinks`, `obsidian tags` exist as read-only queries | **Partial** — detection only, no suggestion engine | No |
| 12 | `maintain` — health score, broken links, stale notes | Partial — `obsidian unresolved` for broken links | **Partial** — single check vs composite scoring | No |
| 13 | `sync` (website) — MDX metadata → vault stubs | None (upstream `sync` = Obsidian Sync, entirely different) | **None** | Yes |
| 14 | `ingest` — scout intel + learnings → vault | None | **None** | Yes |
| 15 | `configure` — interactive config setup | Settings UI (not CLI) | **None** | N/A |
| 16 | `doctor` — diagnostics and health check | None | **None** | N/A |

### Upstream-Only Capabilities (not in custom CLI)

| Official CLI Command | What It Does | Relevance |
|---|---|---|
| `prepend` | Insert content at note start | Low — `append --section` covers most use cases |
| `move` | Relocate note with link updating | High — link integrity is valuable |
| `delete` | Trash/permanent delete | Low — rm suffices for headless |
| `daily` | Daily note CRUD (read/append/prepend/open) | Medium — date-aware path resolution |
| `properties` | YAML frontmatter CRUD (set/remove, typed) | High — typed property management |
| `tags` | List, search, rename tags vault-wide | Medium — rename is powerful |
| `links` / `backlinks` | Outgoing/incoming link queries | Medium — graph traversal |
| `tasks` | Task listing and management | Low — not in current workflow |
| `plugins` / `themes` | Plugin/theme management | Low — admin, not note operations |
| `sync` / `publish` / `history` | Obsidian Sync, Publish, version history | Low — requires paid services |
| `eval` | Execute JavaScript in Obsidian runtime | High — escape hatch for anything |
| `folders` | List vault folders as tree | Low — trivial |

---

## Recommendations

### Per-Capability Verdict

| # | Capability | Verdict | Rationale |
|---|---|---|---|
| 1 | `read` | **Wrap** | Delegate to official CLI when Obsidian is running (richer: plugin-rendered content, live index). Fall back to custom for headless/JSON-structured output (frontmatter, headings, wikilinks as parsed fields). Custom's structured parse output is valuable for agent consumption. |
| 2 | `append` | **Keep** | Section-aware append (`--section "## Tasks"`) is a differentiator with no upstream equivalent. The official `append` only targets note end. Keep as-is. |
| 3 | `create` | **Keep** | Custom frontmatter flags (`--type`, `--status`, `--summary`, `--context-set`, `--tags`) map directly to the vault's note taxonomy. Official `create` uses Obsidian templates, which is a different model. Both are useful but serve different workflows. |
| 4 | `list` | **Retire** | Official `files` is strictly superior (more output formats, folder filtering, counts). Custom `list` adds no unique value. Remove after migration. |
| 5 | `capture` | **Keep** | Core workflow with no upstream equivalent. Fleeting-note-to-Inbox pattern, `--source` tagging, cron usage — all unique. |
| 6 | `search --keyword` | **Retire** | Official search uses Obsidian's live index with richer query syntax (tag filters, property filters). Custom FTS5 is a weaker reimplementation that requires manual `index` runs. Retire keyword mode; retain index infrastructure only for embeddings. |
| 7 | `search --semantic` | **Keep** | No upstream equivalent. Gemini embedding vectors for conceptual search is a unique differentiator. |
| 8 | `search --hybrid` | **Wrap** | Keep RRF fusion logic but feed it results from official keyword search (when available) + custom semantic search. Improves keyword leg quality without maintaining a parallel FTS index. |
| 9 | `index` | **Retire partial** | Drop FTS5 indexing entirely — official CLI's live index is always current. Keep embedding index build for semantic search vectors. Rename to `embed` or `index --embeddings-only`. |
| 10 | `triage` | **Keep** | Fully custom: rule-based classification, folder routing, wikilink enrichment, cron automation. No upstream equivalent. Core workflow. |
| 11 | `enrich` | **Wrap** | Use official `orphans` and `backlinks` for detection (better data from live index), keep custom cosine-similarity suggestion engine for wikilink and tag recommendations. |
| 12 | `maintain` | **Wrap** | Use official `unresolved` for broken-link detection. Keep composite health scoring, stale-note detection, and inbox backlog metrics as custom additions. |
| 13 | `sync` (website) | **Keep** | Completely unrelated to Obsidian Sync. Custom MDX-to-vault stub pipeline is bespoke integration. |
| 14 | `ingest` | **Keep** | Scout intel and learnings DB import — entirely custom data pipeline. No overlap. |
| 15 | `configure` | **Keep** | Needed for custom features (Gemini API key, vault path, website path). Different config model from Obsidian settings. |
| 16 | `doctor` | **Keep** | Validates custom infrastructure (embedding index, config permissions, API key). Unrelated to official CLI health. |

### Summary Tally

| Verdict | Count | Capabilities |
|---|---|---|
| **Keep** | 9 | capture, append, create, semantic search, triage, sync, ingest, configure, doctor |
| **Wrap** | 3 | read, hybrid search (feed from upstream keyword), enrich, maintain |
| **Retire** | 2 | list, keyword search (+ FTS5 index) |
| **Delegate** | 0 | — (no full delegation due to headless requirement) |

---

## Phased Migration Plan

### Phase 0: Prerequisite — Official CLI Availability Detection

Add an `upstream.Available() bool` check (probe for `obsidian` binary with version >= 1.12.4 and a running Obsidian instance). All wrap/retire decisions gate on this check. When upstream is unavailable, all current behavior is preserved unchanged.

**Files touched:** `internal/upstream/detect.go` (new)
**Risk:** None — additive only, no behavior change.

---

### Phase 1: Retire `list`, Slim the Index

**Goal:** Remove capabilities fully superseded by upstream.

| Step | Action | Detail |
|---|---|---|
| 1a | Retire `list` command | Remove `internal/cmd/list.go`. Update SKILL.md to direct `list` usage to `obsidian files`. If upstream is unavailable at runtime, print an error with install instructions rather than maintaining dead code. |
| 1b | Drop FTS5 from `index` | Remove FTS5 table creation and keyword indexing from `internal/index/store.go`. Keep embedding storage and vector operations. Rename the command to `obsidian index` with `--embeddings-only` becoming the default (and only) mode. |
| 1c | Retire `search --mode keyword` | Remove keyword search path from `internal/cmd/search.go` and `internal/index/store.go`. When `--mode keyword` is requested, proxy to `obsidian search query=... format=json` if upstream is available, otherwise return an error. |
| 1d | Update `doctor` | Remove FTS5-related health checks. Add upstream CLI availability check. |

**Estimated scope:** ~200 lines removed, ~50 added.
**Breaking changes:** `obsidian list` and `obsidian search --mode keyword` change behavior. Announce in changelog.

---

### Phase 2: Wrap Read, Enrich, Maintain

**Goal:** Use upstream for higher-quality data when available, preserve headless fallback.

| Step | Action | Detail |
|---|---|---|
| 2a | Wrap `read` | When upstream is available, call `obsidian read file=... format=json` for richer content (plugin-rendered, live metadata). Parse upstream JSON into current output schema. Fall back to custom file-based read when headless. |
| 2b | Wrap `enrich` — orphan detection | Replace custom orphan scan with `obsidian orphans format=json` when available. Keep cosine-similarity suggestion engine as-is (no upstream equivalent). |
| 2c | Wrap `maintain` — broken links | Replace custom broken-link scanner with `obsidian unresolved format=json` when available. Keep health score computation, stale-note detection, and inbox metrics. |
| 2d | Wrap `search --mode hybrid` | RRF fusion: keyword leg calls `obsidian search query=... format=json` (upstream), semantic leg uses custom embeddings. Merge with existing RRF logic. Fall back to semantic-only when headless. |

**Files touched:** `internal/upstream/detect.go`, `internal/upstream/client.go` (new — exec wrapper), modifications to `read.go`, `enrich.go`, `maintain.go`, `search.go`.
**Estimated scope:** ~300 lines added/modified.
**Breaking changes:** None — output schemas unchanged, behavior is strictly improved.

---

### Phase 3: Evaluate New Upstream Capabilities

**Goal:** Assess whether to adopt upstream-only features that complement the custom workflow.

| Capability | Assessment | Action |
|---|---|---|
| `properties` (typed set/remove) | High value for triage — could replace custom frontmatter writing with upstream's type-safe property API | **Evaluate** — prototype triage using `properties:set` instead of direct file writes |
| `move` (with link updating) | High value for triage — current triage moves files but doesn't update inbound links | **Adopt** — use `obsidian move` in triage `--auto` when upstream available |
| `daily` (date-aware paths) | Medium value — eliminates manual date path construction in SKILL.md examples | **Evaluate** — may simplify append-to-daily-note workflow |
| `tags:rename` | Medium value — useful for vault maintenance | **Expose** — surface via `maintain --fix` or as standalone pass-through |
| `eval` | Escape hatch — could enable capabilities impossible headless | **Defer** — evaluate if specific use cases arise |

**Timeline:** After Phase 2 is stable. Each adoption is independent and can be prioritized by value.

---

### Phase 4: Upstream Client Extraction (Optional)

**Goal:** If Phase 2-3 patterns prove stable, extract `internal/upstream/` into a reusable Go package for calling the official Obsidian CLI with typed inputs/outputs. This would let other tools in the ecosystem benefit from the same detection + fallback pattern.

**Trigger:** Only if 3+ commands successfully wrap upstream and the pattern is proven stable.

---

## Decision Log

| Decision | Rationale |
|---|---|
| No full delegation to upstream | Headless operation is a hard requirement for cron triage, CI capture, and server-side automation. Official CLI requires a running GUI. |
| Keep semantic/hybrid search | Unique differentiator. Official CLI has no vector search or embedding infrastructure. |
| Keep triage + capture | Core custom workflow (capture → triage → connect) has no upstream equivalent. |
| Retire list over delegate | `list` adds zero value over `files`. Rather than maintaining a wrapper, remove it and direct users to the upstream command. |
| Retire keyword search over wrap | Maintaining a parallel FTS5 index that's always stale relative to upstream's live index creates a worse experience. Clean cut. |
| Wrap over delegate for read | Custom `read` returns structured parsed output (frontmatter dict, headings array, wikilinks array) that upstream's plain-text output doesn't match. Wrapping preserves the schema while improving data quality. |

---

## Sources

- [Obsidian CLI — Official](https://obsidian.md/cli)
- [Obsidian CLI — Help Docs](https://help.obsidian.md/cli)
- [Complete CLI Guide — Frank Anaya](https://frankanaya.com/obsidian-cli/)
- [Obsidian 1.12 CLI Ultimate Guide — WenHaoFree](https://blog.wenhaofree.com/en/posts/articles/obsidian-1-12-cli-ultimate-guide/)
- [Official CLI Announcement — DEV Community](https://dev.to/shimo4228/obsidians-official-cli-is-here-no-more-hacking-your-vault-from-the-back-door-3123)
- [kepano/obsidian-skills SKILL.md](https://github.com/kepano/obsidian-skills/blob/main/skills/obsidian-cli/SKILL.md)
