# Decisions

The implemented behavior below is source-backed. Rationale is labeled as documented or inferred.

## Explicit synchronous collection

- **What:** `TriggerGC` marks and sweeps under one exclusive heap lock ([gc.go:101](../../gc.go#L101)).
- **Why:** The [README:17](../../README.md#L17) explicitly says the old background threshold could collect an object between allocation and rooting. The [decision record](../../logs/2026-09-21.md#L9) also confirms synchronous collection.
- **Tradeoff:** Deterministic completion and a stable graph during tracing, at the cost of blocking all collector operations for the cycle. Separate allocation/rooting calls can still interleave with a collection requested by another caller.
- **Confidence:** Documented rationale, confirmed implementation. No worker, threshold or write barrier exists in v1 ([README:21](../../README.md#L21)).

## Collector-owned live handles

- **What:** Every allocation records owner and liveness; mutations validate these fields, and cleanup prevents resurrection by marking objects dead ([gc.go:47](../../gc.go#L47), [gc.go:53](../../gc.go#L53), [gc.go:146](../../gc.go#L146)).
- **Why, inferred:** Central ownership checks protect the heap from foreign objects and stale handles. Tests specifically reject both cases ([lifecycle_test.go:25](../../lifecycle_test.go#L25), [lifecycle_test.go:57](../../lifecycle_test.go#L57)).
- **Tradeoff:** Invalid requests silently do nothing, so callers receive no diagnostic explaining a rejected root or edge.
- **Confidence:** Behavior tested; rationale inferred.

## Simple slices and iterative graph tracing

- **What:** Heap, roots and references are slices. Marking uses a LIFO slice worklist and tricolor states; sweep builds a fresh survivor slice ([gc.go:17](../../gc.go#L17), [gc.go:113](../../gc.go#L113), [gc.go:132](../../gc.go#L132)).
- **Why, inferred:** Direct data structures keep the algorithm visible, and iteration avoids recursive call-stack growth. The educational scope is explicit in [V1.md](../../V1.md#contract).
- **Tradeoff:** Root deduplication and reference removal scan slices; a collection allocates temporary slices and retains survivor capacity sized to the pre-sweep heap. Duplicate edges are allowed and scanned ([gc.go:62](../../gc.go#L62), [gc.go:79](../../gc.go#L79)).
- **Confidence:** Implementation confirmed; rationale inferred.

## Logical cleanup instead of allocator replacement

- **What:** Sweep and close clear payload/edge slices and update simulated payload-byte counters ([gc.go:132](../../gc.go#L132), [gc.go:147](../../gc.go#L147)).
- **Why:** The [README:3](../../README.md#L3) explicitly distinguishes the simulator from Go's underlying memory management.
- **Tradeoff:** It teaches reachability and lifecycle, but byte counts do not include structs, slice storage or runtime overhead and do not measure resident memory. A retained dead pointer still refers to a Go object struct.
- **Confidence:** Documented design; accounting limitation follows from counting `size` only.

## Gotchas

- **Obsolete background-GC hints:** [main.go:32](../../main.go#L32) and [gc_test.go:55](../../gc_test.go#L55) sleep as though workers exist. Allocation does not trigger collection and `TriggerGC` returns only after sweep finishes ([gc.go:39](../../gc.go#L39), [gc.go:101](../../gc.go#L101)).
- **Demo bypasses the lock:** [main.go:35](../../main.go#L35) slices `rootSet` directly. That single-threaded example works, but concurrent callers must use [RemoveRoot](../../gc.go#L69). The demo also omits [Close](../../gc.go#L147), unlike the README example.
- **An allocated handle is not a root:** Even a handle in a local Go variable is collected unless reachable from `rootSet` at the cycle ([gc.go:113](../../gc.go#L113)).
- **Zero-value object accessors:** Nil handles are handled, but a non-nil `Object{}` has no owner and panics in [IsObjectMarked/GetRefs](../../gc.go#L172). Use objects from `Allocate`.
- **Misleading counter and worklist names:** `totalAlloc` decreases during sweep and means current live payload bytes; `queue` pops from the end and behaves as a stack ([gc.go:118](../../gc.go#L118), [gc.go:138](../../gc.go#L138)).
- **A copied slice is not a copied graph:** Changing the slice returned by `GetRefs` does not edit stored edges, but its entries still reference the same objects ([gc.go:180](../../gc.go#L180)).

## Conventions

Keep graph mutation and lifecycle changes under the collector's write lock, and handle queries under its read lock ([gc.go:56](../../gc.go#L56), [gc.go:162](../../gc.go#L162)). Tests are in the implementation package, so some inspect private fields directly; concurrent tests use methods while workers run and inspect lifecycle results after joining ([lifecycle_test.go:30](../../lifecycle_test.go#L30)). Basic tests use Testify, while lifecycle tests use standard `testing.T` checks ([gc_test.go:7](../../gc_test.go#L7), [lifecycle_test.go:17](../../lifecycle_test.go#L17)).
