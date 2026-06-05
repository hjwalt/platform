# Platform

Auditable AI harness

The purpose of building this is to gain more control on how tools and models gets executed, with more type safety and visibility.
Tool calls and model eval are separated into message handlers.
Extension points are not limited to golang based implementations, as long as the other integrations can listen to Kafka, this can be extended infinitely.

Some of the packages can be used independently:

- agent: general harness and tools built around OpenAI compatible endpoint and MCP sdk
- flow: simplifies implementation of async functions
- format: a collection of ways to convert to and from byte array and to mask and unmask / encrypt and decrypt byte array
- reflect: a way to safely convert to and from types
- web: simplifies the way to register pages, also serves as an example on how to build a good enough backend rendered UI with interactions via HTMX and design via web components

## Developing

### Tools Used

Please follow the tools setup guide in its respective repository or documentation.

- Golang
- Gomock
- Podman / docker
- Protoc
- Protoc-gen-go
- Psql

### Commands

Makefile is heavily used.

Commands used when coding:

```
make test
make tidy
make update
make mocks
make proto
make cov
make htmlcov
```

### Running

See the source in `main.go` or `example` folder

```
docker compose up -d
make reset
make run
```

1. `docker compose up -d` will start zookeeper, kafka, and postgresql with ports exposed to your host network detached from your terminal
2. `make reset` will clean up the topics on kafka and postgresql table, and add some example events
3. `make run` will start the example word count application

## Principles

- Container first
- Do one thing, and only one thing well
- Ease to change integrations
- Sane defaults
- Simple format helpers, with bytes as default
- Idiomatic golang

## Flow

Flows sits at the core of the [Kappa Architecture](https://hazelcast.com/glossary/kappa-architecture/), where it tackles four elements:

1. Stateless functions
2. Stateful functions
3. Join as combination of stateless and stateful functions
4. Materialisation (TODO)

### Stateless

Stateless functions can be used to perform simple operations such as:

1. Basic event filtering
2. Event mirroring
3. Merging multiple topics into one
4. Exploding events
5. Interfacing with external parties with at least once semantic

### Stateful

Stateful functions can be used to perform state machine operations on exactly one topic to perform operations such as:

1. Reduce or aggregate
2. Validation

### Join

To ensure events are published for the multiple topics that are being joined, there are two options:

1. Maintain publishing state for all topics (currently, the last result state is global per key, this can be changed to be per topic per key)
2. Merge topics into one intermediate topic, and perform a stateful function

In this codebase, option two is the preferred option for reasons of:

1. Avoiding to impose limits on the number of messages being published, which can increase the state size written into the store, which implication should be obvious
2. Stateless map that can be used to merge topics are cheaper in terms of time latency than transaction locking failure
3. Kafka can be configured to be scalable enough in terms of throughput
4. Parallelism of the intermediate topic (partition count) can be higher than the source topics
5. Avoiding transaction will reduce cost and increase speed especially for cloud services
6. Avoiding distributed data contention allows local state caching, reducing data query latency for cache hits

![Join Pattern](docs/join-pattern.drawio.png "Join Pattern")

### Materialiser

Materialise function batch upserts into database.

### Semantics

At least once publishing with effectively once state update.
Additional application based deduplication is recommended (request id header deduplication for instance).

### Integrations

- Kafka using librdkafka via its golang binding
- Postgresql using Bun

### Limitations

To keep the simplicity of implementation, temporal operations are not yet considered in this project.
Examples of temporal operations that are not considered for implementation yet:

1. End of time window only publishing. With states, a window can be emulated, but an output will be published for each message received instead of only at the end of the window.
2. Per-key publication rate limiter. Combining state storage, commit offset, and real time ticks can be implemented, however that complicates the interfaces needed.

### Kafka Migration

The way to migrate stateful functions to a mirrored Kafka cluster are by:

1. Stopping the job
2. Clearing the internal column / field
3. Starting the job

Due to the fact that offset numbers are different in mirrored Kafka cluster,
an additional application functionality side deduplication will be required to ensure that stateful operation does not get executed twice.
Such deduplication can be peformed using a unique Kafka header identifier.

However, if the application functionality can already tolerate at least once execution, then there will be no problems with migration.

### Why

Why do I build this instead of using tools like Spark, Flink, Kafka Streams?

This is my personal view based on experience with those three.
Note that I have not tested things like Pulsar functions or NATS jetstream, those might well be solving the same thing.

The three I mentioned are heavyweight data engineering tools.
It can continously process data flow with a special DSL, with exactly once semantics (at a cost) and very high throughput (also at a cost).

However, in a complex backend data flow changes, constantly.
Often times, its just one step in the middle of the flow.
Sometimes its removing a step, sometimes its adding steps, sometimes its reusing steps.
Heavyweight tools just doesn't work well with that kind of constant change.
Deploy a flow that is too small its costly, deploy a flow that is too big it constantly changes.

So the idea of a lightweight flow comes in.
Its similar to Kafka streams, but every single step is independently deployed, every intermediate data types are independently designed.
With a schema management system, Kafka, and Kubernetes, its the right balance of performance, ease of deployment, and flexibility.

Its written in golang, so that the resource use of each lightweight flow step is small, yet it can be scaled both horizontally and vertically very well.
In theory other language like C++ and Rust will also work. I attempted Rust, the type and memory safety rules just makes the implementation a far bigger hassle than I would like it to be.

## Misc

### Runtime

Runtimes are independent context that needs to be maintained from the start until the end of the program execution.
Runtimes usually also needs to be instantiated and cleaned up.

Examples of runtimes:

1. HTTP server
2. Kafka consumer
3. Kafka producer
4. Database connection

Why is it important to standardise?

1. Golang is too barebone. This kind of runtime management needs to keep getting rewritten.
2. Resource management is not a difficult problem but the cost of mistakes is high. Your program can hang, you can have connections interrupted without clean disconnect, and many more. Reducing the potential scope of failure is in general a good idea.

### Trusted

Trusted is a collection of stuff that is practically impossible to test.
Keeping this collection of untestable code as small as possible is a must.

## Credits

The web page uses the free version of web-awesome. Check them out for a really well implemented web components.

## Caveat

These code are an amalgamation of a few old repositories I have to create a more cohesive codebase. As it is going to be a constant work in progress for experiments, expect breaking interface changes.

Do not use this for production use case, however, it may be helpful if you want your own fun little chatbot.
This will also help you understand how general agent harnesses work under the hood.
