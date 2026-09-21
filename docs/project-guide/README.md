# Tricolor GC project guide

> Generated: 2026-09-21 from commit 82e9409.

## What this is

This is an educational simulator that retains objects reachable from explicit roots and cleans up unreachable objects using tricolor mark/sweep ([contract](../../V1.md#contract)). It is for readers exploring collection algorithms; Go still manages the actual process memory ([README](../../README.md#tricolor-garbage-collector-simulator-v1)). It runs as a single Go executable with an in-memory graph and synchronous collection ([entry point](../../main.go#L5), [collector](../../gc.go#L101)).

## Run it

From the repository root, with Go compatible with the [1.23.5 directive](../../go.mod#L3):

```bash
go run .
go test -race ./...
make build
./tricolor-gc
```

The [demo](../../main.go#L5) takes no arguments or environment configuration. The [Makefile](../../Makefile#L13) builds the two implementation files explicitly. Race tests are documented as passed in [delivery verification](../../HISTORY.md#delivery-verification); this guide update did not rerun them.

## The five-file tour

There are four Go source files; the module manifest completes this tour.

| # | File | Why this one | Then look at |
| --- | --- | --- | --- |
| 1 | [main.go](../../main.go#L5) | Builds the example graph and requests collection. | [Demo flow](02-flow.md#demo-lifecycle) |
| 2 | [gc.go](../../gc.go#L17) | Owns graph, marking, accounting and lifecycle. | [Components](01-architecture.md#components) |
| 3 | [gc_test.go](../../gc_test.go#L44) | Basic reachability expectations. | [Tests](03-structure.md#root) |
| 4 | [lifecycle_test.go](../../lifecycle_test.go#L8) | Cycles, cleanup, concurrency and invalid input. | [Decisions](05-decisions.md) |
| 5 | [go.mod](../../go.mod#L1) | Runtime requirement and test dependencies. | [Tech stack](04-tech-stack.md) |

## Reading order for this guide

1. [Architecture](01-architecture.md) — ownership and contracts.
2. [Flow](02-flow.md) — startup, demo, collection, close and build.
3. [Structure](03-structure.md) — complete source inventory.
4. [Tech stack](04-tech-stack.md) — runtime, dependencies and tooling.
5. [Decisions](05-decisions.md) — evidence, tradeoffs and gotchas.

## Open questions

- Should the demo use `RemoveRoot` and `Close`, and remove obsolete sleeps? [main.go:32](../../main.go#L32) still suggests background work, while [TriggerGC](../../gc.go#L101) is synchronous.
- Should manually constructed zero-value `Object` handles be supported? Mutators reject them through [valid](../../gc.go#L53), but [object accessors](../../gc.go#L172) dereference a nil owner on a non-nil zero-value object. Current tests only inspect allocated handles.
- Are allocation-plus-rooting transactions needed for concurrent clients? The [documented contract](../../README.md#L17) requires callers to arrange roots before collection; individually locked calls do not make that sequence atomic.
