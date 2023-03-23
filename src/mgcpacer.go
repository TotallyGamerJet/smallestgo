package runtime

import (
	"internal/cpu"
	"runtime/internal/atomic"
	_ "unsafe" // for go:linkname
)

var gcController gcControllerState

type gcControllerState struct {
	// Initialized from GOGC. GOGC=off means no GC.
	gcPercent atomic.Int32

	// memoryLimit is the soft memory limit in bytes.
	//
	// Initialized from GOMEMLIMIT. GOMEMLIMIT=off is equivalent to MaxInt64
	// which means no soft memory limit in practice.
	//
	// This is an int64 instead of a uint64 to more easily maintain parity with
	// the SetMemoryLimit API, which sets a maximum at MaxInt64. This value
	// should never be negative.
	memoryLimit atomic.Int64

	// heapMinimum is the minimum heap size at which to trigger GC.
	// For small heaps, this overrides the usual GOGC*live set rule.
	//
	// When there is a very small live set but a lot of allocation, simply
	// collecting when the heap reaches GOGC*live results in many GC
	// cycles and high total per-GC overhead. This minimum amortizes this
	// per-GC overhead while keeping the heap reasonably small.
	//
	// During initialization this is set to 4MB*GOGC/100. In the case of
	// GOGC==0, this will set heapMinimum to 0, resulting in constant
	// collection even when the heap size is small, which is useful for
	// debugging.
	heapMinimum uint64

	// runway is the amount of runway in heap bytes allocated by the
	// application that we want to give the GC once it starts.
	//
	// This is computed from consMark during mark termination.
	runway atomic.Uint64

	// consMark is the estimated per-CPU consMark ratio for the application.
	//
	// It represents the ratio between the application's allocation
	// rate, as bytes allocated per CPU-time, and the GC's scan rate,
	// as bytes scanned per CPU-time.
	// The units of this ratio are (B / cpu-ns) / (B / cpu-ns).
	//
	// At a high level, this value is computed as the bytes of memory
	// allocated (cons) per unit of scan work completed (mark) in a GC
	// cycle, divided by the CPU time spent on each activity.
	//
	// Updated at the end of each GC cycle, in endCycle.
	consMark float64

	// lastConsMark is the computed cons/mark value for the previous GC
	// cycle. Note that this is *not* the last value of cons/mark, but the
	// actual computed value. See endCycle for details.
	lastConsMark float64

	// gcPercentHeapGoal is the goal heapLive for when next GC ends derived
	// from gcPercent.
	//
	// Set to ^uint64(0) if gcPercent is disabled.
	gcPercentHeapGoal atomic.Uint64

	// sweepDistMinTrigger is the minimum trigger to ensure a minimum
	// sweep distance.
	//
	// This bound is also special because it applies to both the trigger
	// *and* the goal (all other trigger bounds must be based *on* the goal).
	//
	// It is computed ahead of time, at commit time. The theory is that,
	// absent a sudden change to a parameter like gcPercent, the trigger
	// will be chosen to always give the sweeper enough headroom. However,
	// such a change might dramatically and suddenly move up the trigger,
	// in which case we need to ensure the sweeper still has enough headroom.
	sweepDistMinTrigger atomic.Uint64

	// triggered is the point at which the current GC cycle actually triggered.
	// Only valid during the mark phase of a GC cycle, otherwise set to ^uint64(0).
	//
	// Updated while the world is stopped.
	triggered uint64

	// lastHeapGoal is the value of heapGoal at the moment the last GC
	// ended. Note that this is distinct from the last value heapGoal had,
	// because it could change if e.g. gcPercent changes.
	//
	// Read and written with the world stopped or with mheap_.lock held.
	lastHeapGoal uint64

	// heapLive is the number of bytes considered live by the GC.
	// That is: retained by the most recent GC plus allocated
	// since then. heapLive ≤ memstats.totalAlloc-memstats.totalFree, since
	// heapAlloc includes unmarked objects that have not yet been swept (and
	// hence goes up as we allocate and down as we sweep) while heapLive
	// excludes these objects (and hence only goes up between GCs).
	//
	// To reduce contention, this is updated only when obtaining a span
	// from an mcentral and at this point it counts all of the unallocated
	// slots in that span (which will be allocated before that mcache
	// obtains another span from that mcentral). Hence, it slightly
	// overestimates the "true" live heap size. It's better to overestimate
	// than to underestimate because 1) this triggers the GC earlier than
	// necessary rather than potentially too late and 2) this leads to a
	// conservative GC rate rather than a GC rate that is potentially too
	// low.
	//
	// Whenever this is updated, call traceHeapAlloc() and
	// this gcControllerState's revise() method.
	heapLive atomic.Uint64

	// heapScan is the number of bytes of "scannable" heap. This is the
	// live heap (as counted by heapLive), but omitting no-scan objects and
	// no-scan tails of objects.
	//
	// This value is fixed at the start of a GC cycle. It represents the
	// maximum scannable heap.
	heapScan atomic.Uint64

	// lastHeapScan is the number of bytes of heap that were scanned
	// last GC cycle. It is the same as heapMarked, but only
	// includes the "scannable" parts of objects.
	//
	// Updated when the world is stopped.
	lastHeapScan uint64

	// lastStackScan is the number of bytes of stack that were scanned
	// last GC cycle.
	lastStackScan atomic.Uint64

	// maxStackScan is the amount of allocated goroutine stack space in
	// use by goroutines.
	//
	// This number tracks allocated goroutine stack space rather than used
	// goroutine stack space (i.e. what is actually scanned) because used
	// goroutine stack space is much harder to measure cheaply. By using
	// allocated space, we make an overestimate; this is OK, it's better
	// to conservatively overcount than undercount.
	maxStackScan atomic.Uint64

	// globalsScan is the total amount of global variable space
	// that is scannable.
	globalsScan atomic.Uint64

	// heapMarked is the number of bytes marked by the previous
	// GC. After mark termination, heapLive == heapMarked, but
	// unlike heapLive, heapMarked does not change until the
	// next mark termination.
	heapMarked uint64

	// heapScanWork is the total heap scan work performed this cycle.
	// stackScanWork is the total stack scan work performed this cycle.
	// globalsScanWork is the total globals scan work performed this cycle.
	//
	// These are updated atomically during the cycle. Updates occur in
	// bounded batches, since they are both written and read
	// throughout the cycle. At the end of the cycle, heapScanWork is how
	// much of the retained heap is scannable.
	//
	// Currently these are measured in bytes. For most uses, this is an
	// opaque unit of work, but for estimation the definition is important.
	//
	// Note that stackScanWork includes only stack space scanned, not all
	// of the allocated stack.
	heapScanWork    atomic.Int64
	stackScanWork   atomic.Int64
	globalsScanWork atomic.Int64

	// bgScanCredit is the scan work credit accumulated by the concurrent
	// background scan. This credit is accumulated by the background scan
	// and stolen by mutator assists.  Updates occur in bounded batches,
	// since it is both written and read throughout the cycle.
	bgScanCredit atomic.Int64

	// assistTime is the nanoseconds spent in mutator assists
	// during this cycle. This is updated atomically, and must also
	// be updated atomically even during a STW, because it is read
	// by sysmon. Updates occur in bounded batches, since it is both
	// written and read throughout the cycle.
	assistTime atomic.Int64

	// dedicatedMarkTime is the nanoseconds spent in dedicated mark workers
	// during this cycle. This is updated at the end of the concurrent mark
	// phase.
	dedicatedMarkTime atomic.Int64

	// fractionalMarkTime is the nanoseconds spent in the fractional mark
	// worker during this cycle. This is updated throughout the cycle and
	// will be up-to-date if the fractional mark worker is not currently
	// running.
	fractionalMarkTime atomic.Int64

	// idleMarkTime is the nanoseconds spent in idle marking during this
	// cycle. This is updated throughout the cycle.
	idleMarkTime atomic.Int64

	// markStartTime is the absolute start time in nanoseconds
	// that assists and background mark workers started.
	markStartTime int64

	// dedicatedMarkWorkersNeeded is the number of dedicated mark workers
	// that need to be started. This is computed at the beginning of each
	// cycle and decremented as dedicated mark workers get started.
	dedicatedMarkWorkersNeeded atomic.Int64

	// idleMarkWorkers is two packed int32 values in a single uint64.
	// These two values are always updated simultaneously.
	//
	// The bottom int32 is the current number of idle mark workers executing.
	//
	// The top int32 is the maximum number of idle mark workers allowed to
	// execute concurrently. Normally, this number is just gomaxprocs. However,
	// during periodic GC cycles it is set to 0 because the system is idle
	// anyway; there's no need to go full blast on all of GOMAXPROCS.
	//
	// The maximum number of idle mark workers is used to prevent new workers
	// from starting, but it is not a hard maximum. It is possible (but
	// exceedingly rare) for the current number of idle mark workers to
	// transiently exceed the maximum. This could happen if the maximum changes
	// just after a GC ends, and an M with no P.
	//
	// Note that if we have no dedicated mark workers, we set this value to
	// 1 in this case we only have fractional GC workers which aren't scheduled
	// strictly enough to ensure GC progress. As a result, idle-priority mark
	// workers are vital to GC progress in these situations.
	//
	// For example, consider a situation in which goroutines block on the GC
	// (such as via runtime.GOMAXPROCS) and only fractional mark workers are
	// scheduled (e.g. GOMAXPROCS=1). Without idle-priority mark workers, the
	// last running M might skip scheduling a fractional mark worker if its
	// utilization goal is met, such that once it goes to sleep (because there's
	// nothing to do), there will be nothing else to spin up a new M for the
	// fractional worker in the future, stalling GC progress and causing a
	// deadlock. However, idle-priority workers will *always* run when there is
	// nothing left to do, ensuring the GC makes progress.
	//
	// See github.com/golang/go/issues/44163 for more details.
	idleMarkWorkers atomic.Uint64

	// assistWorkPerByte is the ratio of scan work to allocated
	// bytes that should be performed by mutator assists. This is
	// computed at the beginning of each cycle and updated every
	// time heapScan is updated.
	assistWorkPerByte atomic.Float64

	// assistBytesPerWork is 1/assistWorkPerByte.
	//
	// Note that because this is read and written independently
	// from assistWorkPerByte users may notice a skew between
	// the two values, and such a state should be safe.
	assistBytesPerWork atomic.Float64

	// fractionalUtilizationGoal is the fraction of wall clock
	// time that should be spent in the fractional mark worker on
	// each P that isn't running a dedicated worker.
	//
	// For example, if the utilization goal is 25% and there are
	// no dedicated workers, this will be 0.25. If the goal is
	// 25%, there is one dedicated worker, and GOMAXPROCS is 5,
	// this will be 0.05 to make up the missing 5%.
	//
	// If this is zero, no fractional workers are needed.
	fractionalUtilizationGoal float64

	// These memory stats are effectively duplicates of fields from
	// memstats.heapStats but are updated atomically or with the world
	// stopped and don't provide the same consistency guarantees.
	//
	// Because the runtime is responsible for managing a memory limit, it's
	// useful to couple these stats more tightly to the gcController, which
	// is intimately connected to how that memory limit is maintained.
	heapInUse    sysMemStat    // bytes in mSpanInUse spans
	heapReleased sysMemStat    // bytes released to the OS
	heapFree     sysMemStat    // bytes not in any span, but not released to the OS
	totalAlloc   atomic.Uint64 // total bytes allocated
	totalFree    atomic.Uint64 // total bytes freed
	mappedReady  atomic.Uint64 // total virtual memory in the Ready state (see mem.go).

	// test indicates that this is a test-only copy of gcControllerState.
	test bool

	_ cpu.CacheLinePad
}

func (c *gcControllerState) addGlobals(amount int64)                                      {}
func (c *gcControllerState) init(gcPercent int32, memoryLimit int64)                      {}
func readGOGC() int32                                                                     { return 0 }
func readGOMEMLIMIT() int64                                                               { return 0 }
func (c *gcControllerState) trigger() (uint64, uint64)                                    { return 0, 0 }
func (c *gcControllerState) startCycle(markStartTime int64, procs int, trigger gcTrigger) {}
func (c *gcControllerState) endCycle(now int64, procs int, userForced bool)               {}

//go:systemstack
func gcControllerCommit()                                                         {}
func (c *gcControllerState) markWorkerStop(mode gcMarkWorkerMode, duration int64) {}
func (c *gcControllerState) resetLive(bytesMarked uint64)                         {}
func (c *gcControllerState) findRunnableGCWorker(pp *p, now int64) (*g, int64)    { return nil, 0 }

//go:nosplit
func (c *gcControllerState) addIdleMarkWorker() bool { return false }

func (c *gcControllerState) removeIdleMarkWorker() {}

//go:nosplit
func (c *gcControllerState) needIdleMarkWorker() bool              { return false }
func (c *gcControllerState) addScannableStack(pp *p, amount int64) {}
func (c *gcControllerState) heapGoal() uint64                      { return 0 }
