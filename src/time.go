package runtime

import _ "unsafe"
import "sync/atomic"

type timer struct {
	// If this timer is on a heap, which P's heap it is on.
	// puintptr rather than *p to match uintptr in the versions
	// of this struct defined in other packages.
	pp puintptr

	// Timer wakes up at when, and then at when+period, ... (period > 0 only)
	// each time calling f(arg, now) in the timer goroutine, so f must be
	// a well-behaved function and not block.
	//
	// when must be positive on an active timer.
	when   int64
	period int64
	f      func(any, uintptr)
	arg    any
	seq    uintptr

	// What to set the when field to in timerModifiedXX status.
	nextwhen int64

	// The status field holds one of the values below.
	status atomic.Uint32
}

//go:linkname resetTimer time.resetTimer
func resetTimer(t *timer, when int64) bool {
	return false
}

//go:linkname stopTimer time.stopTimer
func stopTimer(t *timer) bool {
	return false
}

func modtimer(t *timer, when, period int64, f func(any, uintptr), arg any, seq uintptr) bool {
	return false
}

//go:nowritebarrierrec
func nobarrierWakeTime(pp *p) int64 {
	return 0
}

func adjusttimers(pp *p, now int64) {}

//go:systemstack
func runtimer(pp *p, now int64) int64   { return 0 }
func deltimer(t *timer) bool            { return false }
func clearDeletedTimers(pp *p)          {}
func moveTimers(pp *p, timers []*timer) {}
func timeSleepUntil() int64             { return 0 }

func resettimer(t *timer, when int64) bool { return false }

const maxWhen = 1<<63 - 1
