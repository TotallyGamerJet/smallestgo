package runtime

import "unsafe"
import "internal/abi"

// Should be a built-in for unsafe.Pointer?
//
//go:nosplit
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

// getg returns the pointer to the current g.
// The compiler rewrites calls to this function into instructions
// that fetch the g directly (from TLS or from the dedicated register).
func getg() *g

//go:noescape
func memclrNoHeapPointers(ptr unsafe.Pointer, n uintptr)

//go:noescape
func memmove(to, from unsafe.Pointer, n uintptr)

func cgocallback(fn, frame, ctxt uintptr)

//go:noescape
func getcallerpc() uintptr

//go:noescape
func systemstack(fn func())

//go:noescape
func asmcgocall(fn, arg unsafe.Pointer) int32

func procyield(cycles uint32)

func mcall(fn func(*g))

func gcWriteBarrier() {}

func memequal(a, b unsafe.Pointer, size uintptr) bool { return false }

func memequal_varlen(a, b unsafe.Pointer) bool { return false }

//go:noescape
func reflectcall(stackArgsType *_type, fn, stackArgs unsafe.Pointer, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs)

func gogo(buf *gobuf)

func asminit()
func setg(gg *g)
func morestack()
func breakpoint()

func checkASM() bool
