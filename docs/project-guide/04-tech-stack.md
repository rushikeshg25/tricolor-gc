# Tech stack

## Languages and runtimes

| Language or runtime | Version | Pinned at |
| --- | --- | --- |
| Go | `go 1.23.5` module requirement; exact executing toolchain is not recorded here | [go.mod:3](../../go.mod#L3) |
| Module | `tricolor-gc`, executable `package main` | [go.mod:1](../../go.mod#L1), [main.go:1](../../main.go#L1) |

## Frameworks and major libraries

| Library | Version | Used for | Used in |
| --- | --- | --- | --- |
| Go `sync` | Standard library, follows selected Go toolchain | Collector `RWMutex` and test `WaitGroup` | [gc.go:5](../../gc.go#L5), [lifecycle_test.go:4](../../lifecycle_test.go#L4) |
| Go `time`, `fmt` | Standard library | Cycle elapsed time, demo/test sleeps, stdout statistics | [gc.go:4](../../gc.go#L4), [main.go:3](../../main.go#L3), [gc_test.go:5](../../gc_test.go#L5) |
| Go `testing` | Standard library | Test discovery and assertions through `testing.T` | [lifecycle_test.go:5](../../lifecycle_test.go#L5) |
| `testify/assert` | v1.10.0 | Basic-test assertions only | [go.mod:5](../../go.mod#L5), [gc_test.go:7](../../gc_test.go#L7) |
| `go-spew`, `go-difflib`, `yaml.v3` | v1.1.1, v1.0.0, v3.0.1 respectively | Indirect dependencies declared by the test dependency graph, not runtime integrations | [go.mod:7](../../go.mod#L7) |

## Data and infrastructure

| Service or storage | Role | Configured at |
| --- | --- | --- |
| In-process Go memory | Payload slices, roots, references and statistics | [gc.go:17](../../gc.go#L17) |
| Standard output | Human-readable object/root counts and live/freed payload bytes | [gc.go:188](../../gc.go#L188) |

The executable has no datastore client, network listener or configuration loader ([main.go:5](../../main.go#L5), [gc.go imports](../../gc.go#L3)). YAML is an indirect dependency, not evidence of application YAML configuration.

## Tooling

| Tool | Role | Configured at |
| --- | --- | --- |
| `go run .` | Execute demo | [README:21](../../README.md#L21) |
| Make and `go build` | Compile named implementation files into local executable; Make version unpinned | [Makefile:13](../../Makefile#L13) |
| `go test -v ./...` | Standard Make test target | [Makefile:16](../../Makefile#L16) |
| `go test -race ./...` | Explicit concurrency verification command | [README:21](../../README.md#L21), [recorded result](../../HISTORY.md#delivery-verification) |
| `make clean` | Remove built executable | [Makefile:19](../../Makefile#L19) |

## Notes

No CI, deployment, benchmark or lint configuration was present in the inspected tracked tree. The explicit `go build ... main.go gc.go` command means adding another implementation file also requires updating the Makefile, while `go run .` automatically includes package files ([Makefile:14](../../Makefile#L14)).
