package runtime

import (
	"internal/cpu"
	"runtime/internal/atomic"
	"runtime/internal/sys"
	"unsafe"
)

type gcBits struct {
	_ sys.NotInHeap
	x uint8
}

type special struct {
	_      sys.NotInHeap
	next   *special // linked list in span
	offset uint16   // span offset of object
	kind   byte     // kind of special
}

const (
	// minPhysPageSize is a lower-bound on the physical page size. The
	// true physical page size may be larger than this. In contrast,
	// sys.PhysPageSize is an upper-bound on the physical page size.
	minPhysPageSize = 4096

	// maxPhysPageSize is the maximum page size the runtime supports.
	maxPhysPageSize = 512 << 10

	// maxPhysHugePageSize sets an upper-bound on the maximum huge page size
	// that the runtime supports.
	maxPhysHugePageSize = pallocChunkBytes

	// pagesPerReclaimerChunk indicates how many pages to scan from the
	// pageInUse bitmap at a time. Used by the page reclaimer.
	//
	// Higher values reduce contention on scanning indexes (such as
	// h.reclaimIndex), but increase the minimum latency of the
	// operation.
	//
	// The time required to scan this many pages can vary a lot depending
	// on how many spans are actually freed. Experimentally, it can
	// scan for pages at ~300 GB/ms on a 2.6GHz Core i7, but can only
	// free spans at ~32 MB/ms. Using 512 pages bounds this at
	// roughly 100µs.
	//
	// Must be a multiple of the pageInUse bitmap element size and
	// must also evenly divide pagesPerArena.
	pagesPerReclaimerChunk = 512

	// physPageAlignedStacks indicates whether stack allocations must be
	// physical page aligned. This is a requirement for MAP_STACK on
	// OpenBSD.
	physPageAlignedStacks = false
)

var mheap_ mheap

type mspan struct {
	_    sys.NotInHeap
	next *mspan     // next span in list, or nil if none
	prev *mspan     // previous span in list, or nil if none
	list *mSpanList // For debugging. TODO: Remove.

	startAddr uintptr // address of first byte of span aka s.base()
	npages    uintptr // number of pages in span

	manualFreeList gclinkptr // list of free objects in mSpanManual spans

	// freeindex is the slot index between 0 and nelems at which to begin scanning
	// for the next free object in this span.
	// Each allocation scans allocBits starting at freeindex until it encounters a 0
	// indicating a free object. freeindex is then adjusted so that subsequent scans begin
	// just past the newly discovered free object.
	//
	// If freeindex == nelem, this span has no free objects.
	//
	// allocBits is a bitmap of objects in this span.
	// If n >= freeindex and allocBits[n/8] & (1<<(n%8)) is 0
	// then object n is free;
	// otherwise, object n is allocated. Bits starting at nelem are
	// undefined and should never be referenced.
	//
	// Object n starts at address n*elemsize + (start << pageShift).
	freeindex uintptr
	// TODO: Look up nelems from sizeclass and remove this field if it
	// helps performance.
	nelems uintptr // number of object in the span.

	// Cache of the allocBits at freeindex. allocCache is shifted
	// such that the lowest bit corresponds to the bit freeindex.
	// allocCache holds the complement of allocBits, thus allowing
	// ctz (count trailing zero) to use it directly.
	// allocCache may contain bits beyond s.nelems; the caller must ignore
	// these.
	allocCache uint64

	// allocBits and gcmarkBits hold pointers to a span's mark and
	// allocation bits. The pointers are 8 byte aligned.
	// There are three arenas where this data is held.
	// free: Dirty arenas that are no longer accessed
	//       and can be reused.
	// next: Holds information to be used in the next GC cycle.
	// current: Information being used during this GC cycle.
	// previous: Information being used during the last GC cycle.
	// A new GC cycle starts with the call to finishsweep_m.
	// finishsweep_m moves the previous arena to the free arena,
	// the current arena to the previous arena, and
	// the next arena to the current arena.
	// The next arena is populated as the spans request
	// memory to hold gcmarkBits for the next GC cycle as well
	// as allocBits for newly allocated spans.
	//
	// The pointer arithmetic is done "by hand" instead of using
	// arrays to avoid bounds checks along critical performance
	// paths.
	// The sweep will free the old allocBits and set allocBits to the
	// gcmarkBits. The gcmarkBits are replaced with a fresh zeroed
	// out memory.
	allocBits  *gcBits
	gcmarkBits *gcBits

	// sweep generation:
	// if sweepgen == h->sweepgen - 2, the span needs sweeping
	// if sweepgen == h->sweepgen - 1, the span is currently being swept
	// if sweepgen == h->sweepgen, the span is swept and ready to use
	// if sweepgen == h->sweepgen + 1, the span was cached before sweep began and is still cached, and needs sweeping
	// if sweepgen == h->sweepgen + 3, the span was swept and then cached and is still cached
	// h->sweepgen is incremented by 2 after every GC

	sweepgen              uint32
	divMul                uint32        // for divide by elemsize
	allocCount            uint16        // number of allocated objects
	spanclass             spanClass     // size class and noscan (uint8)
	state                 mSpanStateBox // mSpanInUse etc; accessed atomically (get/set methods)
	needzero              uint8         // needs to be zeroed before allocation
	isUserArenaChunk      bool          // whether or not this span represents a user arena
	allocCountBeforeCache uint16        // a copy of allocCount that is stored just before this span is cached
	elemsize              uintptr       // computed from sizeclass or from npages
	limit                 uintptr       // end of data in span
	speciallock           mutex         // guards specials list
	specials              *special      // linked list of special records sorted by offset.
	userArenaChunkFree    addrRange     // interval for managing chunk allocation

	// freeIndexForScan is like freeindex, except that freeindex is
	// used by the allocator whereas freeIndexForScan is used by the
	// GC scanner. They are two fields so that the GC sees the object
	// is allocated only when the object and the heap bits are
	// initialized (see also the assignment of freeIndexForScan in
	// mallocgc, and issue 54596).
	freeIndexForScan uintptr
}

type mSpanStateBox struct {
	s atomic.Uint8
}
type spanClass uint8

type mSpanList struct {
	_     sys.NotInHeap
	first *mspan // first span in list, or nil if none
	last  *mspan // last span in list, or nil if none
}

type mheap struct {
	_ sys.NotInHeap

	// lock must only be acquired on the system stack, otherwise a g
	// could self-deadlock if its stack grows with the lock held.
	lock mutex

	pages pageAlloc // page allocation data structure

	sweepgen uint32 // sweep generation, see comment in mspan; written during STW

	// allspans is a slice of all mspans ever created. Each mspan
	// appears exactly once.
	//
	// The memory for allspans is manually managed and can be
	// reallocated and move as the heap grows.
	//
	// In general, allspans is protected by mheap_.lock, which
	// prevents concurrent access as well as freeing the backing
	// store. Accesses during STW might not hold the lock, but
	// must ensure that allocation cannot happen around the
	// access (since that may free the backing store).
	allspans []*mspan // all spans out there

	// Proportional sweep
	//
	// These parameters represent a linear function from gcController.heapLive
	// to page sweep count. The proportional sweep system works to
	// stay in the black by keeping the current page sweep count
	// above this line at the current gcController.heapLive.
	//
	// The line has slope sweepPagesPerByte and passes through a
	// basis point at (sweepHeapLiveBasis, pagesSweptBasis). At
	// any given time, the system is at (gcController.heapLive,
	// pagesSwept) in this space.
	//
	// It is important that the line pass through a point we
	// control rather than simply starting at a 0,0 origin
	// because that lets us adjust sweep pacing at any time while
	// accounting for current progress. If we could only adjust
	// the slope, it would create a discontinuity in debt if any
	// progress has already been made.
	pagesInUse         atomic.Uintptr // pages of spans in stats mSpanInUse
	pagesSwept         atomic.Uint64  // pages swept this cycle
	pagesSweptBasis    atomic.Uint64  // pagesSwept to use as the origin of the sweep ratio
	sweepHeapLiveBasis uint64         // value of gcController.heapLive to use as the origin of sweep ratio; written with lock, read without
	sweepPagesPerByte  float64        // proportional sweep ratio; written with lock, read without

	// Page reclaimer state

	// reclaimIndex is the page index in allArenas of next page to
	// reclaim. Specifically, it refers to page (i %
	// pagesPerArena) of arena allArenas[i / pagesPerArena].
	//
	// If this is >= 1<<63, the page reclaimer is done scanning
	// the page marks.
	reclaimIndex atomic.Uint64

	// reclaimCredit is spare credit for extra pages swept. Since
	// the page reclaimer works in large chunks, it may reclaim
	// more than requested. Any spare pages released go to this
	// credit pool.
	reclaimCredit atomic.Uintptr

	// arenas is the heap arena map. It points to the metadata for
	// the heap for every arena frame of the entire usable virtual
	// address space.
	//
	// Use arenaIndex to compute indexes into this array.
	//
	// For regions of the address space that are not backed by the
	// Go heap, the arena map contains nil.
	//
	// Modifications are protected by mheap_.lock. Reads can be
	// performed without locking; however, a given entry can
	// transition from nil to non-nil at any time when the lock
	// isn't held. (Entries never transitions back to nil.)
	//
	// In general, this is a two-level mapping consisting of an L1
	// map and possibly many L2 maps. This saves space when there
	// are a huge number of arena frames. However, on many
	// platforms (even 64-bit), arenaL1Bits is 0, making this
	// effectively a single-level map. In this case, arenas[0]
	// will never be nil.
	arenas [1 << arenaL1Bits]*[1 << arenaL2Bits]*heapArena

	// heapArenaAlloc is pre-reserved space for allocating heapArena
	// objects. This is only used on 32-bit, where we pre-reserve
	// this space to avoid interleaving it with the heap itself.
	heapArenaAlloc linearAlloc

	// arenaHints is a list of addresses at which to attempt to
	// add more heap arenas. This is initially populated with a
	// set of general hint addresses, and grown with the bounds of
	// actual heap arena ranges.
	arenaHints *arenaHint

	// arena is a pre-reserved space for allocating heap arenas
	// (the actual arenas). This is only used on 32-bit.
	arena linearAlloc

	// allArenas is the arenaIndex of every mapped arena. This can
	// be used to iterate through the address space.
	//
	// Access is protected by mheap_.lock. However, since this is
	// append-only and old backing arrays are never freed, it is
	// safe to acquire mheap_.lock, copy the slice header, and
	// then release mheap_.lock.
	allArenas []arenaIdx

	// sweepArenas is a snapshot of allArenas taken at the
	// beginning of the sweep cycle. This can be read safely by
	// simply blocking GC (by disabling preemption).
	sweepArenas []arenaIdx

	// markArenas is a snapshot of allArenas taken at the beginning
	// of the mark cycle. Because allArenas is append-only, neither
	// this slice nor its contents will change during the mark, so
	// it can be read safely.
	markArenas []arenaIdx

	// curArena is the arena that the heap is currently growing
	// into. This should always be physPageSize-aligned.
	curArena struct {
		base, end uintptr
	}

	// central free lists for small size classes.
	// the padding makes sure that the mcentrals are
	// spaced CacheLinePadSize bytes apart, so that each mcentral.lock
	// gets its own cache line.
	// central is indexed by spanClass.
	central [numSpanClasses]struct {
		mcentral mcentral
		pad      [(cpu.CacheLinePadSize - unsafe.Sizeof(mcentral{})%cpu.CacheLinePadSize) % cpu.CacheLinePadSize]byte
	}

	spanalloc             fixalloc // allocator for span*
	cachealloc            fixalloc // allocator for mcache*
	specialfinalizeralloc fixalloc // allocator for specialfinalizer*
	specialprofilealloc   fixalloc // allocator for specialprofile*
	specialReachableAlloc fixalloc // allocator for specialReachable
	speciallock           mutex    // lock for special record allocators.
	arenaHintAlloc        fixalloc // allocator for arenaHints

	// User arena state.
	//
	// Protected by mheap_.lock.
	userArena struct {
		// arenaHints is a list of addresses at which to attempt to
		// add more heap arenas for user arena chunks. This is initially
		// populated with a set of general hint addresses, and grown with
		// the bounds of actual heap arena ranges.
		arenaHints *arenaHint

		// quarantineList is a list of user arena spans that have been set to fault, but
		// are waiting for all pointers into them to go away. Sweeping handles
		// identifying when this is true, and moves the span to the ready list.
		quarantineList mSpanList

		// readyList is a list of empty user arena spans that are ready for reuse.
		readyList mSpanList
	}

	unused *specialfinalizer // never set, just here to force the specialfinalizer type into DWARF
}
type arenaHint struct {
	_    sys.NotInHeap
	addr uintptr
	down bool
	next *arenaHint
}

type heapArena struct {
	_ sys.NotInHeap

	// bitmap stores the pointer/scalar bitmap for the words in
	// this arena. See mbitmap.go for a description.
	// This array uses 1 bit per word of heap, or 1.6% of the heap size (for 64-bit).
	bitmap [heapArenaBitmapWords]uintptr

	// If the ith bit of noMorePtrs is true, then there are no more
	// pointers for the object containing the word described by the
	// high bit of bitmap[i].
	// In that case, bitmap[i+1], ... must be zero until the start
	// of the next object.
	// We never operate on these entries using bit-parallel techniques,
	// so it is ok if they are small. Also, they can't be bigger than
	// uint16 because at that size a single noMorePtrs entry
	// represents 8K of memory, the minimum size of a span. Any larger
	// and we'd have to worry about concurrent updates.
	// This array uses 1 bit per word of bitmap, or .024% of the heap size (for 64-bit).
	noMorePtrs [heapArenaBitmapWords / 8]uint8

	// spans maps from virtual address page ID within this arena to *mspan.
	// For allocated spans, their pages map to the span itself.
	// For free spans, only the lowest and highest pages map to the span itself.
	// Internal pages map to an arbitrary span.
	// For pages that have never been allocated, spans entries are nil.
	//
	// Modifications are protected by mheap.lock. Reads can be
	// performed without locking, but ONLY from indexes that are
	// known to contain in-use or stack spans. This means there
	// must not be a safe-point between establishing that an
	// address is live and looking it up in the spans array.
	spans [pagesPerArena]*mspan

	// pageInUse is a bitmap that indicates which spans are in
	// state mSpanInUse. This bitmap is indexed by page number,
	// but only the bit corresponding to the first page in each
	// span is used.
	//
	// Reads and writes are atomic.
	pageInUse [pagesPerArena / 8]uint8

	// pageMarks is a bitmap that indicates which spans have any
	// marked objects on them. Like pageInUse, only the bit
	// corresponding to the first page in each span is used.
	//
	// Writes are done atomically during marking. Reads are
	// non-atomic and lock-free since they only occur during
	// sweeping (and hence never race with writes).
	//
	// This is used to quickly find whole spans that can be freed.
	//
	// TODO(austin): It would be nice if this was uint64 for
	// faster scanning, but we don't have 64-bit atomic bit
	// operations.
	pageMarks [16 / 8]uint8

	// pageSpecials is a bitmap that indicates which spans have
	// specials (finalizers or other). Like pageInUse, only the bit
	// corresponding to the first page in each span is used.
	//
	// Writes are done atomically whenever a special is added to
	// a span and whenever the last special is removed from a span.
	// Reads are done atomically to find spans containing specials
	// during marking.
	pageSpecials [16 / 8]uint8

	// checkmarks stores the debug.gccheckmark state. It is only
	// used if debug.gccheckmark > 0.
	checkmarks *checkmarksMap

	// zeroedBase marks the first byte of the first page in this
	// arena which hasn't been used yet and is therefore already
	// zero. zeroedBase is relative to the arena base.
	// Increases monotonically until it hits heapArenaBytes.
	//
	// This field is sufficient to determine if an allocation
	// needs to be zeroed because the page allocator follows an
	// address-ordered first-fit policy.
	//
	// Read atomically and written with an atomic CAS.
	zeroedBase uintptr
}

type specialfinalizer struct {
}

type gcBitsArena struct{}

var gcBitsArenas struct {
	lock     mutex
	free     *gcBitsArena
	next     *gcBitsArena // Read atomically. Write atomically under lock.
	current  *gcBitsArena
	previous *gcBitsArena
}

type arenaIdx uint

//go:nosplit
func (i arenaIdx) l1() uint { return 0 }

//go:nosplit
func (i arenaIdx) l2() uint { return 0 }

//go:nosplit
func arenaIndex(p uintptr) arenaIdx { return 0 }

const (
	numSpanClasses = 16 << 1
	tinySpanClass  = spanClass(1<<1 | 1)
)

type specialReachable struct {
	special   special
	done      bool
	reachable bool
}

func removefinalizer(p unsafe.Pointer) {}

func makeSpanClass(sizeclass uint8, noscan bool) spanClass { return 0 }
func addspecial(p unsafe.Pointer, s *special) bool         { return false }
func setprofilebucket(p unsafe.Pointer, b *bucket)         {}
func (s *mspan) base() uintptr                             { return 0 }

const (
	_KindSpecialFinalizer = 1
	_KindSpecialProfile   = 2
	// _KindSpecialReachable is a special used for tracking
	// reachability during testing.
	_KindSpecialReachable = 3
	// Note: The finalizer special must be first because if we're freeing
	// an object, a finalizer special will cause the freeing operation
	// to abort, and we want to keep the other special records around
	// if that happens.
)

func addfinalizer(p unsafe.Pointer, f *funcval, nret uintptr, fint *_type, ot *ptrtype) bool {
	return false
}
func (h *mheap) init() {}
