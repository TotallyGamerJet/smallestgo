package runtime

const (
	fixedRootFinalizers = iota
	fixedRootFreeGStacks
	fixedRootCount

	// rootBlockBytes is the number of bytes to scan per data or
	// BSS root.
	rootBlockBytes = 256 << 10

	// maxObletBytes is the maximum bytes of an object to scan at
	// once. Larger objects will be split up into "oblets" of at
	// most this size. Since we can scan 1–2 MB/ms, 128 KB bounds
	// scan preemption at ~100 µs.
	//
	// This must be > _MaxSmallSize so that the object base is the
	// span base.
	maxObletBytes = 128 << 10

	// drainCheckThreshold specifies how many units of work to do
	// between self-preemption checks in gcDrain. Assuming a scan
	// rate of 1 MB/ms, this is ~100 µs. Lower values have higher
	// overhead in the scan loop (the scheduler check may perform
	// a syscall, so its overhead is nontrivial). Higher values
	// make the system less responsive to incoming work.
	drainCheckThreshold = 100000

	// pagesPerSpanRoot indicates how many pages to scan from a span root
	// at a time. Used by special root marking.
	//
	// Higher values improve throughput by increasing locality, but
	// increase the minimum latency of a marking operation.
	//
	// Must be a multiple of the pageInUse bitmap element size and
	// must also evenly divide pagesPerArena.
	pagesPerSpanRoot = 512
)

//go:nowritebarrier
//go:nosplit
func gcmarknewobject(span *mspan, obj, size uintptr) {}

func gcAssistAlloc(gp *g)                         {}
func gcDumpObject(label string, obj, off uintptr) {}
func gcMarkRootPrepare()                          {}
func gcMarkTinyAllocs()                           {}
func gcWakeAllAssists()                           {}

//go:nowritebarrier
func gcDrain(gcw *gcWork, flags gcDrainFlags) {}

type gcDrainFlags int

const (
	gcDrainUntilPreempt gcDrainFlags = 1 << iota
	gcDrainFlushBgCredit
	gcDrainIdle
	gcDrainFractional
)

func gcMarkRootCheck() {}

//go:nowritebarrier
func scanobject(b uintptr, gcw *gcWork) {}

var oneptrmask = [...]uint8{1}

//go:nowritebarrier
func scanblock(b0, n0 uintptr, ptrmask *uint8, gcw *gcWork, stk *stackScanState) {}
