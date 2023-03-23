package runtime

import (
	"runtime/internal/atomic"
	"runtime/internal/sys"
)

// Per-thread (in Go, per-P) cache for small objects.
// This includes a small object cache and local allocation stats.
// No locking needed because it is per-thread (per-P).
//
// mcaches are allocated from non-GC'd memory, so any heap pointers
// must be specially handled.
type mcache struct {
	_ sys.NotInHeap

	// The following members are accessed on every malloc,
	// so they are grouped here for better caching.
	nextSample uintptr // trigger heap sample after allocating this many bytes
	scanAlloc  uintptr // bytes of scannable heap allocated

	// Allocator cache for tiny objects w/o pointers.
	// See "Tiny allocator" comment in malloc.go.

	// tiny points to the beginning of the current tiny block, or
	// nil if there is no current tiny block.
	//
	// tiny is a heap pointer. Since mcache is in non-GC'd memory,
	// we handle it by clearing it in releaseAll during mark
	// termination.
	//
	// tinyAllocs is the number of tiny allocations performed
	// by the P that owns this mcache.
	tiny       uintptr
	tinyoffset uintptr
	tinyAllocs uintptr

	// The rest is not accessed on every malloc.

	alloc [8]*mspan // spans to allocate from, indexed by spanClass

	stackcache [8]stackfreelist

	// flushGen indicates the sweepgen during which this mcache
	// was last flushed. If flushGen != mheap_.sweepgen, the spans
	// in this mcache are stale and need to the flushed so they
	// can be swept. This is done in acquirep.
	flushGen atomic.Uint32
}

type gclinkptr uintptr

type stackfreelist struct {
	list gclinkptr // linked list of free stacks
	size uintptr   // total size of stacks in list
}

func (c *mcache) refill(spc spanClass)                        {}
func getMCache(mp *m) *mcache                                 { return nil }
func (c *mcache) prepareForSweep()                            {}
func allocmcache() *mcache                                    { return nil }
func freemcache(c *mcache)                                    {}
func (c *mcache) allocLarge(size uintptr, noscan bool) *mspan { return nil }
