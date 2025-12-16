# Go GC and Memory Profiler

This application exposes Go runtime metrics including garbage collector and memory statistics in Prometheus format.

## Features

- Exposes detailed memory and GC metrics in Prometheus format
- Provides HTTP endpoints for metrics, health checks, and profiling
- Supports runtime GC parameter adjustment
- Includes memory allocation endpoint for testing

## Endpoints

- `/metrics` - Prometheus format metrics for GC and memory
- `/health` - Health check endpoint
- `/debug/pprof/` - Go profiling endpoints (index, heap, goroutine, etc.)
- `/allocate` - Test endpoint to allocate memory

## Installation and Running

1. Build the application:
```bash
go build -o profiler main.go
```

2. Run the application:
```bash
./profiler --port=8080 --gcpercent=100
```

## Available Flags

- `--port` (default: 8080): Port to run the server on
- `--gcpercent` (default: 100): GOGC value (percentage of heap growth that triggers GC)

## Example Metrics Output

The `/metrics` endpoint returns metrics in Prometheus format:

```
# HELP go_memstats_alloc_bytes Bytes allocated and not yet freed
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 123456

# HELP go_gc_cycles_total Completed GC cycles
# TYPE go_gc_cycles_total counter
go_gc_cycles_total 5

# HELP go_memstats_heap_objects Number of allocated heap objects
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 7890
```

## Example Usage

1. Start the server:
```bash
./profiler
```

2. View metrics:
```bash
curl http://localhost:8080/metrics
```

3. Allocate memory for testing:
```bash
curl "http://localhost:8080/allocate?size=2097152"  # Allocates 2MB
```

4. Trigger GC after allocation:
```bash
curl "http://localhost:8080/allocate?size=2097152&gc=true"
```

5. View profiling information:
```bash
# View pprof index
curl http://localhost:8080/debug/pprof/

# Get heap profile
go tool pprof http://localhost:8080/debug/pprof/heap

# Get goroutine profile
go tool pprof http://localhost:8080/debug/pprof/goroutine
```

## Metrics Collected

### Memory Metrics
- `go_memstats_alloc_bytes`: Bytes allocated and not yet freed
- `go_memstats_total_alloc_bytes_total`: Total bytes allocated (even if freed)
- `go_memstats_sys_bytes`: Bytes obtained from system
- `go_memstats_lookups_total`: Number of pointer lookups
- `go_memstats_mallocs_total`: Number of mallocs
- `go_memstats_frees_total`: Number of frees
- `go_memstats_heap_alloc_bytes`: Heap bytes allocated and not yet freed
- `go_memstats_heap_sys_bytes`: Heap bytes obtained from system
- `go_memstats_heap_idle_bytes`: Bytes in idle spans
- `go_memstats_heap_inuse_bytes`: Bytes in non-idle spans
- `go_memstats_heap_released_bytes`: Bytes released to the OS
- `go_memstats_heap_objects`: Number of allocated heap objects
- `go_memstats_stack_inuse_bytes`: Bytes used by stack allocator
- `go_memstats_stack_sys_bytes`: Bytes used for stacks obtained from system
- `go_memstats_mspan_inuse_bytes`: Bytes used by mspan structures
- `go_memstats_mspan_sys_bytes`: Bytes used for mspan structures obtained from system
- `go_memstats_mcache_inuse_bytes`: Bytes used by mcache structures
- `go_memstats_mcache_sys_bytes`: Bytes used for mcache structures obtained from system

### GC Metrics
- `go_gc_cycles_total`: Completed GC cycles
- `go_memstats_gc_cpu_fraction`: Fraction of CPU time used by GC
- `go_gc_last_time_seconds`: Time of last garbage collection
- `go_gc_duration_seconds_total`: Total time spent in GC

### Runtime Metrics
- `go_goroutines`: Number of goroutines
- `go_cgo_calls_total`: Number of cgo calls