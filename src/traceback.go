package runtime

import "unsafe"

func callers(skip int, pcbuf []uintptr) int {
	return 0
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
