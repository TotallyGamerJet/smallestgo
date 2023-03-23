package runtime

import "sync/atomic"

type dlogger struct {

	// allLink is the next dlogger in the allDloggers list.
	allLink *dlogger

	// owned indicates that this dlogger is owned by an M. This is
	// accessed atomically.
	owned atomic.Uint32
}
