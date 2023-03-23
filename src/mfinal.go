package runtime

import (
	"runtime/internal/atomic"
)

var fing *g

// finalizer goroutine status.
const (
	fingUninitialized uint32 = iota
	fingCreated       uint32 = 1 << (iota - 1)
	fingRunningFinalizer
	fingWait
	fingWake
)

var finlock mutex
var fingStatus atomic.Uint32

func KeepAlive(x any) {
	// Introduce a use of x that the compiler can't eliminate.
	// This makes sure x is alive on entry. We need x to be alive
	// on entry for "defer runtime.KeepAlive(x)"; see issue 21402.
	if cgoAlwaysFalse {
		println(x)
	}
}
func wakefing() *g { return nil }
