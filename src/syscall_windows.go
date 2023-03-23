package runtime

import (
	"unsafe"
	_ "unsafe"
)

type callbackArgs struct {
	index uintptr
	// args points to the argument block.
	//
	// For cdecl and stdcall, all arguments are on the stack.
	//
	// For fastcall, the trampoline spills register arguments to
	// the reserved spill slots below the stack arguments,
	// resulting in a layout equivalent to stdcall.
	//
	// For arm, the trampoline stores the register arguments just
	// below the stack arguments, so again we can treat it as one
	// big stack arguments frame.
	args unsafe.Pointer
	// Below are out-args from callbackWrap
	result uintptr
	retPop uintptr // For 386 cdecl, how many bytes to pop on return
}

const _LOAD_LIBRARY_SEARCH_SYSTEM32 = 0x00000800

//go:linkname compileCallback syscall.compileCallback
func compileCallback(fn eface, cdecl bool) (code uintptr) {
	return 0
}
