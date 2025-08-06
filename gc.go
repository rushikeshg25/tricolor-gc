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
		threshold:  512, //in bytes
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

func (gc *GarbageCollector) markObject(obj *Object) {
	if obj == nil || obj.color == BLACK {
		return
	}
	obj.color = BLACK
	obj.isMarked = true

	for _, ref := range obj.refs {
		if ref != nil && ref.color == WHITE {
			ref.color = GRAY
			gc.wg.Add(1)
			select {
			case gc.markQueue <- ref:
			default:
				gc.markObject(ref)
				gc.wg.Done()
			}
		}
	}
}

func (gc *GarbageCollector) concurrentSweeper() {
	for obj := range gc.sweepQueue {
		gc.cleanupObject(obj)
	}
}
func (gc *GarbageCollector) cleanupObject(obj *Object) {}

func (gc *GarbageCollector) Allocate(size int) *Object {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()

	// Always create new objects to avoid memory corruption issues
	// In a production GC, you'd implement a more sophisticated free list
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
	gc.markWithColorPhase()

	// Phase 2: Sweep phase
	gc.sweepPhase()

	duration := time.Since(start)
	gc.gcTime += duration

	atomic.StoreInt32(&gc.gcRunningState, 0)

}

func (gc *GarbageCollector) markWithColorPhase() {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()

	for _, obj := range gc.heap {
		obj.color = WHITE
		obj.isMarked = false
	}

	for _, root := range gc.rootSet {
		if root != nil {
			root.color = GRAY
			gc.wg.Add(1)
			select {
			case gc.markQueue <- root:
			default:
				gc.markObject(root)
				gc.wg.Done()
			}
		}
	}
	gc.wg.Wait()
}

func (gc *GarbageCollector) sweepPhase() {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()
	var newHeap []*Object
	var freedCount int
	var freedBytes int64
	fmt.Printf("Before sweep - Roots: %d, Heap objects: %d\n", len(gc.rootSet), len(gc.heap))
	for i, root := range gc.rootSet {
		if root != nil {
			fmt.Printf("  Root[%d]: Obj[%d], marked: %v\n", i, root.id, root.isMarked)
		}
	}

	var newRootSet []*Object

	for _, obj := range gc.heap {
		if obj.isMarked {
			// Keep marked objects
			newHeap = append(newHeap, obj)
		} else {
			// Free unmarked objects
			freedCount++
			freedBytes += int64(obj.size)

			select {
			case gc.sweepQueue <- obj:
			default:
				gc.cleanupObject(obj)
			}
		}
	}

	for _, root := range gc.rootSet {
		if root != nil && root.isMarked {
			keepRoot := false
			for _, heapObj := range newHeap {
				if heapObj == root {
					keepRoot = true
					break
				}
			}
			if keepRoot {
				newRootSet = append(newRootSet, root)
			}
		}
	}

	gc.heap = newHeap
	gc.rootSet = newRootSet
	atomic.AddInt64(&gc.totalFreed, freedBytes)
	atomic.AddInt64(&gc.totalAlloc, -freedBytes)

	fmt.Printf("Swept %d objects, freed %d bytes\n", freedCount, freedBytes)
	fmt.Printf("Roots after sweep: %d\n", len(gc.rootSet))
}

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

func (gc *GarbageCollector) PrintHeap() {
	gc.mutex.RLock()
	defer gc.mutex.RUnlock()

	fmt.Println("Heap State:")
	fmt.Printf("Objects: %d, Roots: %d\n", len(gc.heap), len(gc.rootSet))

	colorNames := map[Color]string{WHITE: "White", GRAY: "Gray", BLACK: "Black"}

	for i, obj := range gc.heap {
		if i < 10 {
			fmt.Printf("  Obj[%d]: %s, Size: %d, Refs: %d\n",
				obj.id, colorNames[obj.color], obj.size, len(obj.refs))
		}
	}

	if len(gc.heap) > 10 {
		fmt.Printf("  ... and %d more objects\n", len(gc.heap)-10)
	}
}
