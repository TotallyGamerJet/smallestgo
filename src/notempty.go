package runtime

import (
	"unsafe"
	_ "unsafe"
)

func schedinit() {
	modules := new([]*moduledata)
	for md := &firstmoduledata; md != nil; md = (*moduledata)(unsafe.Pointer(md.next)) {
		*modules = append(*modules, md)
	}
}

// moduledata records information about the layout of the executable
// image. It is written by the linker. Any changes here must be
// matched changes to the code in cmd/link/internal/ld/symtab.go:symtab.
// moduledata is stored in statically allocated non-pointer memory;
// none of the pointers here are visible to the garbage collector.
type moduledata struct {
	_    [505 - 8]byte
	next *struct{}
}

var firstmoduledata moduledata // linker symbol
