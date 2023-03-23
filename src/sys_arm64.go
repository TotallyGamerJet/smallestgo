package runtime

import "unsafe"

func gostartcall(buf *gobuf, fn, ctxt unsafe.Pointer) {}
