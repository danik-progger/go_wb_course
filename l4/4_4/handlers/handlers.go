package handlers

import (
	"expvar"
	"fmt"
	"net/http"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"profiler/metrics"
)

var (
	// Global variables to store metrics
	currentMemStats metrics.MemStats
	memStatsMutex   sync.RWMutex

	// For tracking GC duration
	lastNumGC        uint32
	lastPauseTotalNs uint64
)

// InitMetrics initializes the metrics package
func InitMetrics() {
	// Initialize with current stats
	memStatsMutex.Lock()
	currentMemStats = metrics.GetMemStats()
	lastNumGC = currentMemStats.NumGC
	lastPauseTotalNs = 0
	memStatsMutex.Unlock()
}

// CollectMetrics periodically collects metrics
func CollectMetrics() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		stats := metrics.GetMemStats()
		memStatsMutex.Lock()
		currentMemStats = stats
		memStatsMutex.Unlock()
	}
}

// MetricsHandler serves GC and memory metrics in Prometheus format
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	memStatsMutex.RLock()
	stats := currentMemStats
	memStatsMutex.RUnlock()

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")

	// Memory allocation metrics
	fmt.Fprintf(w, "# HELP go_memstats_alloc_bytes Bytes allocated and not yet freed\n")
	fmt.Fprintf(w, "# TYPE go_memstats_alloc_bytes gauge\n")
	fmt.Fprintf(w, "go_memstats_alloc_bytes %d\n", stats.Alloc)

	fmt.Fprintf(w, "# HELP go_memstats_total_alloc_bytes_total Bytes allocated (even if freed)\n")
	fmt.Fprintf(w, "# TYPE go_memstats_total_alloc_bytes_total counter\n")
	fmt.Fprintf(w, "go_memstats_total_alloc_bytes_total %d\n", stats.TotalAlloc)

	fmt.Fprintf(w, "# HELP go_memstats_sys_bytes Bytes obtained from system\n")
	fmt.Fprintf(w, "# TYPE go_memstats_sys_bytes gauge\n")
	fmt.Fprintf(w, "go_memstats_sys_bytes %d\n", stats.Sys)

	fmt.Fprintf(w, "# HELP go_memstats_lookups_total Number of pointer lookups\n")
	fmt.Fprintf(w, "# TYPE go_memstats_lookups_total counter\n")
	fmt.Fprintf(w, "go_memstats_lookups_total %d\n", stats.Lookups)

	fmt.Fprintf(w, "# HELP go_memstats_mallocs_total Number of mallocs\n")
	fmt.Fprintf(w, "# TYPE go_memstats_mallocs_total counter\n")
	fmt.Fprintf(w, "go_memstats_mallocs_total %d\n", stats.Mallocs)

	fmt.Fprintf(w, "# HELP go_memstats_frees_total Number of frees\n")
	fmt.Fprintf(w, "# TYPE go_memstats_frees_total counter\n")
	fmt.Fprintf(w, "go_memstats_frees_total %d\n", stats.Frees)

	// Heap metrics
	fmt.Fprintf(w, "# HELP go_memstats_heap_alloc_bytes Bytes allocated and not yet freed (heap)\n")
	fmt.Fprintf(w, "# TYPE go_memstats_heap_alloc_bytes gauge\n")
	fmt.Fprintf(w, "go_memstats_heap_alloc_bytes %d\n", stats.HeapAlloc)

	fmt.Fprintf(w, "# HELP go_memstats_heap_sys_bytes Bytes obtained from system for heap\n")
	fmt.Fprintf(w, "# TYPE go_memstats_heap_sys_bytes gauge\n")
	fmt.Fprintf(w, "go_memstats_heap_sys_bytes %d\n", stats.HeapSys)

	fmt.Fprintf(w, "# HELP go_memstats_heap_idle_bytes Bytes in idle spans\n")
	fmt.Fprintf(w, "# TYPE go_memstats_heap_idle_bytes gauge\n")
	fmt.Fprintf(w, "go_memstats_heap_idle_bytes %d\n", stats.HeapIdle)

	fmt.Fprintf(w, "# HELP go_memstats_heap_inuse_bytes Bytes in non-idle spans\n")
	fmt.Fprintf(w, "# TYPE go_memstats_heap_inuse_bytes gauge\n")
	fmt.Fprintf(w, "go_memstats_heap_inuse_bytes %d\n", stats.HeapInuse)

	fmt.Fprintf(w, "# HELP go_memstats_heap_released_bytes Bytes released to the OS\n")
	fmt.Fprintf(w, "# TYPE go_memstats_heap_released_bytes gauge\n")
	fmt.Fprintf(w, "go_memstats_heap_released_bytes %d\n", stats.HeapReleased)

	fmt.Fprintf(w, "# HELP go_memstats_heap_objects Number of allocated heap objects\n")
	fmt.Fprintf(w, "# TYPE go_memstats_heap_objects gauge\n")
	fmt.Fprintf(w, "go_memstats_heap_objects %d\n", stats.HeapObjects)

	// Stack metrics
	fmt.Fprintf(w, "# HELP go_memstats_stack_inuse_bytes Bytes used by stack allocator\n")
	fmt.Fprintf(w, "# TYPE go_memstats_stack_inuse_bytes gauge\n")
	fmt.Fprintf(w, "go_memstats_stack_inuse_bytes %d\n", stats.StackInuse)

	fmt.Fprintf(w, "# HELP go_memstats_stack_sys_bytes Bytes used for stacks obtained from system\n")
	fmt.Fprintf(w, "# TYPE go_memstats_stack_sys_bytes gauge\n")
	fmt.Fprintf(w, "go_memstats_stack_sys_bytes %d\n", stats.StackSys)

	fmt.Fprintf(w, "# HELP go_memstats_mspan_inuse_bytes Bytes used by mspan structures\n")
	fmt.Fprintf(w, "# TYPE go_memstats_mspan_inuse_bytes gauge\n")
	fmt.Fprintf(w, "go_memstats_mspan_inuse_bytes %d\n", stats.MSpanInuse)

	fmt.Fprintf(w, "# HELP go_memstats_mspan_sys_bytes Bytes used for mspan structures obtained from system\n")
	fmt.Fprintf(w, "# TYPE go_memstats_mspan_sys_bytes gauge\n")
	fmt.Fprintf(w, "go_memstats_mspan_sys_bytes %d\n", stats.MSpanSys)

	fmt.Fprintf(w, "# HELP go_memstats_mcache_inuse_bytes Bytes used by mcache structures\n")
	fmt.Fprintf(w, "# TYPE go_memstats_mcache_inuse_bytes gauge\n")
	fmt.Fprintf(w, "go_memstats_mcache_inuse_bytes %d\n", stats.MCacheInuse)

	fmt.Fprintf(w, "# HELP go_memstats_mcache_sys_bytes Bytes used for mcache structures obtained from system\n")
	fmt.Fprintf(w, "# TYPE go_memstats_mcache_sys_bytes gauge\n")
	fmt.Fprintf(w, "go_memstats_mcache_sys_bytes %d\n", stats.MCacheSys)

	// GC metrics
	fmt.Fprintf(w, "# HELP go_gc_cycles_total Completed GC cycles\n")
	fmt.Fprintf(w, "# TYPE go_gc_cycles_total counter\n")
	fmt.Fprintf(w, "go_gc_cycles_total %d\n", stats.NumGC)

	fmt.Fprintf(w, "# HELP go_memstats_gc_cpu_fraction Fraction of CPU time used by GC\n")
	fmt.Fprintf(w, "# TYPE go_memstats_gc_cpu_fraction gauge\n")
	fmt.Fprintf(w, "go_memstats_gc_cpu_fraction %f\n", stats.GCCPUFraction)

	fmt.Fprintf(w, "# HELP go_gc_last_time_seconds Time of last garbage collection\n")
	fmt.Fprintf(w, "# TYPE go_gc_last_time_seconds gauge\n")
	fmt.Fprintf(w, "go_gc_last_time_seconds %f\n", float64(stats.LastGC)/1e9)

	// Calculate GC pause time if possible
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	pauseTotalNs := m.PauseTotalNs
	totalGCDurationSecs := float64(pauseTotalNs) / 1e9

	fmt.Fprintf(w, "# HELP go_gc_duration_seconds_total Total time spent in GC\n")
	fmt.Fprintf(w, "# TYPE go_gc_duration_seconds_total counter\n")
	fmt.Fprintf(w, "go_gc_duration_seconds_total %f\n", totalGCDurationSecs)

	// Goroutine count
	goroutines := runtime.NumGoroutine()
	fmt.Fprintf(w, "# HELP go_goroutines Number of goroutines\n")
	fmt.Fprintf(w, "# TYPE go_goroutines gauge\n")
	fmt.Fprintf(w, "go_goroutines %d\n", goroutines)

	// CGO calls
	var cgoCalls int64
	if debugMetrics := debugMetrics(); debugMetrics != nil {
		if cgoCallsVal, ok := debugMetrics["cgo_calls"]; ok {
			if cgoCallsStr, ok := cgoCallsVal.(string); ok {
				if val, err := strconv.ParseInt(cgoCallsStr, 10, 64); err == nil {
					cgoCalls = val
				}
			}
		}
	}
	fmt.Fprintf(w, "# HELP go_cgo_calls_total Number of cgo calls\n")
	fmt.Fprintf(w, "# TYPE go_cgo_calls_total counter\n")
	fmt.Fprintf(w, "go_cgo_calls_total %d\n", cgoCalls)

	// Memory stats from expvar if available
	if expvarMetrics := expvarMetrics(); expvarMetrics != nil {
		for _, line := range expvarMetrics {
			if strings.HasPrefix(line, "memstats:") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					if key == "memstats" {
						continue // Skip the header line
					}
					// Format would be like: memstats:Alloc=123456
					kv := strings.Split(value, "=")
					if len(kv) == 2 {
						fmt.Fprintf(w, "# HELP go_expvar_%s Generated from expvar\n", kv[0])
						fmt.Fprintf(w, "# TYPE go_expvar_%s gauge\n", kv[0])
						fmt.Fprintf(w, "go_expvar_%s %s\n", kv[0], kv[1])
					}
				}
			}
		}
	}
}

// HealthHandler provides a simple health check
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, "OK")
}

// AllocateMemoryHandler creates memory allocations for testing
func AllocateMemoryHandler(w http.ResponseWriter, r *http.Request) {
	// Allocate some memory to demonstrate the metrics
	size := r.URL.Query().Get("size")
	if size == "" {
		size = "1048576" // 1MB default
	}

	bytes, err := strconv.Atoi(size)
	if err != nil || bytes <= 0 {
		http.Error(w, "Invalid size parameter", http.StatusBadRequest)
		return
	}

	// Allocate memory
	data := make([]byte, bytes)

	// Fill with data to prevent optimization
	for i := range data {
		data[i] = byte(i % 256)
	}

	// Optionally trigger GC
	if r.URL.Query().Get("gc") == "true" {
		runtime.GC()
	}

	fmt.Fprintf(w, "Allocated %d bytes\n", len(data))
}

// debugMetrics returns debug metrics
func debugMetrics() map[string]interface{} {
	metrics := make(map[string]interface{})

	expvar.Do(func(kv expvar.KeyValue) {
		metrics[kv.Key] = kv.Value
	})

	return metrics
}

// expvarMetrics returns expvar metrics as strings
func expvarMetrics() []string {
	var lines []string
	expvar.Do(func(kv expvar.KeyValue) {
		lines = append(lines, fmt.Sprintf("%s: %v", kv.Key, kv.Value))
	})
	sort.Strings(lines)
	return lines
}
