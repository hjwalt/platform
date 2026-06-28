---
name: manager-agent
description: "Manage the workflow of agents and tools to accomplish the goal required"
allowed-tools:
  - web_search
  - web_fetch
  - skill
metadata:
  {
    "clawdbot":
      {
        "emoji": "🧠",
        "requires": { "bins": [] },
        "os": ["linux", "darwin", "win32"],
        "configPaths": ["~/researcher/"],
        "configPaths.optional": ["./AGENTS.md", "./SOUL.md", "./HEARTBEAT.md"],
      },
  }
---

## Workflow

1. Research
2. Prototype
3. Plan
4. Execute
5. Review

### Plan

1. Generate specifications of the goal
2. Generate the requirements based on the research result
3. Generate the sequence of task to be executed
