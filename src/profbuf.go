package runtime

import (
	"sync/atomic"
	"unsafe"
)

type profBuf struct {
	//r, w         profAtomic
	overflow     atomic.Uint64
	overflowTime atomic.Uint64
	eof          atomic.Uint32

	// immutable (excluding slice content)
	hdrsize uintptr
	data    []uint64
	tags    []unsafe.Pointer

	// owned by reader
	//rNext       profIndex
	overflowBuf []uint64 // for use by reader to return overflow record
	wait        note
}

func (b *profBuf) write(tagPtr *unsafe.Pointer, now int64, hdr []uint64, stk []uintptr)       {}
func (b *profBuf) close()                                                                     {}
func (b *profBuf) read(mode profBufReadMode) (data []uint64, tags []unsafe.Pointer, eof bool) { return }
func newProfBuf(hdrsize, bufwords, tags int) *profBuf                                         { return nil }

type profBufReadMode int

const (
	profBufBlocking profBufReadMode = iota
	profBufNonBlocking
)
