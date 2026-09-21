package main

import (
	"fmt"
	"sync"
	"time"
)

type Color int

const (
	WHITE Color = iota
	GRAY
	BLACK
)

type Object struct {
	id       int
	data     []byte
	color    Color
	isMarked bool
	size     uintptr
	refs     []*Object
	owner    *GarbageCollector
	alive    bool
}
type GarbageCollector struct {
	mutex                                sync.RWMutex
	objID                                int
	heap                                 []*Object
	rootSet                              []*Object
	totalAlloc, totalFreed, gcIterations int64
	gcTime                               time.Duration
	closed                               bool
}

func NewGarbageCollector() *GarbageCollector { return &GarbageCollector{} }

// Allocate returns an unrooted object. Negative allocations and closed heaps return nil.
// Collection is explicit, so callers can establish roots before requesting a cycle.
func (gc *GarbageCollector) Allocate(size int) *Object {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()
	if size < 0 || gc.closed {
		return nil
	}
	o := &Object{id: gc.objID, data: make([]byte, size), size: uintptr(size), owner: gc, alive: true}
	gc.objID++
	gc.heap = append(gc.heap, o)
	gc.totalAlloc += int64(size)
	return o
}
func (gc *GarbageCollector) valid(o *Object) bool {
	return o != nil && o.owner == gc && o.alive && !gc.closed
}
func (gc *GarbageCollector) AddRoot(o *Object) {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()
	if !gc.valid(o) {
		return
	}
	for _, r := range gc.rootSet {
		if r == o {
			return
		}
	}
	gc.rootSet = append(gc.rootSet, o)
}
func (gc *GarbageCollector) RemoveRoot(o *Object) {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()
	for i, r := range gc.rootSet {
		if r == o {
			gc.rootSet = append(gc.rootSet[:i], gc.rootSet[i+1:]...)
			return
		}
	}
}
func (gc *GarbageCollector) AddReference(from, to *Object) {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()
	if gc.valid(from) && gc.valid(to) {
		from.refs = append(from.refs, to)
	}
}
func (gc *GarbageCollector) RemoveReference(from, to *Object) {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()
	if !gc.valid(from) {
		return
	}
	refs := from.refs[:0]
	for _, r := range from.refs {
		if r != to {
			refs = append(refs, r)
		}
	}
	from.refs = refs
}

// TriggerGC runs a complete stop-the-world cycle under the heap lock.
func (gc *GarbageCollector) TriggerGC() {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()
	if gc.closed {
		return
	}
	start := time.Now()
	for _, o := range gc.heap {
		o.color = WHITE
		o.isMarked = false
	}
	queue := append([]*Object(nil), gc.rootSet...)
	for _, r := range queue {
		r.color = GRAY
	}
	for len(queue) > 0 {
		o := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		if o.color == BLACK {
			continue
		}
		for _, r := range o.refs {
			if r.alive && r.color == WHITE {
				r.color = GRAY
				queue = append(queue, r)
			}
		}
		o.color = BLACK
		o.isMarked = true
	}
	live := make([]*Object, 0, len(gc.heap))
	for _, o := range gc.heap {
		if o.isMarked {
			live = append(live, o)
		} else {
			gc.totalFreed += int64(o.size)
			gc.totalAlloc -= int64(o.size)
			gc.cleanupObject(o)
		}
	}
	gc.heap = live
	gc.gcIterations++
	gc.gcTime += time.Since(start)
}
func (gc *GarbageCollector) cleanupObject(o *Object) { o.data = nil; o.refs = nil; o.alive = false }
func (gc *GarbageCollector) Close() {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()
	if gc.closed {
		return
	}
	for _, o := range gc.heap {
		gc.cleanupObject(o)
	}
	gc.totalFreed += gc.totalAlloc
	gc.totalAlloc = 0
	gc.heap = nil
	gc.rootSet = nil
	gc.closed = true
}
func (gc *GarbageCollector) GetHeapSize() int {
	gc.mutex.RLock()
	defer gc.mutex.RUnlock()
	return len(gc.heap)
}
func (gc *GarbageCollector) GetRootCount() int {
	gc.mutex.RLock()
	defer gc.mutex.RUnlock()
	return len(gc.rootSet)
}
func (o *Object) IsObjectMarked() bool {
	if o == nil {
		return false
	}
	o.owner.mutex.RLock()
	defer o.owner.mutex.RUnlock()
	return o.isMarked && o.alive
}
func (o *Object) GetRefs() []*Object {
	if o == nil {
		return nil
	}
	o.owner.mutex.RLock()
	defer o.owner.mutex.RUnlock()
	return append([]*Object(nil), o.refs...)
}
func (gc *GarbageCollector) PrintHeap() {
	gc.mutex.RLock()
	defer gc.mutex.RUnlock()
	fmt.Printf("Objects: %d, roots: %d, live bytes: %d, freed bytes: %d\n", len(gc.heap), len(gc.rootSet), gc.totalAlloc, gc.totalFreed)
}
