# Tricolor Garbage Collector

A concurrent tricolor garbage collector implementation in Go, demonstrating the fundamental concepts of automatic memory management used in modern programming languages.

## Overview

This project implements a tricolor marking garbage collector with concurrent marking and sweeping phases. The tricolor algorithm is a foundational technique used in garbage collectors for languages like Go, Java, and JavaScript, providing efficient memory management while minimizing pause times.

## Features

- **Tricolor Marking Algorithm**: Objects are colored WHITE (unmarked), GRAY (marked but not scanned), or BLACK (marked and scanned)
- **Concurrent Collection**: Separate goroutines handle marking and sweeping phases
- **Root Set Management**: Track root objects (stack variables, globals) that serve as GC starting points
- **Reference Tracking**: Objects can reference other objects, creating object graphs
- **Automatic Triggering**: GC automatically runs when memory allocation exceeds threshold
- **Thread-Safe Operations**: All operations are protected with appropriate synchronization
- **Memory Statistics**: Track allocation, deallocation, and GC performance metrics

## How It Works

### Tricolor Algorithm

1. **Initialization**: All objects start as WHITE (unmarked)
2. **Root Marking**: Root objects are colored GRAY (marked but not processed)
3. **Marking Phase**: 
   - GRAY objects are processed and colored BLACK
   - Their references are colored GRAY
   - Process continues until no GRAY objects remain
4. **Sweeping Phase**: All WHITE objects are deallocated as unreachable

### Concurrent Design

- **Marker Goroutine**: Processes objects from the mark queue concurrently
- **Sweeper Goroutine**: Handles cleanup of unreachable objects
- **Queue-Based Communication**: Uses buffered channels for work distribution

## Project Structure

```
tricolor-gc/
├── main.go          # Example usage and demonstration
├── gc.go            # Core garbage collector implementation
├── gc_test.go       # Comprehensive test suite
├── Makefile         # Build and test automation
├── go.mod           # Go module definition
└── README.md        
```

## Getting Started


### Building

```bash
# Using Make
make build

# Or directly with Go
go build -o tricolor-gc main.go gc.go
```

### Running

```bash
# Run the example
./tricolor-gc

# Or directly
go run main.go gc.go
```

### Testing

```bash
# Using Make
make test

# Or directly with Go
go test -v ./...
```

## Usage Example

```go
package main

import "time"

func main() {
    // Create a new garbage collector
    gc := NewGarbageCollector()
    
    // Allocate objects
    obj1 := gc.Allocate(100)
    obj2 := gc.Allocate(200)
    obj3 := gc.Allocate(150)
    
    // Add root objects (these won't be collected)
    gc.AddRoot(obj1)
    gc.AddRoot(obj2)
    
    // Create object references
    gc.AddReference(obj1, obj3) // obj1 references obj3
    
    // Print current heap state
    gc.PrintHeap()
    
    // Trigger garbage collection manually
    gc.TriggerGC()
    time.Sleep(100 * time.Millisecond) // Wait for GC to complete
    
    // Print heap state after GC
    gc.PrintHeap()
}
```


## Testing

The project includes comprehensive tests covering:

- Basic allocation and initialization
- Root set management
- Reference tracking
- Mark and sweep phases
- Concurrent operations
- Memory statistics

Run tests with detailed output:
```bash
go test -v ./...
```

## Performance Considerations

- **Concurrent Design**: Marking and sweeping run in separate goroutines
- **Queue Buffering**: Buffered channels prevent blocking during peak loads
- **Threshold-Based Triggering**: Automatic GC prevents excessive memory usage
- **Read-Write Locks**: Minimize contention during heap access
