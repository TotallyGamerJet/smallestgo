package runtime

type lockRankStruct struct {
}

func lockInit(l *mutex, rank lockRank) {
}

func getLockRank(l *mutex) lockRank {
	return 0
}

func lockWithRank(l *mutex, rank lockRank) {

}

// This function may be called in nosplit context and thus must be nosplit.
//
//go:nosplit
func acquireLockRank(rank lockRank) {
}

func unlockWithRank(l *mutex) {

}

// This function may be called in nosplit context and thus must be nosplit.
//
//go:nosplit
func releaseLockRank(rank lockRank) {
}

func lockWithRankMayAcquire(l *mutex, rank lockRank) {
}

//go:nosplit
func assertLockHeld(l *mutex) {
}

//go:nosplit
func assertRankHeld(r lockRank) {
}

//go:nosplit
func worldStopped() {
}

//go:nosplit
func worldStarted() {
}

//go:nosplit
func assertWorldStopped() {
}

//go:nosplit
func assertWorldStoppedOrLockHeld(l *mutex) {
}
