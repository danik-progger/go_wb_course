package metrics

import (
	"runtime"
)

// MemStats holds memory statistics
type MemStats struct {
	// General memory statistics
	Alloc      uint64 // bytes allocated and not yet freed
	TotalAlloc uint64 // bytes allocated (even if freed)
	Sys        uint64 // bytes obtained from system
	Lookups    uint64 // number of pointer lookups
	Mallocs    uint64 // number of mallocs
	Frees      uint64 // number of frees

	// Heap memory statistics
	HeapAlloc    uint64 // bytes allocated and not yet freed (same as Alloc but specifically for heap)
	HeapSys      uint64 // bytes obtained from system for heap
	HeapIdle     uint64 // bytes in idle spans
	HeapInuse    uint64 // bytes in non-idle span
	HeapReleased uint64 // bytes released to the OS
	HeapObjects  uint64 // number of allocated objects

	// Stack memory statistics
	StackInuse  uint64 // bytes used by stack allocator
	StackSys    uint64 // bytes obtained from system for stacks
	MSpanInuse  uint64 // bytes used by mspan structures
	MSpanSys    uint64 // bytes used by mspan structures obtained from system
	MCacheInuse uint64 // bytes used by mcache structures
	MCacheSys   uint64 // bytes used by mcache structures obtained from system

	// Garbage collector statistics
	NumGC         uint32  // number of completed GC cycles
	GCCPUFraction float64 // fraction of CPU time used by GC
	EnableGC      bool    // whether GC is enabled
	DebugGC       bool    // whether GC debug info is being printed

	// Last GC time in nanoseconds since Unix epoch
	LastGC uint64 // time of last collection (nanoseconds since 1970)
}

// GetMemStats returns current memory statistics
func GetMemStats() MemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return MemStats{
		Alloc:         m.Alloc,
		TotalAlloc:    m.TotalAlloc,
		Sys:           m.Sys,
		Lookups:       m.Lookups,
		Mallocs:       m.Mallocs,
		Frees:         m.Frees,
		HeapAlloc:     m.HeapAlloc,
		HeapSys:       m.HeapSys,
		HeapIdle:      m.HeapIdle,
		HeapInuse:     m.HeapInuse,
		HeapReleased:  m.HeapReleased,
		HeapObjects:   m.HeapObjects,
		StackInuse:    m.StackInuse,
		StackSys:      m.StackSys,
		MSpanInuse:    m.MSpanInuse,
		MSpanSys:      m.MSpanSys,
		MCacheInuse:   m.MCacheInuse,
		MCacheSys:     m.MCacheSys,
		NumGC:         m.NumGC,
		GCCPUFraction: m.GCCPUFraction,
		EnableGC:      m.EnableGC,
		DebugGC:       m.DebugGC,
		LastGC:        m.LastGC,
	}
}
