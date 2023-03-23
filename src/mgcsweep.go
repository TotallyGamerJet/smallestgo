package runtime

import "sync/atomic"

func deductSweepCredit(spanBytes uintptr, callerSweepPages uintptr) {}

type sweepLocker struct {
	// sweepGen is the sweep generation of the heap.
	sweepGen uint32
	valid    bool
}
type sweepLocked struct {
	*mspan
}

type sweepdata struct {
	lock   mutex
	g      *g
	parked bool

	nbgsweep    uint32
	npausesweep uint32

	// active tracks outstanding sweepers and the sweep
	// termination condition.
	active activeSweep

	// centralIndex is the current unswept span class.
	// It represents an index into the mcentral span
	// sets. Accessed and updated via its load and
	// update methods. Not protected by a lock.
	//
	// Reset at mark termination.
	// Used by mheap.nextSpanForSweep.
	centralIndex sweepClass
}

type activeSweep struct {
	state atomic.Uint32
}

const sweepDrainedMask = 1 << 31

type sweepClass uint32

func (a *activeSweep) begin() sweepLocker { return sweepLocker{} }
func (a *activeSweep) end(sl sweepLocker) {}

var sweep sweepdata

func (sl *sweepLocked) sweep(preserve bool) bool               { return false }
func (l *sweepLocker) tryAcquire(s *mspan) (sweepLocked, bool) { return sweepLocked{}, false }
func bgsweep(c chan int)                                       {}
func sweepone() uintptr                                        { return 0 }
func isSweepDone() bool                                        { return false }

//go:nowritebarrier
func finishsweep_m()               {}
func (a *activeSweep) reset()      {}
func (s *sweepClass) clear()       {}
func gcPaceSweeper(trigger uint64) {}

//go:nowritebarrier
func (s *mspan) ensureSwept() {}
