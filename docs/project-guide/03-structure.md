# Structure

## What lives where

All four Go source files are in the root and share `package main`: the collector is a module inside the executable rather than a separately importable library ([gc.go:1](../../gc.go#L1), [main.go:1](../../main.go#L1)). The v1 contract, usage and history sit alongside them. This inventory covers the 12 tracked files at the guide's baseline plus the existing untracked decision log, with exclusions identified below.

```text
./
├── main.go                example program
├── gc.go                  complete collector and object implementation
├── gc_test.go             basic operation tests
├── lifecycle_test.go      graph and lifecycle regression tests
├── go.mod / go.sum        module and dependency resolution
├── Makefile               local build and test commands
├── README.md / V1.md      usage and delivery contract
├── HISTORY.md             delivery timeline
├── logs/                  development decisions
├── tricolor-gc            generated executable
└── docs/project-guide/    this guide
```

## Root

| File | Responsibility | Key exports or entry points | Called by |
| --- | --- | --- | --- |
| [main.go](../../main.go#L5) | Allocates demo graph and prints before/after one collection. | `main` | Go process startup |
| [gc.go](../../gc.go#L9) | Colors, object and collector state; synchronized graph API; mark/sweep and close. | `Color`, `WHITE`, `GRAY`, `BLACK`, `Object`, `GarbageCollector`, `NewGarbageCollector`; collector methods `Allocate`, `AddRoot`, `RemoveRoot`, `AddReference`, `RemoveReference`, `TriggerGC`, `Close`, `GetHeapSize`, `GetRootCount`, `PrintHeap`; object methods `IsObjectMarked`, `GetRefs` | Demo and both test files |
| [gc_test.go](../../gc_test.go#L10) | Basic initialization, allocation, roots, reference reads, mark/sweep and heap size. | `TestGCInit`, `TestAllocate`, `TestAddRootAndGetRootCount`, `TestAddReferenceAndGetRefs`, `TestTriggerGCAndMarkSweep`, `TestGetHeapSize` | `go test` |
| [lifecycle_test.go](../../lifecycle_test.go#L8) | Reachable/unreachable cycles, deduplication, cleanup/accounting, no resurrection, concurrent graph operations, idempotent close and invalid input. | `TestCyclesAndCleanup`, `TestConcurrentGraphAccess`, `TestForeignAndNegative` | `go test` |
| [go.mod](../../go.mod#L1) | Module identity, Go directive and test dependencies. | `tricolor-gc` module | Go toolchain |
| [Makefile](../../Makefile#L5) | Builds explicit source list, runs verbose tests, removes executable. | `all`, `help`, `build`, `test`, `clean` | Developer |
| [README.md](../../README.md#L1) | Current API example and synchronous-collection limits. | Usage contract | Developer |
| [V1.md](../../V1.md#contract) | v1 scope and acceptance criteria. | Contract, acceptance | Developer |
| [HISTORY.md](../../HISTORY.md#timeline) | Evidence-backed timeline and recorded race-test result. | Timeline, delivery verification | Developer |

## Logs

| File | Responsibility | Key exports | Called by |
| --- | --- | --- | --- |
| [2026-09-21.md](../../logs/2026-09-21.md#L2) | Records bounded v1 scope and synchronous collection decision. | None | Maintainers reviewing rationale |

## Excluded

The generated [tricolor-gc binary](../../tricolor-gc) and [go.sum](../../go.sum) and boilerplate [.gitignore](../../.gitignore) are excluded from design tables. The six files in this guide are navigation artifacts, not runtime source. No vendored sources, infrastructure configuration or CI workflows were present in the inspected tracked inventory.
