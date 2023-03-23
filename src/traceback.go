package runtime

import "unsafe"

const usesLR = false

func callers(skip int, pcbuf []uintptr) int {
	return 0
}

func isSystemGoroutine(gp *g, fixed bool) bool {
	return false
}

func gentraceback(pc0, sp0, lr0 uintptr, gp *g, skip int, pcbuf *uintptr, max int, callback func(*stkframe, unsafe.Pointer) bool, v unsafe.Pointer, flags uint) int {
	return 0
}

func gcallers(gp *g, skip int, pcbuf []uintptr) int {
	return 0
}

func goroutineheader(gp *g) {

}

func tracebackothers(me *g) {

}

func traceback(pc, sp, lr uintptr, gp *g) {

}

func tracebacktrap(pc, sp, lr uintptr, gp *g) {

}

var cgoSymbolizer unsafe.Pointer

type cgoSymbolizerArg struct {
	pc       uintptr
	file     *byte
	lineno   uintptr
	funcName *byte
	entry    uintptr
	more     uintptr
	data     uintptr
}

func callCgoSymbolizer(arg *cgoSymbolizerArg) {}

func elideWrapperCalling(id funcID) bool { return false }
