package runtime

import "unsafe"

type slice struct {
	array uintptr
}

// builtin
func growslice(oldPtr unsafe.Pointer, newLen, oldCap, num int, et *_type) slice {
	return slice{}
}
