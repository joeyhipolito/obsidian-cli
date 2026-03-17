# Search Evaluation Results

**Date:** 2026-03-18
**Vault:** 117 notes (pro-vault)
**Embedding model:** gemini-embedding-001 (768 dimensions)
**Test queries:** 27 queries × 3 modes = 81 executions

---

## Summary

| Metric | Count |
|--------|-------|
| Total executions | 81 |
| PASS (expected in top-1) | 47 |
| PARTIAL (expected in top-3) | 12 |
| FAIL (expected not in top-3 or false positive) | 20 |
| ERROR (search crashed) | 2 |
| **Accuracy (PASS+PARTIAL)** | **59/81 (72.8%)** |

---

### Keyword Mode: 13 pass, 3 partial, 10 fail, 1 error

### Semantic Mode: 17 pass, 5 partial, 5 fail, 0 error

### Hybrid Mode: 17 pass, 4 partial, 5 fail, 1 error

---

## Exact Keyword Match

### Query: `kubernetes`
**Expected:** People/Frank Li.md, Inbox/20260308-212511.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 1 | 4.9396 | `People/Frank Li.md` | #1 | PASS (rank 1) |
| semantic | 20 | 0.5554 | `Inbox/20260308-212511.md` | #1 | PASS (rank 1) |
| hybrid | 20 | 0.0164 | `People/Frank Li.md` | #1 | PASS (rank 1) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

1. `People/Frank Li.md` — Engineering Manager (B2B Team) (score: 4.9396) **[EXPECTED]**

**semantic:**

1. `Inbox/20260308-212511.md` — 20260308-212511 (score: 0.5554) **[EXPECTED]**
2. `Projects/Via/ContentKit.md` — ContentKit (score: 0.5439)
3. `System/Implementation Plan.md` — Implementation Plan (score: 0.5368)
4. `Projects/sched-fyi-analysis.md` — sched.fyi — Platform Analysis (score: 0.5360)
5. `System/Vision & Design Document.md` — Vision & Design Document (score: 0.5326)

**hybrid:**

1. `People/Frank Li.md` — Engineering Manager (B2B Team) (score: 0.0164) **[EXPECTED]**
2. `Inbox/20260308-212511.md` — 20260308-212511 (score: 0.0164) **[EXPECTED]**
3. `Projects/Via/ContentKit.md` — ContentKit (score: 0.0161)
4. `System/Implementation Plan.md` — Implementation Plan (score: 0.0159)
5. `Projects/sched-fyi-analysis.md` — sched.fyi — Platform Analysis (score: 0.0156)

</details>

### Query: `docker`
**Expected:** Projects/Website/Blog/2022-02-10-early-docker-adoption-2014.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 4 | 6.6172 | `...e/Blog/2022-02-10-early-docker-adoption-2014.md` | #1 | PASS (rank 1) |
| semantic | 20 | 0.5877 | `...e/Blog/2022-02-10-early-docker-adoption-2014.md` | #1 | PASS (rank 1) |
| hybrid | 20 | 0.0328 | `...e/Blog/2022-02-10-early-docker-adoption-2014.md` | #1 | PASS (rank 1) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

1. `Projects/Website/Blog/2022-02-10-early-docker-adoption-2014.md` — Docker in 2014: What Version 0.8 Taught Me That Polished Tools Never Could (score: 6.6172) **[EXPECTED]**
2. `Projects/Website/Stories/oxygen-ventures-2015.md` — The Patterns I Found at My First Real Job (score: 6.2120)
3. `Projects/Via/Substack Backfill Plan.md` — Substack Backfill Plan (score: 5.4897)
4. `Projects/Website/Stories/custom-developer-experience-tools.md` — Three Tools I Built Because Someone on the Team Was Stuck (score: 4.4980)

**semantic:**

1. `Projects/Website/Blog/2022-02-10-early-docker-adoption-2014.md` — Docker in 2014: What Version 0.8 Taught Me That Polished Tools Never Could (score: 0.5877) **[EXPECTED]**
2. `Projects/Website/Stories/oxygen-ventures-2015.md` — The Patterns I Found at My First Real Job (score: 0.5691)
3. `Projects/Via/Substack Backfill Plan.md` — Substack Backfill Plan (score: 0.5580)
4. `Inbox/20260308-212511.md` — 20260308-212511 (score: 0.5562)
5. `Projects/Via/Gastown Analysis.md` — Gastown Analysis (score: 0.5542)

**hybrid:**

1. `Projects/Website/Blog/2022-02-10-early-docker-adoption-2014.md` — Docker in 2014: What Version 0.8 Taught Me That Polished Tools Never Could (score: 0.0328) **[EXPECTED]**
2. `Projects/Website/Stories/oxygen-ventures-2015.md` — The Patterns I Found at My First Real Job (score: 0.0323)
3. `Projects/Via/Substack Backfill Plan.md` — Substack Backfill Plan (score: 0.0317)
4. `Inbox/20260308-212511.md` — 20260308-212511 (score: 0.0156)
5. `Projects/Website/Stories/custom-developer-experience-tools.md` — Three Tools I Built Because Someone on the Team Was Stuck (score: 0.0156)

</details>

### Query: `obsidian`
**Expected:** Projects/Via/Obsidian CLI.md, System/Vision & Design Document.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 10 | 4.7834 | `Projects/Via/Obsidian CLI.md` | #1 | PASS (rank 1) |
| semantic | 20 | 0.6228 | `Projects/Via/Obsidian CLI.md` | #1 | PASS (rank 1) |
| hybrid | 20 | 0.0328 | `Projects/Via/Obsidian CLI.md` | #1 | PASS (rank 1) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

1. `Projects/Via/Obsidian CLI.md` — Obsidian CLI — Project Hub (score: 4.7834) **[EXPECTED]**
2. `Projects/Website/Blog/2026-01-12-lifeos-personal-operating-system.md` — LifeOS: Building an AI-Powered Personal Operating System with Claude Code & Obsidian (score: 4.6644)
3. `Projects/Website/Blog/2026-01-31-six-plugins-one-brain.md` — Six Plugins, One Brain (score: 4.3868)
4. `System/Implementation Plan.md` — Implementation Plan (score: 4.0088)
5. `System/Note Writing Rules & Conventions.md` — Note Writing Rules & Conventions (score: 3.8969)

**semantic:**

1. `Projects/Via/Obsidian CLI.md` — Obsidian CLI — Project Hub (score: 0.6228) **[EXPECTED]**
2. `System/Implementation Plan.md` — Implementation Plan (score: 0.5928)
3. `System/plugin-ideas.md` — Plugin & Skill Proposals (score: 0.5845)
4. `System/Vision & Design Document.md` — Vision & Design Document (score: 0.5690) **[EXPECTED]**
5. `Projects/Website/Blog/2026-01-12-lifeos-personal-operating-system.md` — LifeOS: Building an AI-Powered Personal Operating System with Claude Code & Obsidian (score: 0.5653)

**hybrid:**

1. `Projects/Via/Obsidian CLI.md` — Obsidian CLI — Project Hub (score: 0.0328) **[EXPECTED]**
2. `System/Implementation Plan.md` — Implementation Plan (score: 0.0318)
3. `Projects/Website/Blog/2026-01-12-lifeos-personal-operating-system.md` — LifeOS: Building an AI-Powered Personal Operating System with Claude Code & Obsidian (score: 0.0315)
4. `Projects/Website/Blog/2026-01-31-six-plugins-one-brain.md` — Six Plugins, One Brain (score: 0.0306)
5. `System/plugin-ideas.md` — Plugin & Skill Proposals (score: 0.0304)

</details>

### Query: `golang`
**Expected:** Projects/Website/Blog/2026-02-12-building-ai-intelligence-layer-pure-go.md, Projects/Via/Scout.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 4 | 5.6816 | `...26-02-03-how-multi-agent-orchestration-works.md` | #2 | PARTIAL (rank 2) |
| semantic | 20 | 0.6022 | `...02-12-building-ai-intelligence-layer-pure-go.md` | #1 | PASS (rank 1) |
| hybrid | 20 | 0.0325 | `...02-12-building-ai-intelligence-layer-pure-go.md` | #1 | PASS (rank 1) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

1. `Projects/Website/Blog/2026-02-03-how-multi-agent-orchestration-works.md` — How Multi-Agent Orchestration Works (score: 5.6816)
2. `Projects/Website/Blog/2026-02-12-building-ai-intelligence-layer-pure-go.md` — Building an AI Intelligence Layer in Pure Go (score: 5.5749) **[EXPECTED]**
3. `People/Frank Li.md` — Engineering Manager (B2B Team) (score: 4.8331)
4. `Projects/ecal-security-deep-dive.md` — ECAL Security Deep Dive — Popup Evasion, Jailbreaks, API Signing (score: 1.0600)

**semantic:**

1. `Projects/Website/Blog/2026-02-12-building-ai-intelligence-layer-pure-go.md` — Building an AI Intelligence Layer in Pure Go (score: 0.6022) **[EXPECTED]**
2. `Inbox/20260308-212511.md` — 20260308-212511 (score: 0.5915)
3. `Projects/Website/Blog/2026-02-03-how-multi-agent-orchestration-works.md` — How Multi-Agent Orchestration Works (score: 0.5556)
4. `Projects/Via/ContentKit.md` — ContentKit (score: 0.5382)
5. `Inbox/20260226-144631.md` — 20260226-144631 (score: 0.5338)

**hybrid:**

1. `Projects/Website/Blog/2026-02-12-building-ai-intelligence-layer-pure-go.md` — Building an AI Intelligence Layer in Pure Go (score: 0.0325) **[EXPECTED]**
2. `Projects/Website/Blog/2026-02-03-how-multi-agent-orchestration-works.md` — How Multi-Agent Orchestration Works (score: 0.0323)
3. `People/Frank Li.md` — Engineering Manager (B2B Team) (score: 0.0300)
4. `Projects/ecal-security-deep-dive.md` — ECAL Security Deep Dive — Popup Evasion, Jailbreaks, API Signing (score: 0.0262)
5. `Inbox/20260308-212511.md` — 20260308-212511 (score: 0.0161)

</details>

### Query: `halloween`
**Expected:** Family/Milestones/2025-10-31 — Hiide First Halloween.md, Family/Milestones/2025-10-31 — Kaede First Halloween.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 6 | 5.5982 | `...ilestones/2025-10-31 — Hiide First Halloween.md` | #1 | PASS (rank 1) |
| semantic | 20 | 0.5979 | `Daily/2025/10/2025-10-31.md` | #2 | PARTIAL (rank 2) |
| hybrid | 20 | 0.0325 | `Daily/2025/10/2025-10-31.md` | #2 | PARTIAL (rank 2) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

1. `Family/Milestones/2025-10-31 — Hiide First Halloween.md` — 2025-10-31 — Hiide First Halloween (score: 5.5982) **[EXPECTED]**
2. `Daily/2025/10/2025-10-31.md` — Friday, Oct 31, 2025 — Daily (score: 5.5795)
3. `Family/Milestones/2025-10-31 — Kaede First Halloween.md` — 2025-10-31 — Kaede First Halloween (score: 5.5022) **[EXPECTED]**
4. `People/Kaede Miyuki Hipolito.md` — Kaede Miyuki Hipolito (score: 4.7958)
5. `People/Hiide Illumi Hipolito.md` — Hiide Illumi Hipolito (score: 4.3671)

**semantic:**

1. `Daily/2025/10/2025-10-31.md` — Friday, Oct 31, 2025 — Daily (score: 0.5979)
2. `Family/Milestones/2025-10-31 — Kaede First Halloween.md` — 2025-10-31 — Kaede First Halloween (score: 0.5859) **[EXPECTED]**
3. `Family/Milestones/2025-10-31 — Hiide First Halloween.md` — 2025-10-31 — Hiide First Halloween (score: 0.5745) **[EXPECTED]**
4. `Inbox/20260308-212511.md` — 20260308-212511 (score: 0.5565)
5. `Writing/the-vampiric-effect-of-ai.md` — The Vampiric Effect of AI (score: 0.5329)

**hybrid:**

1. `Daily/2025/10/2025-10-31.md` — Friday, Oct 31, 2025 — Daily (score: 0.0325)
2. `Family/Milestones/2025-10-31 — Hiide First Halloween.md` — 2025-10-31 — Hiide First Halloween (score: 0.0323) **[EXPECTED]**
3. `Family/Milestones/2025-10-31 — Kaede First Halloween.md` — 2025-10-31 — Kaede First Halloween (score: 0.0320) **[EXPECTED]**
4. `People/Kaede Miyuki Hipolito.md` — Kaede Miyuki Hipolito (score: 0.0269)
5. `People/Hiide Illumi Hipolito.md` — Hiide Illumi Hipolito (score: 0.0259)

</details>

### Query: `gratitude`
**Expected:** Gratitudes/Gratitude log.md, Gratitudes/2025-10-30 Birthday Celebration with Camachons.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 20 | 1.0089 | `Gratitudes/Gratitude log.md` | #1 | PASS (rank 1) |
| semantic | 20 | 0.6532 | `Gratitudes/Gratitude log.md` | #1 | PASS (rank 1) |
| hybrid | 20 | 0.0328 | `Gratitudes/Gratitude log.md` | #1 | PASS (rank 1) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

1. `Gratitudes/Gratitude log.md` — Gratitude Practice (score: 1.0089) **[EXPECTED]**
2. `People/Michaela.md` — Real Estate Agent (score: 0.9896)
3. `Gratitudes/2025-10-03 Gratitude for Ariel Camacho.md` — 2025-10-03 Gratitude for Ariel Camacho (score: 0.9836)
4. `Gratitudes/2025-10-03 Gratitude for Ezra Camacho.md` — 2025-10-03 Gratitude for Ezra Camacho (score: 0.9818)
5. `Gratitudes/2025-10-04 Gratitude for Jing.md` — 2025-10-04 Gratitude for Jing (score: 0.9782)

**semantic:**

1. `Gratitudes/Gratitude log.md` — Gratitude Practice (score: 0.6532) **[EXPECTED]**
2. `System/Templates/Gratitude Entry.md` — Grateful for {{person/thing}} (score: 0.6509)
3. `Gratitudes/2025-10-03 Gratitude for Ezra Camacho.md` — 2025-10-03 Gratitude for Ezra Camacho (score: 0.6231)
4. `Gratitudes/2025-10-25 Cherry's Support.md` — Grateful for Cherry (score: 0.6214)
5. `Gratitudes/2025-10-03 Gratitude for Ariel Camacho.md` — 2025-10-03 Gratitude for Ariel Camacho (score: 0.6214)

**hybrid:**

1. `Gratitudes/Gratitude log.md` — Gratitude Practice (score: 0.0328) **[EXPECTED]**
2. `Gratitudes/2025-10-03 Gratitude for Ezra Camacho.md` — 2025-10-03 Gratitude for Ezra Camacho (score: 0.0315)
3. `Gratitudes/2025-10-03 Gratitude for Ariel Camacho.md` — 2025-10-03 Gratitude for Ariel Camacho (score: 0.0313)
4. `Gratitudes/2025-10-04 Gratitude for Jing.md` — 2025-10-04 Gratitude for Jing (score: 0.0305)
5. `People/Michaela.md` — Real Estate Agent (score: 0.0296)

</details>

### Query: `thesis`
**Expected:** Hubs/HUB — Thesis (Business Informatics).md, Projects/Website/Blog/2022-02-10-early-docker-adoption-2014.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 12 | 4.3973 | `Hubs/HUB — Thesis (Business Informatics).md` | #1 | PASS (rank 1) |
| semantic | 20 | 0.5896 | `Hubs/HUB — Thesis (Business Informatics).md` | #1 | PASS (rank 1) |
| hybrid | 20 | 0.0328 | `Hubs/HUB — Thesis (Business Informatics).md` | #1 | PASS (rank 1) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

1. `Hubs/HUB — Thesis (Business Informatics).md` — Thesis — Business Informatics (score: 4.3973) **[EXPECTED]**
2. `System/Templates/Thesis Progress.md` — Thesis Progress — {{date:YYYY-MM-DD}} (score: 4.2104)
3. `System/Personal Context — Joey Hipolito.md` — Personal Context — Joey Hipolito (score: 4.1276)
4. `System/System Design — Personalized Features.md` — Personalized Features for Joey (score: 4.1090)
5. `Daily/2025/10/2025-10-25.md` — Saturday, Oct 25, 2025 — Daily (score: 4.0317)

**semantic:**

1. `Hubs/HUB — Thesis (Business Informatics).md` — Thesis — Business Informatics (score: 0.5896) **[EXPECTED]**
2. `Projects/ecal-theming-i18n-deepdive.md` — ECAL White-Label Theming & i18n Deep Dive (score: 0.5607)
3. `Inbox/20260308-212511.md` — 20260308-212511 (score: 0.5587)
4. `Research/alternative-research-angles-ai-dev.md` — Alternative Research Angles — AI + Software Development (score: 0.5524)
5. `Projects/sched-fyi-analysis.md` — sched.fyi — Platform Analysis (score: 0.5514)

**hybrid:**

1. `Hubs/HUB — Thesis (Business Informatics).md` — Thesis — Business Informatics (score: 0.0328) **[EXPECTED]**
2. `System/plugin-ideas.md` — Plugin & Skill Proposals (score: 0.0292)
3. `System/System Design — Personalized Features.md` — Personalized Features for Joey (score: 0.0283)
4. `System/Templates/Thesis Progress.md` — Thesis Progress — {{date:YYYY-MM-DD}} (score: 0.0276)
5. `System/Personal Context — Joey Hipolito.md` — Personal Context — Joey Hipolito (score: 0.0263)

</details>

### Query: `trading platform`
**Expected:** Projects/Website/Stories/titanfx-dashboard-modernization.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 4 | 9.6114 | `...site/Stories/titanfx-dashboard-modernization.md` | #1 | PASS (rank 1) |
| semantic | 20 | 0.5521 | `Projects/sched-fyi-analysis.md` | #3 | PARTIAL (rank 3) |
| hybrid | 20 | 0.0323 | `...site/Stories/titanfx-dashboard-modernization.md` | #1 | PASS (rank 1) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

1. `Projects/Website/Stories/titanfx-dashboard-modernization.md` — What Nobody Tells You About Migrating a Live Trading Platform (score: 9.6114) **[EXPECTED]**
2. `People/Frank Li.md` — Engineering Manager (B2B Team) (score: 7.7440)
3. `Projects/ecal-architecture.md` — ECAL Architecture — Widget Ecosystem (score: 3.8163)
4. `Projects/ecal-security-deep-dive.md` — ECAL Security Deep Dive — Popup Evasion, Jailbreaks, API Signing (score: 3.7642)

**semantic:**

1. `Projects/sched-fyi-analysis.md` — sched.fyi — Platform Analysis (score: 0.5521)
2. `Projects/Via/Revenue Plan Updates.md` — Revenue Plan Updates (score: 0.5214)
3. `Projects/Website/Stories/titanfx-dashboard-modernization.md` — What Nobody Tells You About Migrating a Live Trading Platform (score: 0.5185) **[EXPECTED]**
4. `Inbox/20260308-212511.md` — 20260308-212511 (score: 0.5074)
5. `Projects/fyp-architecture.md` — FYP — Interest & Feed System Architecture (score: 0.5009)

**hybrid:**

1. `Projects/Website/Stories/titanfx-dashboard-modernization.md` — What Nobody Tells You About Migrating a Live Trading Platform (score: 0.0323) **[EXPECTED]**
2. `Projects/ecal-architecture.md` — ECAL Architecture — Widget Ecosystem (score: 0.0274)
3. `Projects/sched-fyi-analysis.md` — sched.fyi — Platform Analysis (score: 0.0164)
4. `People/Frank Li.md` — Engineering Manager (B2B Team) (score: 0.0161)
5. `Projects/Via/Revenue Plan Updates.md` — Revenue Plan Updates (score: 0.0161)

</details>

### Query: `scout`
**Expected:** Projects/Via/Scout.md, Projects/Website/Blog/2026-02-12-building-ai-intelligence-layer-pure-go.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 5 | 6.3332 | `Projects/Via/Scout.md` | #1 | PASS (rank 1) |
| semantic | 20 | 0.6475 | `Projects/Via/Scout.md` | #1 | PASS (rank 1) |
| hybrid | 20 | 0.0328 | `Projects/Via/Scout.md` | #1 | PASS (rank 1) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

1. `Projects/Via/Scout.md` — Scout — Project Hub (score: 6.3332) **[EXPECTED]**
2. `Projects/Via/AI.ssistant.md` — AI.ssistant — Project Hub (score: 5.6186)
3. `Projects/Website/Blog/2026-02-12-building-ai-intelligence-layer-pure-go.md` — Building an AI Intelligence Layer in Pure Go (score: 5.6153) **[EXPECTED]**
4. `Projects/Website/Blog/2026-01-31-six-plugins-one-brain.md` — Six Plugins, One Brain (score: 4.4277)
5. `Projects/Via/Via — Personal Intelligence OS.md` — Via — Personal Intelligence OS (score: 3.9510)

**semantic:**

1. `Projects/Via/Scout.md` — Scout — Project Hub (score: 0.6475) **[EXPECTED]**
2. `Projects/Website/Blog/2026-02-12-building-ai-intelligence-layer-pure-go.md` — Building an AI Intelligence Layer in Pure Go (score: 0.6037) **[EXPECTED]**
3. `Inbox/20260308-212511.md` — 20260308-212511 (score: 0.5713)
4. `Projects/Via/AI.ssistant.md` — AI.ssistant — Project Hub (score: 0.5533)
5. `Projects/Via/ContentKit.md` — ContentKit (score: 0.5479)

**hybrid:**

1. `Projects/Via/Scout.md` — Scout — Project Hub (score: 0.0328) **[EXPECTED]**
2. `Projects/Website/Blog/2026-02-12-building-ai-intelligence-layer-pure-go.md` — Building an AI Intelligence Layer in Pure Go (score: 0.0320) **[EXPECTED]**
3. `Projects/Via/AI.ssistant.md` — AI.ssistant — Project Hub (score: 0.0318)
4. `Projects/Via/Via — Personal Intelligence OS.md` — Via — Personal Intelligence OS (score: 0.0293)
5. `Projects/Website/Blog/2026-01-31-six-plugins-one-brain.md` — Six Plugins, One Brain (score: 0.0286)

</details>

## Semantic/Conceptual

### Query: `how to organize my notes effectively`
**Expected:** System/Note Writing Rules & Conventions.md, System/Vision & Design Document.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 0 | — | — | — | FAIL (no results) |
| semantic | 20 | 0.6255 | `System/Note Writing Rules & Conventions.md` | #1 | PASS (rank 1) |
| hybrid | 20 | 0.0164 | `System/Note Writing Rules & Conventions.md` | #1 | PASS (rank 1) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

(no results)

**semantic:**

1. `System/Note Writing Rules & Conventions.md` — Note Writing Rules & Conventions (score: 0.6255) **[EXPECTED]**
2. `People/Diana Castano.md` — Professor (score: 0.5921)
3. `System/Vision & Design Document.md` — Vision & Design Document (score: 0.5885) **[EXPECTED]**
4. `Projects/Gmail-Triage/Action-Plan.md` — Gmail Triage Action Plan (score: 0.5885)
5. `System/Pre-Migration Snapshot.md` — Pre-Migration Snapshot (score: 0.5797)

**hybrid:**

1. `System/Note Writing Rules & Conventions.md` — Note Writing Rules & Conventions (score: 0.0164) **[EXPECTED]**
2. `People/Diana Castano.md` — Professor (score: 0.0161)
3. `System/Vision & Design Document.md` — Vision & Design Document (score: 0.0159) **[EXPECTED]**
4. `Projects/Gmail-Triage/Action-Plan.md` — Gmail Triage Action Plan (score: 0.0156)
5. `System/Pre-Migration Snapshot.md` — Pre-Migration Snapshot (score: 0.0154)

</details>

### Query: `personal knowledge management system`
**Expected:** System/Vision & Design Document.md, Projects/Via/Via — Personal Intelligence OS.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 3 | 12.4265 | `.../2026-01-12-lifeos-personal-operating-system.md` | #2 | PARTIAL (rank 2) |
| semantic | 20 | 0.6599 | `System/Implementation Plan.md` | #2 | PARTIAL (rank 2) |
| hybrid | 20 | 0.0323 | `System/Vision & Design Document.md` | #1 | PASS (rank 1) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

1. `Projects/Website/Blog/2026-01-12-lifeos-personal-operating-system.md` — LifeOS: Building an AI-Powered Personal Operating System with Claude Code & Obsidian (score: 12.4265)
2. `System/Vision & Design Document.md` — Vision & Design Document (score: 6.7789) **[EXPECTED]**
3. `System/Note Writing Rules & Conventions.md` — Note Writing Rules & Conventions (score: 5.6647)

**semantic:**

1. `System/Implementation Plan.md` — Implementation Plan (score: 0.6599)
2. `System/Vision & Design Document.md` — Vision & Design Document (score: 0.6564) **[EXPECTED]**
3. `System/Personal Context — Joey Hipolito.md` — Personal Context — Joey Hipolito (score: 0.6205)
4. `Projects/Via/Via — Personal Intelligence OS.md` — Via — Personal Intelligence OS (score: 0.6179) **[EXPECTED]**
5. `Projects/Website/Blog/2026-01-12-lifeos-personal-operating-system.md` — LifeOS: Building an AI-Powered Personal Operating System with Claude Code & Obsidian (score: 0.6163)

**hybrid:**

1. `System/Vision & Design Document.md` — Vision & Design Document (score: 0.0323) **[EXPECTED]**
2. `Projects/Website/Blog/2026-01-12-lifeos-personal-operating-system.md` — LifeOS: Building an AI-Powered Personal Operating System with Claude Code & Obsidian (score: 0.0318)
3. `System/Note Writing Rules & Conventions.md` — Note Writing Rules & Conventions (score: 0.0285)
4. `System/Implementation Plan.md` — Implementation Plan (score: 0.0164)
5. `System/Personal Context — Joey Hipolito.md` — Personal Context — Joey Hipolito (score: 0.0159)

</details>

### Query: `building AI tools with no external dependencies`
**Expected:** Projects/Website/Blog/2026-02-12-building-ai-intelligence-layer-pure-go.md, Projects/Via/Scout.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 0 | — | — | — | FAIL (no results) |
| semantic | 20 | 0.6458 | `...02-12-building-ai-intelligence-layer-pure-go.md` | #1 | PASS (rank 1) |
| hybrid | 20 | 0.0164 | `...02-12-building-ai-intelligence-layer-pure-go.md` | #1 | PASS (rank 1) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

(no results)

**semantic:**

1. `Projects/Website/Blog/2026-02-12-building-ai-intelligence-layer-pure-go.md` — Building an AI Intelligence Layer in Pure Go (score: 0.6458) **[EXPECTED]**
2. `Research/alternative-research-angles-ai-dev.md` — Alternative Research Angles — AI + Software Development (score: 0.6254)
3. `Projects/Website/Blog/2026-02-08-teaching-ai-to-learn-from-its-mistakes.md` — Teaching AI to Learn From Its Mistakes (score: 0.6193)
4. `Projects/Website/Blog/2026-01-31-six-plugins-one-brain.md` — Six Plugins, One Brain (score: 0.6167)
5. `Projects/Website/Blog/2026-01-22-why-i-built-a-personal-intelligence-os.md` — Why I Built a Personal Intelligence OS (score: 0.6127)

**hybrid:**

1. `Projects/Website/Blog/2026-02-12-building-ai-intelligence-layer-pure-go.md` — Building an AI Intelligence Layer in Pure Go (score: 0.0164) **[EXPECTED]**
2. `Research/alternative-research-angles-ai-dev.md` — Alternative Research Angles — AI + Software Development (score: 0.0161)
3. `Projects/Website/Blog/2026-02-08-teaching-ai-to-learn-from-its-mistakes.md` — Teaching AI to Learn From Its Mistakes (score: 0.0159)
4. `Projects/Website/Blog/2026-01-31-six-plugins-one-brain.md` — Six Plugins, One Brain (score: 0.0156)
5. `Projects/Website/Blog/2026-01-22-why-i-built-a-personal-intelligence-os.md` — Why I Built a Personal Intelligence OS (score: 0.0154)

</details>

### Query: `cost optimization for AI models`
**Expected:** Projects/Website/Blog/2026-02-05-the-186x-cost-reduction-multi-llm-routing.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 0 | — | — | — | FAIL (no results) |
| semantic | 20 | 0.6248 | `...te/Blog/2026-01-12-how-i-use-ai-claude-swarm.md` | #2 | PARTIAL (rank 2) |
| hybrid | 20 | 0.0164 | `...te/Blog/2026-01-12-how-i-use-ai-claude-swarm.md` | #2 | PARTIAL (rank 2) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

(no results)

**semantic:**

1. `Projects/Website/Blog/2026-01-12-how-i-use-ai-claude-swarm.md` — Why I Built a Multi-LLM Orchestration System (And You Might Want One Too) (score: 0.6248)
2. `Projects/Website/Blog/2026-02-05-the-186x-cost-reduction-multi-llm-routing.md` — Why I Route Research to Gemini (And Keep Claude for the Hard Stuff) (score: 0.5999) **[EXPECTED]**
3. `Research/alternative-research-angles-ai-dev.md` — Alternative Research Angles — AI + Software Development (score: 0.5937)
4. `Projects/Website/Blog/2026-02-12-building-ai-intelligence-layer-pure-go.md` — Building an AI Intelligence Layer in Pure Go (score: 0.5832)
5. `Inbox/20260226-144631.md` — 20260226-144631 (score: 0.5798)

**hybrid:**

1. `Projects/Website/Blog/2026-01-12-how-i-use-ai-claude-swarm.md` — Why I Built a Multi-LLM Orchestration System (And You Might Want One Too) (score: 0.0164)
2. `Projects/Website/Blog/2026-02-05-the-186x-cost-reduction-multi-llm-routing.md` — Why I Route Research to Gemini (And Keep Claude for the Hard Stuff) (score: 0.0161) **[EXPECTED]**
3. `Research/alternative-research-angles-ai-dev.md` — Alternative Research Angles — AI + Software Development (score: 0.0159)
4. `Projects/Website/Blog/2026-02-12-building-ai-intelligence-layer-pure-go.md` — Building an AI Intelligence Layer in Pure Go (score: 0.0156)
5. `Inbox/20260226-144631.md` — 20260226-144631 (score: 0.0154)

</details>

### Query: `family celebrations in New Zealand`
**Expected:** Family/Milestones/2025-10-31 — Hiide First Halloween.md, Gratitudes/2025-10-30 Birthday Celebration with Camachons.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 0 | — | — | — | FAIL (no results) |
| semantic | 20 | 0.6349 | `...25-10-30 Birthday Celebration with Camachons.md` | #1 | PASS (rank 1) |
| hybrid | 20 | 0.0164 | `...25-10-30 Birthday Celebration with Camachons.md` | #1 | PASS (rank 1) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

(no results)

**semantic:**

1. `Gratitudes/2025-10-30 Birthday Celebration with Camachons.md` — 2025-10-30 Birthday Celebration with Camachons (score: 0.6349) **[EXPECTED]**
2. `Family/Milestones/2025-10-31 — Hiide First Halloween.md` — 2025-10-31 — Hiide First Halloween (score: 0.6319) **[EXPECTED]**
3. `Family/Milestones/2025-10-31 — Kaede First Halloween.md` — 2025-10-31 — Kaede First Halloween (score: 0.6286)
4. `Daily/2025/10/2025-10-31.md` — Friday, Oct 31, 2025 — Daily (score: 0.6277)
5. `People/Veneracion, Marivic Roldan.md` — Veneracion, Marivic Roldan (score: 0.6252)

**hybrid:**

1. `Gratitudes/2025-10-30 Birthday Celebration with Camachons.md` — 2025-10-30 Birthday Celebration with Camachons (score: 0.0164) **[EXPECTED]**
2. `Family/Milestones/2025-10-31 — Hiide First Halloween.md` — 2025-10-31 — Hiide First Halloween (score: 0.0161) **[EXPECTED]**
3. `Family/Milestones/2025-10-31 — Kaede First Halloween.md` — 2025-10-31 — Kaede First Halloween (score: 0.0159)
4. `Daily/2025/10/2025-10-31.md` — Friday, Oct 31, 2025 — Daily (score: 0.0156)
5. `People/Veneracion, Marivic Roldan.md` — Veneracion, Marivic Roldan (score: 0.0154)

</details>

### Query: `career interview preparation`
**Expected:** People/Frank Li.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 0 | — | — | — | FAIL (no results) |
| semantic | 20 | 0.5718 | `People/Frank Li.md` | #1 | PASS (rank 1) |
| hybrid | 20 | 0.0164 | `People/Frank Li.md` | #1 | PASS (rank 1) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

(no results)

**semantic:**

1. `People/Frank Li.md` — Engineering Manager (B2B Team) (score: 0.5718) **[EXPECTED]**
2. `People/Eilish McGovern.md` — People & Culture (score: 0.5656)
3. `People/Josh Walker.md` — Engineer (score: 0.5609)
4. `People/Dianne Tennent.md` — Designer (score: 0.5520)
5. `People/Pat Horsley.md` — Engineer (score: 0.5469)

**hybrid:**

1. `People/Frank Li.md` — Engineering Manager (B2B Team) (score: 0.0164) **[EXPECTED]**
2. `People/Eilish McGovern.md` — People & Culture (score: 0.0161)
3. `People/Josh Walker.md` — Engineer (score: 0.0159)
4. `People/Dianne Tennent.md` — Designer (score: 0.0156)
5. `People/Pat Horsley.md` — Engineer (score: 0.0154)

</details>

### Query: `teaching machines to remember things`
**Expected:** Projects/Website/Blog/2026-02-08-teaching-ai-to-learn-from-its-mistakes.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 0 | — | — | — | FAIL (no results) |
| semantic | 20 | 0.6725 | `...02-08-teaching-ai-to-learn-from-its-mistakes.md` | #1 | PASS (rank 1) |
| hybrid | 20 | 0.0164 | `...02-08-teaching-ai-to-learn-from-its-mistakes.md` | #1 | PASS (rank 1) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

(no results)

**semantic:**

1. `Projects/Website/Blog/2026-02-08-teaching-ai-to-learn-from-its-mistakes.md` — Teaching AI to Learn From Its Mistakes (score: 0.6725) **[EXPECTED]**
2. `Projects/Website/Blog/2026-02-11-what-1600-ai-learnings-reveal.md` — What 1,600+ AI Learnings Reveal (score: 0.5738)
3. `Writing/the-vampiric-effect-of-ai.md` — The Vampiric Effect of AI (score: 0.5701)
4. `Projects/Website/Blog/2026-02-13-i-built-a-primordial-soup-and-it-chose-mediocrity.md` — I Built a Primordial Soup for AI Personas — It Chose Mediocrity (score: 0.5587)
5. `Projects/Website/Blog/2026-02-12-building-ai-intelligence-layer-pure-go.md` — Building an AI Intelligence Layer in Pure Go (score: 0.5523)

**hybrid:**

1. `Projects/Website/Blog/2026-02-08-teaching-ai-to-learn-from-its-mistakes.md` — Teaching AI to Learn From Its Mistakes (score: 0.0164) **[EXPECTED]**
2. `Projects/Website/Blog/2026-02-11-what-1600-ai-learnings-reveal.md` — What 1,600+ AI Learnings Reveal (score: 0.0161)
3. `Writing/the-vampiric-effect-of-ai.md` — The Vampiric Effect of AI (score: 0.0159)
4. `Projects/Website/Blog/2026-02-13-i-built-a-primordial-soup-and-it-chose-mediocrity.md` — I Built a Primordial Soup for AI Personas — It Chose Mediocrity (score: 0.0156)
5. `Projects/Website/Blog/2026-02-12-building-ai-intelligence-layer-pure-go.md` — Building an AI Intelligence Layer in Pure Go (score: 0.0154)

</details>

### Query: `migrating legacy frontend code`
**Expected:** Projects/Website/Stories/titanfx-dashboard-modernization.md, Projects/ecal-architecture.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 0 | — | — | — | FAIL (no results) |
| semantic | 20 | 0.6565 | `Projects/ecal-engineering-patterns.md` | #2 | PARTIAL (rank 2) |
| hybrid | 20 | 0.0164 | `Projects/ecal-engineering-patterns.md` | #2 | PARTIAL (rank 2) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

(no results)

**semantic:**

1. `Projects/ecal-engineering-patterns.md` — Engineering Patterns Worth Preserving: A 2015 Codebase Audit (score: 0.6565)
2. `Projects/ecal-architecture.md` — ECAL Architecture — Widget Ecosystem (score: 0.6150) **[EXPECTED]**
3. `Projects/ecal-theming-i18n-deepdive.md` — ECAL White-Label Theming & i18n Deep Dive (score: 0.6042)
4. `Projects/Website/Stories/titanfx-dashboard-modernization.md` — What Nobody Tells You About Migrating a Live Trading Platform (score: 0.5931) **[EXPECTED]**
5. `People/Chris Kelly.md` — Frontend Developer (score: 0.5770)

**hybrid:**

1. `Projects/ecal-engineering-patterns.md` — Engineering Patterns Worth Preserving: A 2015 Codebase Audit (score: 0.0164)
2. `Projects/ecal-architecture.md` — ECAL Architecture — Widget Ecosystem (score: 0.0161) **[EXPECTED]**
3. `Projects/ecal-theming-i18n-deepdive.md` — ECAL White-Label Theming & i18n Deep Dive (score: 0.0159)
4. `Projects/Website/Stories/titanfx-dashboard-modernization.md` — What Nobody Tells You About Migrating a Live Trading Platform (score: 0.0156) **[EXPECTED]**
5. `People/Chris Kelly.md` — Frontend Developer (score: 0.0154)

</details>

## Multi-Word Phrase

### Query: `daily note template`
**Expected:** System/Templates/Daily Note.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 4 | 10.1114 | `System/Vision & Design Document.md` | not found | FAIL (expected not found) |
| semantic | 20 | 0.7295 | `System/Templates/Daily Note.md` | #1 | PASS (rank 1) |
| hybrid | 20 | 0.0297 | `System/Vision & Design Document.md` | #5 | FAIL (rank 5) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

1. `System/Vision & Design Document.md` — Vision & Design Document (score: 10.1114)
2. `System/plugin-ideas.md` — Plugin & Skill Proposals (score: 9.9114)
3. `System/System Design — Personalized Features.md` — Personalized Features for Joey (score: 9.8206)
4. `System/Implementation Plan.md` — Implementation Plan (score: 7.6626)

**semantic:**

1. `System/Templates/Daily Note.md` — {{date:dddd, MMM D, YYYY}} — Daily (score: 0.7295) **[EXPECTED]**
2. `System/Templates/Thesis Progress.md` — Thesis Progress — {{date:YYYY-MM-DD}} (score: 0.6782)
3. `System/Templates/Gratitude Entry.md` — Grateful for {{person/thing}} (score: 0.6754)
4. `System/Templates/Weekly Review.md` — Week {{date:gggg-[W]ww}} Review (score: 0.6585)
5. `System/Templates/Family Activity.md` — {{Activity Name}} — {{date:YYYY-MM-DD}} (score: 0.6566)

**hybrid:**

1. `System/Vision & Design Document.md` — Vision & Design Document (score: 0.0297)
2. `System/plugin-ideas.md` — Plugin & Skill Proposals (score: 0.0280)
3. `System/System Design — Personalized Features.md` — Personalized Features for Joey (score: 0.0276)
4. `System/Implementation Plan.md` — Implementation Plan (score: 0.0256)
5. `System/Templates/Daily Note.md` — {{date:dddd, MMM D, YYYY}} — Daily (score: 0.0164) **[EXPECTED]**

</details>

### Query: `multi-agent orchestration`
**Expected:** Projects/Website/Blog/2026-02-03-how-multi-agent-orchestration-works.md, Projects/Via/Via — Personal Intelligence OS.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | ERROR | — | `Error: keyword search failed: FTS5 search failed: SQL logic ` | — | ERROR |
| semantic | 20 | 0.7113 | `...26-02-03-how-multi-agent-orchestration-works.md` | #1 | PASS (rank 1) |
| hybrid | ERROR | — | `Error: hybrid search failed: FTS5 search failed: SQL logic e` | — | ERROR |

<details><summary>Top-5 results per mode</summary>

**keyword:** ERROR — Error: keyword search failed: FTS5 search failed: SQL logic error: no such colum

**semantic:**

1. `Projects/Website/Blog/2026-02-03-how-multi-agent-orchestration-works.md` — How Multi-Agent Orchestration Works (score: 0.7113) **[EXPECTED]**
2. `Projects/Via/Orchestrator.md` — Orchestrator — Project Hub (score: 0.6793)
3. `Projects/Via/Via — Personal Intelligence OS.md` — Via — Personal Intelligence OS (score: 0.5702) **[EXPECTED]**
4. `Projects/Website/Blog/2026-01-12-how-i-use-ai-claude-swarm.md` — Why I Built a Multi-LLM Orchestration System (And You Might Want One Too) (score: 0.5619)
5. `Projects/Website/Blog/2026-01-31-six-plugins-one-brain.md` — Six Plugins, One Brain (score: 0.5595)

**hybrid:** ERROR — Error: hybrid search failed: FTS5 search failed: SQL logic error: no such column

</details>

### Query: `reciprocal rank fusion`
**Expected:** Projects/Via/Obsidian CLI.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 1 | 16.1192 | `Projects/Via/Obsidian CLI.md` | #1 | PASS (rank 1) |
| semantic | 20 | 0.5719 | `Ideas/20260221-103656.md` | not found | FAIL (expected not found) |
| hybrid | 20 | 0.0164 | `Projects/Via/Obsidian CLI.md` | #1 | PASS (rank 1) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

1. `Projects/Via/Obsidian CLI.md` — Obsidian CLI — Project Hub (score: 16.1192) **[EXPECTED]**

**semantic:**

1. `Ideas/20260221-103656.md` — 20260221-103656 (score: 0.5719)
2. `Inbox/20260226-144631.md` — 20260226-144631 (score: 0.5553)
3. `Projects/Via/Revenue Plan Updates.md` — Revenue Plan Updates (score: 0.5341)
4. `Projects/fyp-architecture.md` — FYP — Interest & Feed System Architecture (score: 0.5292)
5. `Inbox/20260308-212511.md` — 20260308-212511 (score: 0.5273)

**hybrid:**

1. `Projects/Via/Obsidian CLI.md` — Obsidian CLI — Project Hub (score: 0.0164) **[EXPECTED]**
2. `Ideas/20260221-103656.md` — 20260221-103656 (score: 0.0164)
3. `Inbox/20260226-144631.md` — 20260226-144631 (score: 0.0161)
4. `Projects/Via/Revenue Plan Updates.md` — Revenue Plan Updates (score: 0.0159)
5. `Projects/fyp-architecture.md` — FYP — Interest & Feed System Architecture (score: 0.0156)

</details>

### Query: `security vulnerabilities in web widgets`
**Expected:** Projects/ecal-security-deep-dive.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 0 | — | — | — | FAIL (no results) |
| semantic | 20 | 0.6438 | `Projects/ecal-security-deep-dive.md` | #1 | PASS (rank 1) |
| hybrid | 20 | 0.0164 | `Projects/ecal-security-deep-dive.md` | #1 | PASS (rank 1) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

(no results)

**semantic:**

1. `Projects/ecal-security-deep-dive.md` — ECAL Security Deep Dive — Popup Evasion, Jailbreaks, API Signing (score: 0.6438) **[EXPECTED]**
2. `Projects/ecal-theming-i18n-deepdive.md` — ECAL White-Label Theming & i18n Deep Dive (score: 0.6262)
3. `Projects/ecal-architecture.md` — ECAL Architecture — Widget Ecosystem (score: 0.6133)
4. `Projects/ecal-engineering-patterns.md` — Engineering Patterns Worth Preserving: A 2015 Codebase Audit (score: 0.6110)
5. `Projects/sched-fyi-analysis.md` — sched.fyi — Platform Analysis (score: 0.5057)

**hybrid:**

1. `Projects/ecal-security-deep-dive.md` — ECAL Security Deep Dive — Popup Evasion, Jailbreaks, API Signing (score: 0.0164) **[EXPECTED]**
2. `Projects/ecal-theming-i18n-deepdive.md` — ECAL White-Label Theming & i18n Deep Dive (score: 0.0161)
3. `Projects/ecal-architecture.md` — ECAL Architecture — Widget Ecosystem (score: 0.0159)
4. `Projects/ecal-engineering-patterns.md` — Engineering Patterns Worth Preserving: A 2015 Codebase Audit (score: 0.0156)
5. `Projects/sched-fyi-analysis.md` — sched.fyi — Platform Analysis (score: 0.0154)

</details>

### Query: `household cleaning routines`
**Expected:** Home/Household Routines Checklist.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 3 | 14.9255 | `Hubs/HUB — Home & Routines.md` | #2 | PARTIAL (rank 2) |
| semantic | 20 | 0.7350 | `Home/Household Routines Checklist.md` | #1 | PASS (rank 1) |
| hybrid | 20 | 0.0325 | `Hubs/HUB — Home & Routines.md` | #2 | PARTIAL (rank 2) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

1. `Hubs/HUB — Home & Routines.md` — Home & Routines — Area Hub (score: 14.9255)
2. `Home/Household Routines Checklist.md` — Household Routines Checklist (score: 14.6704) **[EXPECTED]**
3. `System/System Design — Personalized Features.md` — Personalized Features for Joey (score: 10.0422)

**semantic:**

1. `Home/Household Routines Checklist.md` — Household Routines Checklist (score: 0.7350) **[EXPECTED]**
2. `Hubs/HUB — Home & Routines.md` — Home & Routines — Area Hub (score: 0.6737)
3. `System/System Design — Personalized Features.md` — Personalized Features for Joey (score: 0.5624)
4. `System/Note Writing Rules & Conventions.md` — Note Writing Rules & Conventions (score: 0.5372)
5. `Daily/2025/10/2025-10-25.md` — Saturday, Oct 25, 2025 — Daily (score: 0.5267)

**hybrid:**

1. `Hubs/HUB — Home & Routines.md` — Home & Routines — Area Hub (score: 0.0325)
2. `Home/Household Routines Checklist.md` — Household Routines Checklist (score: 0.0325) **[EXPECTED]**
3. `System/System Design — Personalized Features.md` — Personalized Features for Joey (score: 0.0317)
4. `System/Note Writing Rules & Conventions.md` — Note Writing Rules & Conventions (score: 0.0156)
5. `Daily/2025/10/2025-10-25.md` — Saturday, Oct 25, 2025 — Daily (score: 0.0154)

</details>

### Query: `image generation API comparison`
**Expected:** Inbox/20260226-144631.md

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 0 | — | — | — | FAIL (no results) |
| semantic | 20 | 0.7238 | `Inbox/20260226-144631.md` | #1 | PASS (rank 1) |
| hybrid | 20 | 0.0164 | `Inbox/20260226-144631.md` | #1 | PASS (rank 1) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

(no results)

**semantic:**

1. `Inbox/20260226-144631.md` — 20260226-144631 (score: 0.7238) **[EXPECTED]**
2. `Projects/ecal-security-deep-dive.md` — ECAL Security Deep Dive — Popup Evasion, Jailbreaks, API Signing (score: 0.5641)
3. `Projects/Website/Blog/2026-02-05-the-186x-cost-reduction-multi-llm-routing.md` — Why I Route Research to Gemini (And Keep Claude for the Hard Stuff) (score: 0.5619)
4. `Projects/sched-fyi-analysis.md` — sched.fyi — Platform Analysis (score: 0.5564)
5. `Projects/ecal-theming-i18n-deepdive.md` — ECAL White-Label Theming & i18n Deep Dive (score: 0.5529)

**hybrid:**

1. `Inbox/20260226-144631.md` — 20260226-144631 (score: 0.0164) **[EXPECTED]**
2. `Projects/ecal-security-deep-dive.md` — ECAL Security Deep Dive — Popup Evasion, Jailbreaks, API Signing (score: 0.0161)
3. `Projects/Website/Blog/2026-02-05-the-186x-cost-reduction-multi-llm-routing.md` — Why I Route Research to Gemini (And Keep Claude for the Hard Stuff) (score: 0.0159)
4. `Projects/sched-fyi-analysis.md` — sched.fyi — Platform Analysis (score: 0.0156)
5. `Projects/ecal-theming-i18n-deepdive.md` — ECAL White-Label Theming & i18n Deep Dive (score: 0.0154)

</details>

## Negative (Should Return Nothing)

### Query: `quantum chromodynamics hadron collider`
**Expected:** (none — should return empty)

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 0 | — | — | — | PASS (correct empty) |
| semantic | 20 | 0.5028 | `Inbox/20260301-215021.md` | N/A | FAIL (false positive, 20 results) |
| hybrid | 20 | 0.0164 | `Inbox/20260301-215021.md` | N/A | FAIL (false positive, 20 results) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

(no results)

**semantic:**

1. `Inbox/20260301-215021.md` — 20260301-215021 (score: 0.5028)
2. `Inbox/20260227-190930.md` — 20260227-190930 (score: 0.4988)
3. `Projects/Via/Revenue Plan Updates.md` — Revenue Plan Updates (score: 0.4976)
4. `Inbox/20260308-212511.md` — 20260308-212511 (score: 0.4969)
5. `Inbox/20260226-144631.md` — 20260226-144631 (score: 0.4952)

**hybrid:**

1. `Inbox/20260301-215021.md` — 20260301-215021 (score: 0.0164)
2. `Inbox/20260227-190930.md` — 20260227-190930 (score: 0.0161)
3. `Projects/Via/Revenue Plan Updates.md` — Revenue Plan Updates (score: 0.0159)
4. `Inbox/20260308-212511.md` — 20260308-212511 (score: 0.0156)
5. `Inbox/20260226-144631.md` — 20260226-144631 (score: 0.0154)

</details>

### Query: `underwater basket weaving certification`
**Expected:** (none — should return empty)

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 0 | — | — | — | PASS (correct empty) |
| semantic | 20 | 0.5269 | `Inbox/20260308-212511.md` | N/A | FAIL (false positive, 20 results) |
| hybrid | 20 | 0.0164 | `Inbox/20260308-212511.md` | N/A | FAIL (false positive, 20 results) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

(no results)

**semantic:**

1. `Inbox/20260308-212511.md` — 20260308-212511 (score: 0.5269)
2. `Projects/Via/Gastown Evaluation.md` — Gastown Evaluation (score: 0.5134)
3. `Daily/2026/01/2026-01-10.md` — Saturday, Jan 10, 2026 — Daily (score: 0.5002)
4. `Ideas/20260221-103656.md` — 20260221-103656 (score: 0.4978)
5. `Inbox/mullvad.md` — mullvad (score: 0.4977)

**hybrid:**

1. `Inbox/20260308-212511.md` — 20260308-212511 (score: 0.0164)
2. `Projects/Via/Gastown Evaluation.md` — Gastown Evaluation (score: 0.0161)
3. `Daily/2026/01/2026-01-10.md` — Saturday, Jan 10, 2026 — Daily (score: 0.0159)
4. `Ideas/20260221-103656.md` — 20260221-103656 (score: 0.0156)
5. `Inbox/mullvad.md` — mullvad (score: 0.0154)

</details>

### Query: `ancient roman aqueduct engineering`
**Expected:** (none — should return empty)

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 0 | — | — | — | PASS (correct empty) |
| semantic | 20 | 0.5240 | `People/Bernat Duran.md` | N/A | FAIL (false positive, 20 results) |
| hybrid | 20 | 0.0164 | `People/Bernat Duran.md` | N/A | FAIL (false positive, 20 results) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

(no results)

**semantic:**

1. `People/Bernat Duran.md` — Lead Engineer (score: 0.5240)
2. `Projects/ecal-engineering-patterns.md` — Engineering Patterns Worth Preserving: A 2015 Codebase Audit (score: 0.5196)
3. `People/Ryan Arbolado.md` — Technical Lead (score: 0.5012)
4. `People/Frank Li.md` — Engineering Manager (B2B Team) (score: 0.4950)
5. `Projects/Website/Blog/2026-01-28-from-chatgpt-to-claude-code-evolution.md` — From ChatGPT to Claude Code: The Evolution (score: 0.4940)

**hybrid:**

1. `People/Bernat Duran.md` — Lead Engineer (score: 0.0164)
2. `Projects/ecal-engineering-patterns.md` — Engineering Patterns Worth Preserving: A 2015 Codebase Audit (score: 0.0161)
3. `People/Ryan Arbolado.md` — Technical Lead (score: 0.0159)
4. `People/Frank Li.md` — Engineering Manager (B2B Team) (score: 0.0156)
5. `Projects/Website/Blog/2026-01-28-from-chatgpt-to-claude-code-evolution.md` — From ChatGPT to Claude Code: The Evolution (score: 0.0154)

</details>

### Query: `cryptocurrency mining profitability 2024`
**Expected:** (none — should return empty)

| Mode | Results | Top Score | Top Result | Best Expected Rank | Verdict |
|------|---------|-----------|------------|-------------------|---------|
| keyword | 0 | — | — | — | PASS (correct empty) |
| semantic | 20 | 0.5056 | `Inbox/20260226-144631.md` | N/A | FAIL (false positive, 20 results) |
| hybrid | 20 | 0.0164 | `Inbox/20260226-144631.md` | N/A | FAIL (false positive, 20 results) |

<details><summary>Top-5 results per mode</summary>

**keyword:**

(no results)

**semantic:**

1. `Inbox/20260226-144631.md` — 20260226-144631 (score: 0.5056)
2. `Projects/Website/Stories/titanfx-dashboard-modernization.md` — What Nobody Tells You About Migrating a Live Trading Platform (score: 0.4978)
3. `Projects/Website/Blog/2026-01-28-from-chatgpt-to-claude-code-evolution.md` — From ChatGPT to Claude Code: The Evolution (score: 0.4972)
4. `Inbox/20260227-190930.md` — 20260227-190930 (score: 0.4964)
5. `Inbox/20260301-215021.md` — 20260301-215021 (score: 0.4963)

**hybrid:**

1. `Inbox/20260226-144631.md` — 20260226-144631 (score: 0.0164)
2. `Projects/Website/Stories/titanfx-dashboard-modernization.md` — What Nobody Tells You About Migrating a Live Trading Platform (score: 0.0161)
3. `Projects/Website/Blog/2026-01-28-from-chatgpt-to-claude-code-evolution.md` — From ChatGPT to Claude Code: The Evolution (score: 0.0159)
4. `Inbox/20260227-190930.md` — 20260227-190930 (score: 0.0156)
5. `Inbox/20260301-215021.md` — 20260301-215021 (score: 0.0154)

</details>

---

## Score Distribution Analysis

### Semantic Score Ranges by Query Category

| Category | Avg Top Score | Avg Bottom Score | Avg Spread |
|----------|--------------|-----------------|------------|
| Exact Keyword Match | 0.6009 | 0.5130 | 0.0880 |
| Semantic/Conceptual | 0.6365 | 0.5385 | 0.0979 |
| Multi-Word Phrase | 0.6859 | 0.5178 | 0.1681 |
| Negative (Should Return Nothing) | 0.5148 | 0.4762 | 0.0386 |

### Hybrid vs Semantic Ranking Comparison (Conceptual Queries)

| Query | Semantic Rank | Hybrid Rank | Delta | Regression? |
|-------|--------------|-------------|-------|-------------|
| how to organize my notes effectively | #1 | #1 | +0 | no |
| personal knowledge management system | #2 | #1 | -1 | no |
| building AI tools with no external dependenci | #1 | #1 | +0 | no |
| cost optimization for AI models | #2 | #2 | +0 | no |
| family celebrations in New Zealand | #1 | #1 | +0 | no |
| career interview preparation | #1 | #1 | +0 | no |
| teaching machines to remember things | #1 | #1 | +0 | no |
| migrating legacy frontend code | #2 | #2 | +0 | no |

### Negative Query False Positive Analysis

| Query | Semantic Results | Top Score | Threshold 0.55 Would Filter? |
|-------|-----------------|-----------|------------------------------|
| quantum chromodynamics hadron collider | 20 | 0.5028 | Yes (removes 20/20) |
| underwater basket weaving certification | 20 | 0.5269 | Yes (removes 20/20) |
| ancient roman aqueduct engineering | 20 | 0.5240 | Yes (removes 20/20) |
| cryptocurrency mining profitability 2024 | 20 | 0.5056 | Yes (removes 20/20) |

