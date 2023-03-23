package runtime

import "unsafe"

type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}

type notInHeapSlice struct {
	array *notInHeap
	len   int
	cap   int
}

func makeslicecopy(et *_type, tolen int, fromlen int, from unsafe.Pointer) unsafe.Pointer { return nil }

func slicecopy(toPtr unsafe.Pointer, toLen int, fromPtr unsafe.Pointer, fromLen int, width uintptr) int {
	return 0
}

// builtin
func growslice(oldPtr unsafe.Pointer, newLen, oldCap, num int, et *_type) slice {
	return slice{}
}

func makeslice(et *_type, len, cap int) unsafe.Pointer { return nil }
