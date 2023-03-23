package runtime

import "unsafe"

func inUserArenaChunk(p uintptr) bool { return false }

type liveUserArenaChunk struct {
	*mspan
	// Reference to mspan.base() to keep the chunk alive.
	x unsafe.Pointer
}

var userArenaState struct {
	lock mutex

	// reuse contains a list of partially-used and already-live
	// user arena chunks that can be quickly reused for another
	// arena.
	//
	// Protected by lock.
	reuse []liveUserArenaChunk

	// fault contains full user arena chunks that need to be faulted.
	//
	// Protected by lock.
	fault []liveUserArenaChunk
}

func (s *mspan) setUserArenaChunkToFault() {}
