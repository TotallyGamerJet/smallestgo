package runtime

import "unsafe"

type cgoCallers [32]uintptr

var ncgocall uint64 // number of cgo calls in total for dead m

//go:nosplit
//go:nowritebarrierrec
func cgoIsGoPointer(p unsafe.Pointer) bool {
	return false
}

func cgoInRange(p unsafe.Pointer, start, end uintptr) bool {
	return false
}

func cgocall(fn, arg unsafe.Pointer) int32 { return 0 }
