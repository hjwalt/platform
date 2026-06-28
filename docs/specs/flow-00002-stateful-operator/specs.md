# Title

Stateful Operator Specification

# High Level Description

This specification defines the `flow/stateful` package, which provides a generic stateful operator that implements `flow.Handler[IV]`. The stateful operator enriches input messages with persisted key-scoped state, applies a state update function that can short-circuit to error, and delegates to a handler that may produce output, an error, or neither. It is the core building block for stateful stream-processing topologies built on the platform `flow` abstractions.

The package defines three function type contracts — `StateKey`, `StateUpdate`, and `Operate` — and an `Operator` struct that orchestrates them with a `flow.Store[ST]`, output producer, and error producer. Every path through the operator is explicit: state-key extraction, state read, state update, state write, handler execution, and message production are each a distinct step with defined error semantics.

# User Scenarios

1. As a flow developer, I want to build stateful processing topologies by composing a key-extraction strategy, a state-update function, and a handler, without manually wiring state read/write and error routing.
2. As a platform operator, I want stateful operators to short-circuit to error when the state update rejects an input, so invalid state transitions are prevented before any handler execution or output production.
3. As a stream-processing maintainer, I want stateful operators to handle the nil-output/nil-error case (filter/sink semantics) correctly, so handlers can act as filters without producing spurious messages.
4. As a backend engineer, I want the stateful operator to be generic over input, output, state, and error types so it can be reused across domains without casting or adapter boilerplate.
5. As an on-call engineer, I want every failure mode (key extraction, state read, state update rejection, state write, handler error) to be clearly distinguishable in logs and error paths.

# Functional Requirements

## FR-FLOW-00002-001 Stateful Operator Construction

- The `NewOperator` constructor must accept all required dependencies and return a `flow.Handler[IV]`.
- Given a name, `StateKey`, `StateUpdate`, `Operate`, output/error metadata extractors, output/error producers, and a `flow.Store`, when `NewOperator` is called, then an `Operator` implementing `flow.Handler[IV]` is returned.
- The constructor must be generic over four type parameters: input value (`IV`), output value (`OV`), state type (`ST`), and error type (`ERR`).

## FR-FLOW-00002-002 Stateful Operator Handle Lifecycle

- The `Handle` method must execute a well-defined sequence of steps for every input message.
- Given an input `flow.Message[IV]`, when `Handle` is called, then the operator must: (1) extract the state key, (2) read current state, (3) apply state update, (4) write updated state or short-circuit to error, (5) execute the handler, (6) produce output, error, or neither.
- Each step must return an error on failure, halting further processing.

## FR-FLOW-00002-003 State Key Extraction

- The `StateKey[IV]` function must extract a unique string key from the input value.
- Given an input value, when `StateKey` is called, then it returns a string key used for state store lookup and an error if key derivation fails.
- Given a key extraction failure, when `Handle` processes a message, then the error is returned immediately and no state read, update, write, or handler execution occurs.

## FR-FLOW-00002-004 State Update And Error Short-Circuit

- The `StateUpdate[IV, ST, ERR]` function must transform input and current state into either a new state or an error.
- Given an input value and current state, when `StateUpdate` is called and returns `either.Left[ST]`, then the new state is persisted and handler execution proceeds.
- Given an input value and current state, when `StateUpdate` is called and returns `either.Right[ERR]`, then an error message is produced via the error producer, no state is written, and the handler is never invoked.
- State-update rejection must be treated as a terminal outcome for that message — it produces exactly one error message and returns the producer's error (or nil).

## FR-FLOW-00002-005 State Read And Write

- State must be read from and written to the configured `flow.Store[ST]` using the extracted key.
- Given a state read error, when `Handle` processes a message, then the error is returned immediately and no further processing occurs.
- Given a successful state update (`either.Left`), when the new state is written, then a `flow.State[ST]` is persisted with the extracted key, updated value, and current timestamp.
- Given a state write error, when `Handle` attempts to persist, then the error is returned immediately and no handler execution or output production occurs.

## FR-FLOW-00002-006 Handler Execution And Output Production

- The `Operate[IV, OV, ST, ERR]` function must accept input, state, and context, and may produce an output value, an error value, both, or neither.
- Given a successful state write, when `Operate` is called and returns a present output (`optional.Optional[OV]`), then exactly one output message is produced via the output producer with metadata extracted by the output metadata function.
- Given a successful state write, when `Operate` is called and returns a present error (`optional.Optional[ERR]`), then exactly one error message is produced via the error producer with metadata extracted by the error metadata function.
- Given a successful state write, when `Operate` returns an absent output and an absent error, then no message is produced and `Handle` returns nil — this is the filter/sink semantic.
- Given a successful state write, when `Operate` returns both a present output and a present error, then only the output message is produced (output takes precedence over error).

## FR-FLOW-00002-007 Metadata Propagation

- Output and error messages must carry metadata derived from the input message metadata and the produced value.
- Given an output value, when `ExtractMetadata[OV]` is called, then it receives the context, the input message's `Metadata`, and the output value, and returns the output message's `Metadata`.
- Given an error value, when `ExtractMetadata[ERR]` is called, then it receives the context, the input message's `Metadata`, and the error value, and returns the error message's `Metadata`.
- The input message's `Metadata` must be passed through to both extractors so correlation identifiers are preserved.

## FR-FLOW-00002-008 Context Propagation And Logging

- The operator must enrich the context with its name before any processing.
- Given a named operator, when `Handle` is called, then the context passed to all downstream functions includes a `function` log field set to the operator's name.
- Each function in the chain (`StateKey`, `StateRead`, `StateUpdate`, `StateWrite`, `Operate`, `ProduceMessage`) must receive the enriched context.

## FR-FLOW-00002-009 Type Safety And Generics

- The `StateKey`, `StateUpdate`, and `Operate` function types must be expressed as Go type definitions with explicit generic parameters.
- `StateKey[IV any]` must accept an input value and return a string key with an error.
- `StateUpdate[IV any, ST any, ERR any]` must accept input and current state, returning `either.Either[ST, ERR]`.
- `Operate[IV any, OV any, ST any, ERR any]` must accept input and state, returning `(optional.Optional[OV], optional.Optional[ERR])`.

# Non-Functional Requirements

1. Correctness: Every code path through `Handle` must be deterministic — given the same inputs, store state, and function implementations, the same outcomes and message productions must occur.
2. Safety: No nil-pointer dereferences — `StateKey`, `StateUpdate`, `HandlerFunction`, `OutputMetadata`, `ErrorMetadata`, `OutputProducer`, `ErrorProducer`, and `StateStore` must all be non-nil for correct operation.
3. Generality: The operator must accept any types for input, output, state, and error without casting — all wiring is through generic parameters.
4. Observability: The operator name must be propagated into logging context so all log lines from a given operator instance are attributable.
5. Testability: Every function type is independently testable — `StateKey`, `StateUpdate`, and `Operate` are plain functions, not interfaces, so they can be tested in isolation before integration.

# Definition of Done

1. FR-FLOW-00002-001 through FR-FLOW-00002-009 are covered by automated tests in `flow/stateful/`.
2. Constructor tests validate that `NewOperator` returns a `flow.Handler[IV]` with all fields assigned.
3. Handle lifecycle tests validate the full sequence: key → read → update → write → handler → produce, for success, error, and filter paths.
4. State-update short-circuit tests validate that `either.Right` produces an error message without writing state or invoking the handler.
5. Handler output/error precedence tests validate that output takes priority over error when both are present.
6. Metadata propagation tests validate that input metadata flows through to output and error messages via the extractor functions.
7. Context propagation tests validate that the operator name appears in the log context.
8. State read and write error tests validate that failures halt processing at the correct step.
9. Key extraction error tests validate that failures before state access are handled correctly.
10. Filter/sink tests validate the nil-output nil-error path produces no messages.

# Testing Methodology

1. Table-driven tests using mock/fake implementations of `flow.Store`, `flow.Producer`, and function types.
2. Each test case specifies: input message, store state, StateKey behavior, StateUpdate behavior, Operate behavior, and expected outputs/errors/messages.
3. Fake store and producer capture calls for assertion — tests verify the exact messages produced, their metadata, values, and timestamps.
4. State update short-circuit tests verify that no state write occurs and no handler is called when `StateUpdate` returns `either.Right`.
5. Handler output-error precedence tests cover all four combinations: (none, none), (output, none), (none, error), (output, error).
6. Metadata extractor tests verify that input metadata fields (Id, Group, Attempt, Sequence, Source) are accessible to extractors and appear in output metadata as directed.
7. Context tests use a test logger or context value inspection to confirm the `function` field is set to the operator name.
8. Error path tests cover: key extraction failure, state read failure, state write failure, and producer failures for both output and error channels.
