package runtime

type suspendGState struct {
	g *g

	// dead indicates the goroutine was not suspended because it
	// is dead. This goroutine could be reused after the dead
	// state was observed, so the caller must not assume that it
	// remains dead.
	dead bool

	// stopped indicates that this suspendG transitioned the G to
	// _Gwaiting via g.preemptStop and thus is responsible for
	// readying it when done.
	stopped bool
}

func wantAsyncPreempt(gp *g) bool                                { return false }
func isAsyncSafePoint(gp *g, pc, sp, lr uintptr) (bool, uintptr) { return false, 0 }

func asyncPreempt()

//go:systemstack
func suspendG(gp *g) suspendGState { return suspendGState{} }
func resumeG(state suspendGState)  {}

//go:nosplit
func canPreemptM(mp *m) bool { return false }
