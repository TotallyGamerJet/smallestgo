package runtime

import "runtime/internal/sys"
import "runtime/internal/atomic"

type pollDesc struct {
	_    sys.NotInHeap
	link *pollDesc // in pollcache, protected by pollcache.lock
	fd   uintptr   // constant for pollDesc usage lifetime

	// atomicInfo holds bits from closing, rd, and wd,
	// which are only ever written while holding the lock,
	// summarized for use by netpollcheckerr,
	// which cannot acquire the lock.
	// After writing these fields under lock in a way that
	// might change the summary, code must call publishInfo
	// before releasing the lock.
	// Code that changes fields and then calls netpollunblock
	// (while still holding the lock) must call publishInfo
	// before calling netpollunblock, because publishInfo is what
	// stops netpollblock from blocking anew
	// (by changing the result of netpollcheckerr).
	// atomicInfo also holds the eventErr bit,
	// recording whether a poll event on the fd got an error;
	// atomicInfo is the only source of truth for that bit.
	atomicInfo atomic.Uint32 // atomic pollInfo

	// rg, wg are accessed atomically and hold g pointers.
	// (Using atomic.Uintptr here is similar to using guintptr elsewhere.)
	rg atomic.Uintptr // pdReady, pdWait, G waiting for read or pdNil
	wg atomic.Uintptr // pdReady, pdWait, G waiting for write or pdNil

	lock    mutex // protects the following fields
	closing bool
	user    uint32    // user settable cookie
	rseq    uintptr   // protects from stale read timers
	rt      timer     // read deadline timer (set if rt.f != nil)
	rd      int64     // read deadline (a nanotime in the future, -1 when expired)
	wseq    uintptr   // protects from stale write timers
	wt      timer     // write deadline timer
	wd      int64     // write deadline (a nanotime in the future, -1 when expired)
	self    *pollDesc // storage for indirect interface. See (*pollDesc).makeArg.
}

var netpollWaiters atomic.Uint32

func netpollinited() bool { return false }

//go:nowritebarrier
func netpollready(toRun *gList, pd *pollDesc, mode int32) {}
