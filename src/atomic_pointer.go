package runtime

import "unsafe"

//go:nosplit
func atomicstorep(ptr unsafe.Pointer, new unsafe.Pointer) {}
