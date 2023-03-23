package runtime

import "sync/atomic"

type throwType uint32

const (
	// throwTypeNone means that we are not throwing.
	throwTypeNone throwType = iota

	// throwTypeUser is a throw due to a problem with the application.
	//
	// These throws do not include runtime frames, system goroutines, or
	// frame metadata.
	throwTypeUser

	// throwTypeRuntime is a throw due to a problem with Go itself.
	//
	// These throws include as much information as possible to aid in
	// debugging the runtime, including runtime frames, system goroutines,
	// and frame metadata.
	throwTypeRuntime
)

func throw(s string) {

}

func fatal(s string) {
}

var panicking atomic.Uint32

var runningPanicDefers atomic.Uint32

var deadlock mutex

var paniclk mutex

func shouldPushSigpanic(gp *g, pc, lr uintptr) bool {
	return false
}

func isAbortPC(pc uintptr) bool {
	return false
}
func startpanic_m() bool {
	return false
}

func canpanic() bool {
	return false
}

func panicmem() {

}

func panicmemAddr(addr uintptr) {

}
func panicshift() {

}

func panicdivide() {

}

func panicoverflow() {

}

func panicfloat() {

}

// builtin
func gopanic(e any) {}

// builtin
func deferreturn() {}

// builtin
func gorecover(argp uintptr) any {
	return nil
}

// builtin
func goPanicIndex(x int, y int) {}

// builtin
func goPanicIndexU(x uint, y int) {

}

// builtin
func goPanicSliceAlen(x int, y int) {

}

// builtin
func goPanicSliceAlenU(x uint, y int) {
}

// builtin
func goPanicSliceAcap(x int, y int) {

}

// builtin
func goPanicSliceAcapU(x uint, y int) {

}

// builtin
func goPanicSliceB(x int, y int) {

}

// builtin
func goPanicSliceBU(x uint, y int) {

}

// builtin
func goPanicSlice3Alen(x int, y int) {

}

// builtin
func goPanicSlice3AlenU(x uint, y int) {

}

// builtin
func goPanicSlice3Acap(x int, y int) {

}

// builtin
func goPanicSlice3AcapU(x uint, y int) {

}

// builtin
func goPanicSlice3B(x int, y int) {

}

// builtin
func goPanicSlice3BU(x uint, y int) {

}

// builtin
func goPanicSlice3C(x int, y int) {

}

// builtin
func goPanicSlice3CU(x uint, y int) {

}

// builtin
func goPanicSliceConvert(x int, y int) {

}
