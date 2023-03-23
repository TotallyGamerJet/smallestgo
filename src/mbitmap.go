package runtime

import "unsafe"

type markBits struct {
	bytep *uint8
	mask  uint8
	index uintptr
}

type writeHeapBits struct {
	addr  uintptr // address that the low bit of mask represents the pointer state of.
	mask  uintptr // some pointer bits starting at the address addr.
	valid uintptr // number of bits in buf that are valid (including low)
	low   uintptr // number of low-order bits to not overwrite
}

type heapBits struct {
	// heapBits will report on pointers in the range [addr,addr+size).
	// The low bit of mask contains the pointerness of the word at addr
	// (assuming valid>0).
	addr, size uintptr

	// The next few pointer bits representing words starting at addr.
	// Those bits already returned by next() are zeroed.
	mask uintptr
	// Number of bits in mask that are valid. mask is always less than 1<<valid.
	valid uintptr
}

func runGCProg(prog, dst *byte) uintptr { return 0 }

func writeHeapBitsForAddr(addr uintptr) (h writeHeapBits) {
	return
}

//go:nosplit
func heapBitsForAddr(addr, size uintptr) heapBits {
	return heapBits{}
}

const ptrBits = 8 * 8

//go:nowritebarrier
//go:nosplit
func addb(p *byte, n uintptr) *byte {
	// Note: wrote out full expression instead of calling add(p, n)
	// to reduce the number of temporaries generated by the
	// compiler for this trivial expression during inlining.
	return (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + n))
}

func readUintptr(p *byte) uintptr {
	return 0
}

//go:nosplit
func (h heapBits) next() (heapBits, uintptr) {
	return heapBits{}, 0
}

func (h writeHeapBits) write(bits, valid uintptr) writeHeapBits {
	return writeHeapBits{}
}

func (h writeHeapBits) pad(size uintptr) writeHeapBits { return writeHeapBits{} }
func (h writeHeapBits) flush(addr, size uintptr)       {}

//go:nosplit
func typeBitsBulkBarrier(typ *_type, dst, src, size uintptr) {}

//go:nosplit
func findObject(p, refBase, refOff uintptr) (base uintptr, s *mspan, objIndex uintptr) {
	return 0, nil, 0
}

//go:nosplit
func bulkBarrierPreWrite(dst, src, size uintptr) {}

func heapBitsSetType(x, size, dataSize uintptr, typ *_type) {}

func (m markBits) isMarked() bool                             { return false }
func materializeGCProg(ptrdata uintptr, prog *byte) *struct{} { return nil }
func dematerializeGCProg(s *struct{})                         {}

//go:nosplit
func (h heapBits) nextFast() (heapBits, uintptr) { return heapBits{}, 0 }

func progToPointerMask(prog *byte, size uintptr) bitvector { return bitvector{} }
func (s *mspan) isFree(index uintptr) bool                 { return false }

func (s *mspan) nextFreeIndex() uintptr                     { return 0 }
func (s *mspan) initHeapBits(forceClear bool)               {}
func (s *mspan) refillAllocCache(whichByte uintptr)         {}
func (s *mspan) divideByElemSize(n uintptr) uintptr         { return 0 }
func (s *mspan) markBitsForIndex(objIndex uintptr) markBits { return markBits{} }

//go:nosplit
func (s *mspan) allocBitsForIndex(allocBitIndex uintptr) markBits { return markBits{} }
func (s *mspan) countAlloc() int                                  { return 0 }
func (m markBits) setMarkedNonAtomic()                            {}
func (s *mspan) markBitsForBase() markBits                        { return markBits{} }
func (m *markBits) advance()                                      {}
