package runtime

var writeBarrier struct {
	enabled bool    // compiler emits a check of this before calling write barrier
	pad     [3]byte // compiler uses 32-bit load for "enabled" field
	needed  bool    // whether we need a write barrier for current GC phase
	cgo     bool    // whether we need a write barrier for a cgo check
	alignme uint64  // guarantee alignment so that compiler can use a 32 or 64-bit load
}

type gcMarkWorkerMode int

const (
	// gcMarkWorkerNotWorker indicates that the next scheduled G is not
	// starting work and the mode should be ignored.
	gcMarkWorkerNotWorker gcMarkWorkerMode = iota

	// gcMarkWorkerDedicatedMode indicates that the P of a mark
	// worker is dedicated to running that mark worker. The mark
	// worker should run without preemption.
	gcMarkWorkerDedicatedMode

	// gcMarkWorkerFractionalMode indicates that a P is currently
	// running the "fractional" mark worker. The fractional worker
	// is necessary when GOMAXPROCS*gcBackgroundUtilization is not
	// an integer and using only dedicated workers would result in
	// utilization too far from the target of gcBackgroundUtilization.
	// The fractional worker should run until it is preempted and
	// will be scheduled to pick up the fractional part of
	// GOMAXPROCS*gcBackgroundUtilization.
	gcMarkWorkerFractionalMode

	// gcMarkWorkerIdleMode indicates that a P is running the mark
	// worker because it has nothing else to do. The idle worker
	// should run until it is preempted and account its time
	// against gcController.idleMarkTime.
	gcMarkWorkerIdleMode
)

const (
	_DebugGC         = 0
	_ConcurrentSweep = true
	_FinBlockSize    = 4 * 1024

	// debugScanConservative enables debug logging for stack
	// frames that are scanned conservatively.
	debugScanConservative = false

	// sweepMinHeapDistance is a lower bound on the heap distance
	// (in bytes) reserved for concurrent sweeping between GC
	// cycles.
	sweepMinHeapDistance = 1024 * 1024
)

type gcTrigger struct {
	kind gcTriggerKind
	now  int64  // gcTriggerTime: current time
	n    uint32 // gcTriggerCycle: cycle number to start
}

type gcTriggerKind int

const (
	// gcTriggerHeap indicates that a cycle should be started when
	// the heap size reaches the trigger heap size computed by the
	// controller.
	gcTriggerHeap gcTriggerKind = iota

	// gcTriggerTime indicates that a cycle should be started when
	// it's been more than forcegcperiod nanoseconds since the
	// previous GC cycle.
	gcTriggerTime

	// gcTriggerCycle indicates that a cycle should be started if
	// we have not yet started cycle number gcTrigger.n (relative
	// to work.cycles).
	gcTriggerCycle
)

var gcMarkWorkerModeStrings = [...]string{
	"Not worker",
	"GC (dedicated)",
	"GC (fractional)",
	"GC (idle)",
}

func gcenable()                 {}
func gcStart(trigger gcTrigger) {}
func gcinit()                   {}

var gcBlackenEnabled uint32

func gcMarkWorkAvailable(p *p) bool { return false }

type gcBgMarkWorkerNode struct {
	// Unused workers are managed in a lock-free stack. This field must be first.
	node lfnode

	// The g of this worker.
	gp guintptr

	// Release this m on park. This is used to communicate with the unlock
	// function, which cannot access the G's stack. It is unused outside of
	// gcBgMarkWorker().
	m muintptr
}

var gcphase uint32

const (
	_GCoff             = iota // GC not running; sweeping in background, write barrier disabled
	_GCmark                   // GC marking roots and workbufs: allocate black, write barrier ENABLED
	_GCmarktermination        // GC mark termination: allocate black, P's help GC, write barrier ENABLED
)

func fmtNSAsMS(buf []byte, ns uint64) []byte { return nil }

func (t gcTrigger) test() bool { return false }
