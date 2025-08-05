package main

import (
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
	isMarked bool
	color    Color
	size     uintptr
	refs     []*Object //slice of pointers to other objects
}

type GarbageCollector struct {
	heap           []*Object
	rootSet        []*Object
	freeList       []*Object //list of available (free) memory blocks that can be reused instead of allocating new ones
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

func (gc *GarbageCollector) Allocate(size int) {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()

}
