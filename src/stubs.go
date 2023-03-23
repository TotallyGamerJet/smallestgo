package runtime

import "unsafe"

func gcWriteBarrier() {}

func memequal_varlen(a, b unsafe.Pointer) bool { return false }
