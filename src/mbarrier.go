package runtime

import "unsafe"

//go:nosplit
func typedmemmove(typ *_type, dst, src unsafe.Pointer) {}

// builtin
//
//go:nosplit
func typedslicecopy(typ *_type, dstPtr unsafe.Pointer, dstLen int, srcPtr unsafe.Pointer, srcLen int) int {
	return 0
}

//go:nosplit
func typedmemclr(typ *_type, ptr unsafe.Pointer) {}

//go:nosplit
func memclrHasPointers(ptr unsafe.Pointer, n uintptr) {}
