package runtime

import "unsafe"

type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}

type notInHeap struct{}

type notInHeapSlice struct {
	array *notInHeap
	len   int
	cap   int
}

// builtin
func growslice(oldPtr unsafe.Pointer, newLen, oldCap, num int, et *_type) slice {
	return slice{}
}
