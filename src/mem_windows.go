package runtime

import "unsafe"

func sysFreeOS(v unsafe.Pointer, n uintptr) {}
func sysAllocOS(n uintptr) unsafe.Pointer   { return nil }
