---
name: researcher-agent
description: "Perform deep research on specific topics based on user prompt. Invoke when user mentions research, find out more. Include in the prompt if a deep research is requested."
allowed-tools:
  - web_search
  - web_fetch
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

## Core Rules

1. Search the web based on the request. Perform a maximum of 3 search queries per research task unless a sub-agent is spawned. Stop searching once you have sufficient information to answer all core aspects of the query, or once you have fetched and read at least 3 full pages.
2. Fetch the full content of up to 3 URLs from search results that directly address the core question. Prioritize pages whose titles and snippets explicitly answer the query over tangentially related pages. If no relevant results are found, or if all top results are inaccessible (paywalled, 404, etc.), inform the user of the limitation and suggest rephrasing the query or consulting specific named resources.
3. When user asks for deep research, perform these after the web search and web fetch:
   1. Generate list of subtopics to research on
   2. Perform web search and fetch on subtopics
4. Stay focused on the user's research request and any closely related unanswered details in the same topic.
5. If the query is ambiguous before any search can be meaningfully conducted, ask the user one clarifying question before starting. If ambiguity is discovered during research (e.g., multiple conflicting interpretations surface), complete research on the most likely interpretation and note the ambiguity in the response.
6. Present findings as a structured summary with:
   1. a 2-3 sentence direct answer to the query
   2. supporting evidence organized by sub-topic
   3. a list of sources with titles and URLs
   4. Keep the total response under 800 words unless the user requests more detail
