package runtime

import "runtime/internal/atomic"

var scavenge struct {
	// gcPercentGoal is the amount of retained heap memory (measured by
	// heapRetained) that the runtime will try to maintain by returning
	// memory to the OS. This goal is derived from gcController.gcPercent
	// by choosing to retain enough memory to allocate heap memory up to
	// the heap goal.
	gcPercentGoal atomic.Uint64

	// memoryLimitGoal is the amount of memory retained by the runtime (
	// measured by gcController.mappedReady) that the runtime will try to
	// maintain by returning memory to the OS. This goal is derived from
	// gcController.memoryLimit by choosing to target the memory limit or
	// some lower target to keep the scavenger working.
	memoryLimitGoal atomic.Uint64

	// assistTime is the time spent by the allocator scavenging in the last GC cycle.
	//
	// This is reset once a GC cycle ends.
	assistTime atomic.Int64

	// backgroundTime is the time spent by the background scavenger in the last GC cycle.
	//
	// This is reset once a GC cycle ends.
	backgroundTime atomic.Int64
}

type scavengeIndex struct {
	// chunks is a bitmap representing the entire address space. Each bit represents
	// a single chunk, and a 1 value indicates the presence of pages available for
	// scavenging. Updates to the bitmap are serialized by the pageAlloc lock.
	//
	// The underlying storage of chunks is platform dependent and may not even be
	// totally mapped read/write. min and max reflect the extent that is safe to access.
	// min is inclusive, max is exclusive.
	//
	// searchAddr is the maximum address (in the offset address space, so we have a linear
	// view of the address space; see mranges.go:offAddr) containing memory available to
	// scavenge. It is a hint to the find operation to avoid O(n^2) behavior in repeated lookups.
	//
	// searchAddr is always inclusive and should be the base address of the highest runtime
	// page available for scavenging.
	//
	// searchAddr is managed by both find and mark.
	//
	// Normally, find monotonically decreases searchAddr as it finds no more free pages to
	// scavenge. However, mark, when marking a new chunk at an index greater than the current
	// searchAddr, sets searchAddr to the *negative* index into chunks of that page. The trick here
	// is that concurrent calls to find will fail to monotonically decrease searchAddr, and so they
	// won't barge over new memory becoming available to scavenge. Furthermore, this ensures
	// that some future caller of find *must* observe the new high index. That caller
	// (or any other racing with it), then makes searchAddr positive before continuing, bringing
	// us back to our monotonically decreasing steady-state.
	//
	// A pageAlloc lock serializes updates between min, max, and searchAddr, so abs(searchAddr)
	// is always guaranteed to be >= min and < max (converted to heap addresses).
	//
	// TODO(mknyszek): Ideally we would use something bigger than a uint8 for faster
	// iteration like uint32, but we lack the bit twiddling intrinsics. We'd need to either
	// copy them from math/bits or fix the fact that we can't import math/bits' code from
	// the runtime due to compiler instrumentation.
	chunks     []atomic.Uint8
	minHeapIdx atomic.Int32
	min, max   atomic.Int32
}

func bgscavenge(c chan int)                                                  {}
func (p *pageAlloc) scavenge(nbytes uintptr, shouldStop func() bool) uintptr { return 0 }
func gcPaceScavenger(memoryLimit int64, heapGoal, lastHeapGoal uint64)       {}
func heapRetained() uint64                                                   { return 0 }
func printScavTrace(released uintptr, forced bool)                           {}

var scavenger scavengerState

func (s *scavengeIndex) mark(base, limit uintptr) {}

type scavengerState struct {
	// lock protects all fields below.
	lock mutex

	// g is the goroutine the scavenger is bound to.
	g *g

	// parked is whether or not the scavenger is parked.
	parked bool

	// timer is the timer used for the scavenger to sleep.
	timer *timer

	// sysmonWake signals to sysmon that it should wake the scavenger.
	sysmonWake atomic.Uint32

	// targetCPUFraction is the target CPU overhead for the scavenger.
	targetCPUFraction float64

	// sleepRatio is the ratio of time spent doing scavenging work to
	// time spent sleeping. This is used to decide how long the scavenger
	// should sleep for in between batches of work. It is set by
	// critSleepController in order to maintain a CPU overhead of
	// targetCPUFraction.
	//
	// Lower means more sleep, higher means more aggressive scavenging.
	sleepRatio float64

	// sleepController controls sleepRatio.
	//
	// See sleepRatio for more details.
	//sleepController piController

	// cooldown is the time left in nanoseconds during which we avoid
	// using the controller and we hold sleepRatio at a conservative
	// value. Used if the controller's assumptions fail to hold.
	controllerCooldown int64

	// printControllerReset instructs printScavTrace to signal that
	// the controller was reset.
	printControllerReset bool

	// sleepStub is a stub used for testing to avoid actually having
	// the scavenger sleep.
	//
	// Unlike the other stubs, this is not populated if left nil
	// Instead, it is called when non-nil because any valid implementation
	// of this function basically requires closing over this scavenger
	// state, and allocating a closure is not allowed in the runtime as
	// a matter of policy.
	sleepStub func(n int64) int64

	// scavenge is a function that scavenges n bytes of memory.
	// Returns how many bytes of memory it actually scavenged, as
	// well as the time it took in nanoseconds. Usually mheap.pages.scavenge
	// with nanotime called around it, but stubbed out for testing.
	// Like mheap.pages.scavenge, if it scavenges less than n bytes of
	// memory, the caller may assume the heap is exhausted of scavengable
	// memory for now.
	//
	// If this is nil, it is populated with the real thing in init.
	scavenge func(n uintptr) (uintptr, int64)

	// shouldStop is a callback called in the work loop and provides a
	// point that can force the scavenger to stop early, for example because
	// the scavenge policy dictates too much has been scavenged already.
	//
	// If this is nil, it is populated with the real thing in init.
	shouldStop func() bool

	// gomaxprocs returns the current value of gomaxprocs. Stub for testing.
	//
	// If this is nil, it is populated with the real thing in init.
	gomaxprocs func() int32
}

func (s *scavengerState) wake() {}
