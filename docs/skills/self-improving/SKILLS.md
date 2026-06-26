---
name: self-improving
description: "Self-reflection + Self-criticism + Self-learning + Self-organizing memory. Agent evaluates its own work, catches mistakes, and improves permanently. Use when (1) a command, tool, API, or operation fails; (2) the user corrects you or rejects your work; (3) you realize your knowledge is outdated or incorrect; (4) you discover a better approach; (5) the user explicitly installs or references the skill for the current task."
changelog: "Uses one operation-based memory tool per type (`*_memory`) with operation=get|update|clear."
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

Always use the typed operation-based tools instead of local markdown files:

- `<type>_memory({"operation":"get"})` to read memory content
- `<type>_memory({"operation":"update","mode":"append|replace","content":"..."})` to write memory
- `<type>_memory({"operation":"clear"})` only when user explicitly requests forgetting

## Quick Reference

| Need                                    | Operation                                                                       |
| --------------------------------------- | ------------------------------------------------------------------------------- |
| Read known corrections for context      | `corrections_memory({"operation":"get"})`                                       |
| Read user preferences before responding | `preferences_memory({"operation":"get"})`                                       |
| Read prior self-lessons before planning | `improvements_memory({"operation":"get"})`                                      |
| Save a correction                       | `corrections_memory({"operation":"update","mode":"append","content":payload})`  |
| Save a preference                       | `preferences_memory({"operation":"update","mode":"replace","content":payload})` |
| Save a self-reflection lesson           | `improvements_memory({"operation":"update","mode":"append","content":payload})` |
| Forget by request                       | `<type>_memory({"operation":"clear"})`                                          |

## Requirements

- No credentials required
- No extra binaries required
- Memory backend must support typed `*_memory` tools with `operation=get|update|clear`

## Learning Signals

Log automatically when you notice these patterns:

**Corrections** → save in `corrections_memory` with `operation=update`:

- "No, that's not right..."
- "Actually, it should be..."
- "You're wrong about..."
- "I prefer X, not Y"
- "Remember that I always..."
- "I told you before..."
- "Stop doing X"
- "Why do you keep..."

**Preference signals** → save in `preferences_memory` with `operation=update` when explicit:

- "I like when you..."
- "Always do X for me"
- "Never do Y"
- "My style is..."
- "For [project], use..."

**Improvement signals** → save in `improvements_memory` with `operation=update`:

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
3. **Is this a pattern?** — If yes, log to `improvements_memory` via `operation=update`

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

Self-reflection entries follow the same promotion rule: 3x applied successfully -> strengthen the `improvements_memory` entry via `operation=update`.

## Quick Queries

| User says                   | Action                                                                                                                             |
| --------------------------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| "What do you know about X?" | `corrections_memory({"operation":"get"})` + `preferences_memory({"operation":"get"})` + `improvements_memory({"operation":"get"})` |
| "What have you learned?"    | `improvements_memory({"operation":"get"})` + `corrections_memory({"operation":"get"})`                                             |
| "Show my patterns"          | `preferences_memory({"operation":"get"})`                                                                                          |
| "Memory stats"              | read counts from all three memory types                                                                                            |
| "Forget X"                  | confirm then `<type>_memory({"operation":"clear"})`                                                                                |

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
| Clearing by assumption | Loses trust and history | Only use `operation=clear` after explicit request |

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

- Pattern repeated 3x with consistency → upsert stronger preference via `operation=update`
- Pattern invalidated by user correction → revise or downgrade via `operation=update`
- Never remove memory unless requested; use `operation=clear` only with explicit confirmation

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

1. Merge duplicate entries with `operation=update`
2. Mark stale entries as inactive with `operation=update`
3. Summarize verbose entries
4. Never lose confirmed preferences

### 7. Transparency

- Every action from memory should cite type and scope (for example: "Using preference: concise test output")
- Weekly digest available: corrections, preferences, and improvements activity
- On request, show memory contents by type using `operation=get`

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
- Uses typed `*_memory` tools with `operation=get|update|clear`
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

- `corrections` via `corrections_memory` with `operation=update` for explicit fixes
- `preferences` via `preferences_memory` with `operation=update` for durable user preferences
- `improvements` via `improvements_memory` with `operation=update` for retrospective lessons
- Retrieval through `<type>_memory` with `operation=get`
- Removal through `<type>_memory` with `operation=clear` only after explicit user request

## Related Skills

Install with `clawhub install <slug>` if user confirms:

- `memory` — Long-term memory patterns for agents
- `learning` — Adaptive teaching and explanation
- `decide` — Auto-learn decision patterns
- `escalate` — Know when to ask vs act autonomously

## Feedback

- If useful: `clawhub star self-improving`
- Stay updated: `clawhub sync`
