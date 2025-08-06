package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGCInit(t *testing.T) {
	gc := NewGarbageCollector()
	if gc == nil {
		t.Error("Garbage collector not initialized")
	}
	assert.Equal(t, len(gc.heap), 0, "Heap in GC is initialized")
}

func TestAllocate(t *testing.T) {
	gc := NewGarbageCollector()
	obj := gc.Allocate(10)
	assert.Equal(t, obj.size, uintptr(10), "Object size is not correct")
	assert.Equal(t, len(gc.heap), 1, "Heap size is not correct")
}

func TestAddRootAndGetRootCount(t *testing.T) {
	gc := NewGarbageCollector()
	obj1 := gc.Allocate(8)
	obj2 := gc.Allocate(16)
	gc.AddRoot(obj1)
	gc.AddRoot(obj2)
	assert.Equal(t, 2, gc.GetRootCount(), "Root count should be 2 after adding two roots")
}

func TestAddReferenceAndGetRefs(t *testing.T) {
	gc := NewGarbageCollector()
	obj1 := gc.Allocate(8)
	obj2 := gc.Allocate(16)
	gc.AddReference(obj1, obj2)
	references := obj1.GetRefs()
	assert.Equal(t, 1, len(references), "Object should have one reference after AddReference")
	assert.Equal(t, obj2, references[0], "Referenced object should be obj2")
}

func TestTriggerGCAndMarkSweep(t *testing.T) {
	gc := NewGarbageCollector()
	obj1 := gc.Allocate(8)
	obj2 := gc.Allocate(16)
	obj3 := gc.Allocate(32)
	gc.AddRoot(obj1)
	gc.AddReference(obj1, obj2)
	// obj3 is unreachable

	// Mark obj1 and obj2 as roots and references
	gc.TriggerGC()
	// Wait a bit for GC goroutines to finish
	time.Sleep(100 * time.Millisecond)

	assert.True(t, obj1.IsObjectMarked(), "Root object should be marked after GC")
	assert.True(t, obj2.IsObjectMarked(), "Referenced object should be marked after GC")
	assert.False(t, obj3.IsObjectMarked(), "Unreachable object should not be marked after GC")
	assert.Equal(t, 2, gc.GetHeapSize(), "Heap should only contain reachable objects after GC")
}

func TestGetHeapSize(t *testing.T) {
	gc := NewGarbageCollector()
	assert.Equal(t, 0, gc.GetHeapSize(), "Heap size should be 0 initially")
	obj1 := gc.Allocate(8)
	gc.Allocate(16)
	assert.Equal(t, 2, gc.GetHeapSize(), "Heap size should be 2 after two allocations")
	gc.AddRoot(obj1)
	gc.TriggerGC()
	time.Sleep(100 * time.Millisecond)
	// obj2 is not a root, should be collected
	assert.Equal(t, 1, gc.GetHeapSize(), "Heap size should be 1 after GC collects unreachable object")
}
