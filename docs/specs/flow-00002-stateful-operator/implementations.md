# Choices Made

- The stateful operator is implemented as a generic struct (`Operator[IV, OV, ST, ERR]`) rather than an interface, because the processing pipeline is a fixed sequence — polymorphism is at the function-type level (`StateKey`, `StateUpdate`, `Operate`), not the operator level.
- State read happens before state update, and state write happens before handler execution — this ensures the handler always sees the committed state and that state-update rejection prevents both writes and handler invocation.
- The `StateUpdate` function returns `either.Either[ST, ERR]` to cleanly model the success-or-reject pattern without multiple return values or sentinel errors.
- The `Operate` function returns `(optional.Optional[OV], optional.Optional[ERR])` to support all four outcomes (output, error, both, neither) with output taking precedence over error when both are present.
- Timestamps on produced messages and persisted state use `time.Now()` at the point of creation, providing wall-clock ordering.
- Context enrichment with the operator name happens once at the top of `Handle` and is propagated to all downstream calls.

# Libraries Used

- `github.com/hjwalt/platform/flow` — core flow types: `Handler`, `Message`, `State`, `Store`, `Producer`, `Metadata`, `ExtractMetadata`.
- `github.com/hjwalt/platform/type/either` — `Either[L, R]` for state-update success/rejection modeling.
- `github.com/hjwalt/platform/type/optional` — `Optional[T]` for handler output/error optionality.
- `github.com/hjwalt/platform/logger` — context-based structured logging.
- Go standard library: `context`, `time`.

# Implementation Preferences

- Prefer table-driven tests with fake `Store` and `Producer` implementations that record all calls for assertion.
- Each test case should set up its own `StateKey`, `StateUpdate`, and `Operate` function literals to control behavior precisely.
- Test fake implementations should live in `flow/stateful/` test files (not exported) to keep the package self-contained.
- Use distinct Go types for IV, OV, ST, ERR in tests (e.g., `string`, `int`, `struct{...}`, `error`) to validate generic correctness at compile time.
- Error messages in tests should use `errors.New` or `fmt.Errorf` for clear failure attribution.

# Caveats

- The spec captures the existing implementation in `flow/stateful/operator.go` and `flow/stateful/functions.go` — no behavioral changes are proposed, only test coverage is being added.
- The output-over-error precedence rule (output present + error present → only output produced) is by design but should be clearly documented as it is a lossy choice.
- The operator does not currently enforce non-nil validation of its dependencies at construction time — this is deferred to a future spec or implementation cycle.
- `time.Now()` usage in `Handle` means test assertions on timestamps should use tolerance-based comparisons rather than exact equality.
