package main

import (
	"sync"
	"testing"
)

func TestCyclesAndCleanup(t *testing.T) {
	g := NewGarbageCollector()
	a := g.Allocate(3)
	b := g.Allocate(4)
	g.AddRoot(a)
	g.AddRoot(a)
	g.AddReference(a, b)
	g.AddReference(b, a)
	g.TriggerGC()
	if g.GetHeapSize() != 2 || g.GetRootCount() != 1 {
		t.Fatal("reachable cycle lost")
	}
	g.RemoveRoot(a)
	g.TriggerGC()
	if g.GetHeapSize() != 0 || a.data != nil || b.refs != nil || g.totalAlloc != 0 || g.totalFreed != 7 {
		t.Fatal("cleanup failed")
	}
	g.AddRoot(a)
	if g.GetRootCount() != 0 {
		t.Fatal("resurrected dead object")
	}
}
func TestConcurrentGraphAccess(t *testing.T) {
	g := NewGarbageCollector()
	a := g.Allocate(1)
	g.AddRoot(a)
	b := g.Allocate(1)
	g.AddRoot(b)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				g.AddReference(a, b)
				a.GetRefs()
				g.RemoveReference(a, b)
				g.TriggerGC()
				a.IsObjectMarked()
			}
		}()
	}
	wg.Wait()
	g.Close()
	g.Close()
	if g.Allocate(1) != nil || g.GetHeapSize() != 0 {
		t.Fatal("closed heap")
	}
}
func TestForeignAndNegative(t *testing.T) {
	g := NewGarbageCollector()
	if g.Allocate(-1) != nil {
		t.Fatal("negative size")
	}
	foreign := NewGarbageCollector().Allocate(1)
	a := g.Allocate(1)
	g.AddRoot(foreign)
	g.AddReference(a, foreign)
	if len(a.GetRefs()) != 0 || g.GetRootCount() != 0 {
		t.Fatal("foreign object accepted")
	}
}
