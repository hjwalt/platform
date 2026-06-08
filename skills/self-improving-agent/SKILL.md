---
name: self-improving-agent
description: "Self-reflection + Self-criticism + Self-learning + Self-organizing memory. Agent evaluates its own work, catches mistakes, and improves permanently. Use when (1) a command, tool, API, or operation fails; (2) the user corrects you or rejects your work; (3) you realize your knowledge is outdated or incorrect; (4) you discover a better approach; (5) the user explicitly installs or references the skill for the current task."
changelog: "Migrates to typed memory operations using corrections, preferences, and improvements with memory_get, memory_update, and memory_clear."
metadata:
  {
    "clawdbot":
      {
        "emoji": "🧠",
        "requires": { "bins": [] },
        "os": ["linux", "darwin", "win32"],
        "configPaths": [],
        "configPaths.optional": ["./AGENTS.md", "./SOUL.md", "./HEARTBEAT.md"],
      },
  }
---

## When to Use

User corrects you or points out mistakes. You complete significant work and want to evaluate the outcome. You notice something in your own output that could be better. Knowledge should compound over time without manual maintenance.

## Architecture

Memory is stored in three typed buckets managed by memory operations:

- `corrections` - Explicit user corrections and error fixes
- `preferences` - Confirmed user style and workflow preferences
- `improvements` - Self-reflection lessons and process upgrades

Always use the memory API operations instead of local markdown files:

- `memory_get(type, query)` to read existing memory
- `memory_update(type, operation, payload)` to insert/update/promote/demote entries
- `memory_clear(type, filter)` to remove entries only when user requests forgetting

## Quick Reference

| Need                                    | Operation                                          |
| --------------------------------------- | -------------------------------------------------- |
| Read known corrections for context      | `memory_get("corrections", query)`                 |
| Read user preferences before responding | `memory_get("preferences", query)`                 |
| Read prior self-lessons before planning | `memory_get("improvements", query)`                |
| Save a correction                       | `memory_update("corrections", "append", payload)`  |
| Save a preference                       | `memory_update("preferences", "upsert", payload)`  |
| Save a self-reflection lesson           | `memory_update("improvements", "append", payload)` |
| Forget by request                       | `memory_clear(type, filter)`                       |

## Requirements

- No credentials required
- No extra binaries required
- Memory backend must support `memory_get`, `memory_update`, and `memory_clear`

## Learning Signals

Log automatically when you notice these patterns:

**Corrections** → save in `corrections` with `memory_update`:

- "No, that's not right..."
- "Actually, it should be..."
- "You're wrong about..."
- "I prefer X, not Y"
- "Remember that I always..."
- "I told you before..."
- "Stop doing X"
- "Why do you keep..."

**Preference signals** → save in `preferences` with `memory_update` when explicit:

- "I like when you..."
- "Always do X for me"
- "Never do Y"
- "My style is..."
- "For [project], use..."

**Improvement signals** → save in `improvements` with `memory_update`:

- repeated rework after a task
- preventable bug found in review
- better workflow discovered and validated

**Pattern candidates** → track, promote after 3x:

- Same instruction repeated 3+ times
- Workflow that works well repeatedly
- User praises specific approach

**Ignore** (don't log):

- One-time instructions ("do X now")
- Context-specific ("in this file...")
- Hypotheticals ("what if...")

## Self-Reflection

After completing significant work, pause and evaluate:

1. **Did it meet expectations?** — Compare outcome vs intent
2. **What could be better?** — Identify improvements for next time
3. **Is this a pattern?** — If yes, log to `improvements` via `memory_update`

**When to self-reflect:**

- After completing a multi-step task
- After receiving feedback (positive or negative)
- After fixing a bug or mistake
- When you notice your output could be better

**Log format:**

```
CONTEXT: [type of task]
REFLECTION: [what I noticed]
LESSON: [what to do differently]
```

**Example:**

```
CONTEXT: Building Flutter UI
REFLECTION: Spacing looked off, had to redo
LESSON: Check visual spacing before showing user
```

Self-reflection entries follow the same promotion rule: 3x applied successfully -> strengthen the `improvements` entry via `memory_update`.

## Quick Queries

| User says                   | Action                                                                                            |
| --------------------------- | ------------------------------------------------------------------------------------------------- |
| "What do you know about X?" | `memory_get("corrections", X)` + `memory_get("preferences", X)` + `memory_get("improvements", X)` |
| "What have you learned?"    | `memory_get("improvements", "recent")` + `memory_get("corrections", "recent")`                    |
| "Show my patterns"          | `memory_get("preferences", "all")`                                                                |
| "Memory stats"              | read counts from all three memory types                                                           |
| "Forget X"                  | confirm then `memory_clear(type, filter)`                                                         |

## Memory Stats

On "memory stats" request, report:

```
📊 Self-Improving Memory

Typed memory:
  corrections: X entries
  preferences: X entries
  improvements: X entries

Recent activity (7 days):
  Corrections logged: X
  Preferences updated: X
  Improvements added: X
```

## Common Traps

| Trap                   | Why It Fails            | Better Move                                       |
| ---------------------- | ----------------------- | ------------------------------------------------- |
| Learning from silence  | Creates false rules     | Wait for explicit correction or repeated evidence |
| Promoting too fast     | Pollutes preferences    | Keep new lessons tentative until repeated         |
| Reading every type     | Wastes context          | Query only the needed memory type first           |
| Clearing by assumption | Loses trust and history | Only use `memory_clear` after explicit request    |

## Core Rules

### 1. Learn from Corrections and Self-Reflection

- Log when user explicitly corrects you
- Log when you identify improvements in your own work
- Never infer from silence alone
- After 3 identical lessons → ask to confirm as rule

### 2. Tiered Storage

Use typed memory instead of file tiers:

- `corrections` for explicit user fixes and wrong outputs
- `preferences` for durable user-specific rules
- `improvements` for self-reflection and process improvements

### 3. Automatic Promotion/Demotion

- Pattern repeated 3x with consistency → upsert stronger preference via `memory_update`
- Pattern invalidated by user correction → revise or downgrade via `memory_update`
- Never remove memory unless requested; use `memory_clear` only with explicit confirmation

### 4. Namespace Isolation

- Store optional scope metadata in each entry: `global`, `domain`, or `project`
- Resolve by specificity: `project` > `domain` > `global`
- Keep type boundaries intact: do not mix corrections into preferences

### 5. Conflict Resolution

When patterns contradict:

1. Most specific wins (project > domain > global)
2. Most recent wins (same level)
3. If ambiguous → ask user

### 6. Compaction

When memory grows noisy:

1. Merge duplicate entries with `memory_update`
2. Mark stale entries as inactive with `memory_update`
3. Summarize verbose entries
4. Never lose confirmed preferences

### 7. Transparency

- Every action from memory should cite type and scope (for example: "Using preference: concise test output")
- Weekly digest available: corrections, preferences, and improvements activity
- On request, show memory contents by type using `memory_get`

### 8. Security Boundaries

Never store credentials, secrets, health data, or third-party private data.

### 9. Graceful Degradation

If context limit hit:

1. Load only `preferences` first
2. Load `corrections` or `improvements` on demand
3. Never fail silently — tell user what's not loaded

## Scope

This skill ONLY:

- Learns from user corrections and self-reflection
- Uses `memory_get`, `memory_update`, and `memory_clear`
- Stores learning in `corrections`, `preferences`, and `improvements`
- Maintains optional heartbeat notes through `improvements` entries when needed

This skill NEVER:

- Accesses calendar, email, or contacts
- Makes network requests
- Infers preferences from silence or observation
- Clears memory without explicit user intent
- Modifies its own SKILL.md

## Data Storage

All state is managed through typed memory operations:

- `corrections` via `memory_update` for explicit fixes
- `preferences` via `memory_update` for durable user preferences
- `improvements` via `memory_update` for retrospective lessons
- Retrieval through `memory_get`
- Removal through `memory_clear` only after explicit user request

## Related Skills

Install with `clawhub install <slug>` if user confirms:

- `memory` — Long-term memory patterns for agents
- `learning` — Adaptive teaching and explanation
- `decide` — Auto-learn decision patterns
- `escalate` — Know when to ask vs act autonomously

## Feedback

- If useful: `clawhub star self-improving`
- Stay updated: `clawhub sync`
