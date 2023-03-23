package runtime

import "unsafe"

//go:nosplit
//go:nowritebarrier
func cgoCheckMemmove(typ *_type, dst, src unsafe.Pointer, off, size uintptr) {

}

func cgoCheckSliceCopy(typ *_type, dst, src unsafe.Pointer, n int) {

}

func cgoCheckWriteBarrier(dst *uintptr, src uintptr) {

}
