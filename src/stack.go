package runtime

import _ "unsafe"

const _StackGuard = 0

const stackDebug = 0

type bitvector struct {
	n        int32 // # of bits
	bytedata *uint8
}

func newstack() {}

type stackObjectRecord struct {
	// offset in frame
	// if negative, offset from varp
	// if non-negative, offset from argp
	off       int32
	size      int32
	_ptrdata  int32  // ptrdata, or -ptrdata is GC prog is used
	gcdataoff uint32 // offset to gcdata from moduledata.rodata
}

const (
	uintptrMask = 1<<(8*8) - 1

	// The values below can be stored to g.stackguard0 to force
	// the next stack check to fail.
	// These are all larger than any real SP.

	// Goroutine preemption request.
	// 0xfffffade in hex.
	stackPreempt = uintptrMask & -1314

	// Thread is forking. Causes a split stack check failure.
	// 0xfffffb2e in hex.
	stackFork = uintptrMask & -1234

	// Force a stack movement. Used for debugging.
	// 0xfffffeed in hex.
	stackForceMove = uintptrMask & -275

	// stackPoisonMin is the lowest allowed stack poison value.
	stackPoisonMin = uintptrMask & -4096
)

//go:systemstack
func stackcache_clear(c *mcache)  {}
func gcComputeStartingStackSize() {}
func freeStackSpans()             {}

var maxstacksize uintptr = 1 << 20

//go:systemstack
func stackfree(stk stack)  {}
func round2(x int32) int32 { return 0 }

const _StackSystem = 512 * 8
const _StackMin = 2048

var maxstackceiling = maxstacksize

func stackinit() {}

//go:systemstack
func stackalloc(n uint32) stack { return stack{} }

var startingStackSize uint32 = _FixedStack

const _FixedStack = 0x1000

func gostartcallfn(gobuf *gobuf, fv *funcval) {}
func (bv *bitvector) ptrbit(i uintptr) uint8  { return 0 }

// builtin
//
//go:nosplit
//go:linkname morestackc
func morestackc() {}
