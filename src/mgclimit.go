package runtime

import "runtime/internal/atomic"

var gcCPULimiter gcCPULimiterState

type gcCPULimiterState struct {
	lock atomic.Uint32

	enabled atomic.Bool
	bucket  struct {
		// Invariants:
		// - fill >= 0
		// - capacity >= 0
		// - fill <= capacity
		fill, capacity uint64
	}
	// overflow is the cumulative amount of GC CPU time that we tried to fill the
	// bucket with but exceeded its capacity.
	overflow uint64

	// gcEnabled is an internal copy of gcBlackenEnabled that determines
	// whether the limiter tracks total assist time.
	//
	// gcBlackenEnabled isn't used directly so as to keep this structure
	// unit-testable.
	gcEnabled bool

	// transitioning is true when the GC is in a STW and transitioning between
	// the mark and sweep phases.
	transitioning bool

	// assistTimePool is the accumulated assist time since the last update.
	assistTimePool atomic.Int64

	// idleMarkTimePool is the accumulated idle mark time since the last update.
	idleMarkTimePool atomic.Int64

	// idleTimePool is the accumulated time Ps spent on the idle list since the last update.
	idleTimePool atomic.Int64

	// lastUpdate is the nanotime timestamp of the last time update was called.
	//
	// Updated under lock, but may be read concurrently.
	lastUpdate atomic.Int64

	// lastEnabledCycle is the GC cycle that last had the limiter enabled.
	lastEnabledCycle atomic.Uint32

	// nprocs is an internal copy of gomaxprocs, used to determine total available
	// CPU time.
	//
	// gomaxprocs isn't used directly so as to keep this structure unit-testable.
	nprocs int32

	// test indicates whether this instance of the struct was made for testing purposes.
	test bool
}

type limiterEventType uint8

const (
	limiterEventNone           limiterEventType = iota // None of the following events.
	limiterEventIdleMarkWork                           // Refers to an idle mark worker (see gcMarkWorkerMode).
	limiterEventMarkAssist                             // Refers to mark assist (see gcAssistAlloc).
	limiterEventScavengeAssist                         // Refers to a scavenge assist (see allocSpan).
	limiterEventIdle                                   // Refers to time a P spent on the idle list.

	limiterEventBits = 3
)

type limiterEvent struct{}

func (l *gcCPULimiterState) startGCTransition(enableGC bool, now int64) {}
func (l *gcCPULimiterState) finishGCTransition(now int64)               {}
func (e *limiterEvent) start(typ limiterEventType, now int64) bool      { return false }
func (e *limiterEvent) stop(typ limiterEventType, now int64)            {}
func (l *gcCPULimiterState) needUpdate(now int64) bool                  { return false }
func (l *gcCPULimiterState) update(now int64)                           {}
func (l *gcCPULimiterState) resetCapacity(now int64, nprocs int32)      {}
