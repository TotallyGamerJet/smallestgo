package runtime

import "unsafe"

func checkptrStraddles(ptr unsafe.Pointer, size uintptr) bool { return false }
