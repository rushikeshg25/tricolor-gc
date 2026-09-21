# Tricolor garbage collector simulator v1

This is an educational managed-heap simulator. It tracks objects, roots, references, tricolor marking and unreachable-object cleanup; Go still manages the underlying process memory.

```go
gc := NewGarbageCollector()
defer gc.Close()
a := gc.Allocate(100)
b := gc.Allocate(200)
gc.AddRoot(a)
gc.AddReference(a, b)
gc.TriggerGC() // both retained
gc.RemoveRoot(a)
gc.TriggerGC() // both cleaned up, including cycles
```

`TriggerGC` is synchronous and holds the heap lock through mark and sweep. Collection is explicit: the previous background threshold mechanism could reclaim an object between allocation and establishing its root. Root/reference mutations and accessors synchronize with collection. Callers must establish roots before requesting a cycle; collection and a sequence of separate graph mutations are not one transaction.

Roots are unique. Foreign or already collected objects are ignored. `RemoveReference` removes all matching edges. `GetRefs` returns an owned slice. Cleanup releases payloads and reference slices, updates accounting and prevents resurrection. Negative allocation or allocation after `Close` returns nil. `Close` releases every object and is idempotent. No background goroutines need shutdown.

Run `go run .`; verify with `go test -race ./...`. Tests cover reachable/unreachable cycles, byte accounting, closed heaps and concurrent graph access. Concurrent runtime GC, write barriers, free-list allocation and automatic threshold collection are outside this simulator's v1.
