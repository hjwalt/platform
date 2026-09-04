1. agent/skill/parser.go:131 — metadata["allowed-tools"].([]string) can never succeed (yaml.v3 decodes to []interface{}), so AllowedTools is always empty.
2. agent/harness/flow.go — st.AppendSkillLoaded(...) return value is discarded, so LoadedSkills is never persisted (skills re-load every time).
3. type/either/either.go — zero-value Either panics (nil Optional interfaces); usable only via constructors.
4. agent/util/brave_search/client.go:43 — u.Path = u.Path + path breaks for BaseUrls without a trailing slash (/res/v1web/search).
5. flow/flow_runtime_kafka — metadata roundtrip resets Sequence to 0 (never written to Offset).
6. trusted.Exit — calls os.Exit(1) on any error except ErrPrimaryTesting, including nil.
7. agent/util/mcp AddToMcp — panics on schema-less tools with go-sdk v1.6.1.
