package main

import (
	"fmt"
	"sync"
	"sync/atomic"
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
	isMarked bool
	color    Color
	size     uintptr
	refs     []*Object //slice of pointers to other objects
	next     *Object
}

type GarbageCollector struct {
	objID          int
	heap           []*Object
	rootSet        []*Object
	freeList       *Object //list of available (free) memory blocks that can be reused instead of allocating new ones
	totalAlloc     int64
	threshold      int64
	gcRunningState int32
	markQueue      chan *Object
	sweepQueue     chan *Object
	wg             sync.WaitGroup
	mutex          sync.RWMutex
	gcIterations   int64
	totalFreed     int64
	gcTime         time.Duration
}

func NewGarbageCollector() *GarbageCollector {
	g := &GarbageCollector{
		objID:      0,
		heap:       make([]*Object, 0),
		rootSet:    make([]*Object, 0),
		threshold:  1024 * 1024,
		markQueue:  make(chan *Object, 100),
		sweepQueue: make(chan *Object, 100),
	}

	go g.concurrentMarker()
	go g.concurrentSweeper()
	return g
}

func (gc *GarbageCollector) concurrentMarker() {
	for obj := range gc.markQueue {
		gc.markObject(obj)
		gc.wg.Done()
	}
}

func (gc *GarbageCollector) markObject(obj *Object) {}

func (gc *GarbageCollector) concurrentSweeper() {
	for obj := range gc.sweepQueue {
		gc.cleanupObject(obj)
	}
}
func (gc *GarbageCollector) cleanupObject(obj *Object) {}

func (gc *GarbageCollector) Allocate(size int) *Object {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()

	if gc.freeList != nil && int(gc.freeList.size) >= size {
		obj := gc.freeList
		gc.freeList = obj.next
		obj.next = nil
		obj.color = WHITE
		obj.isMarked = false
		obj.data = make([]byte, size)
		return obj
	}
	obj := &Object{
		id:    gc.objID,
		data:  make([]byte, size),
		refs:  make([]*Object, 0),
		color: WHITE,
		size:  uintptr(size),
	}
	gc.objID++
	gc.heap = append(gc.heap, obj)
	atomic.AddInt64(&gc.totalAlloc, int64(size))

	if atomic.LoadInt64(&gc.totalAlloc) > gc.threshold {
		go gc.TriggerGC()
	}
	return obj
}

func (gc *GarbageCollector) TriggerGC() {
	//Check if GC is already running
	if !atomic.CompareAndSwapInt32(&gc.gcRunningState, 0, 1) {
		return
	}

	start := time.Now()
	fmt.Printf("Started Gc at %q", start)

	// Phase 1: Mark phase (tricolor marking)
	gc.markWithTricolorPhase()

	// Phase 2: Sweep phase
	gc.sweepPhase()

	duration := time.Since(start)
	gc.gcTime += duration

	atomic.StoreInt32(&gc.gcRunningState, 0)

}

func (gc *GarbageCollector) markWithTricolorPhase() {}

func (gc *GarbageCollector) sweepPhase() {}

func (gc *GarbageCollector) AddRoot(obj *Object) {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()
	gc.rootSet = append(gc.rootSet, obj)
}

func (gc *GarbageCollector) AddReference(from, to *Object) {
	if from != nil && to != nil {
		from.refs = append(from.refs, to)
	}
}
