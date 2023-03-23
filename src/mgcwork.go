package runtime

type gcWork struct {
	// wbuf1 and wbuf2 are the primary and secondary work buffers.
	//
	// This can be thought of as a stack of both work buffers'
	// pointers concatenated. When we pop the last pointer, we
	// shift the stack up by one work buffer by bringing in a new
	// full buffer and discarding an empty one. When we fill both
	// buffers, we shift the stack down by one work buffer by
	// bringing in a new empty buffer and discarding a full one.
	// This way we have one buffer's worth of hysteresis, which
	// amortizes the cost of getting or putting a work buffer over
	// at least one buffer of work and reduces contention on the
	// global work lists.
	//
	// wbuf1 is always the buffer we're currently pushing to and
	// popping from and wbuf2 is the buffer that will be discarded
	// next.
	//
	// Invariant: Both wbuf1 and wbuf2 are nil or neither are.
	wbuf1, wbuf2 *workbuf

	// Bytes marked (blackened) on this gcWork. This is aggregated
	// into work.bytesMarked by dispose.
	bytesMarked uint64

	// Heap scan work performed on this gcWork. This is aggregated into
	// gcController by dispose and may also be flushed by callers.
	// Other types of scan work are flushed immediately.
	heapScanWork int64

	// flushedWork indicates that a non-empty work buffer was
	// flushed to the global work list since the last gcMarkDone
	// termination check. Specifically, this indicates that this
	// gcWork may have communicated work to another gcWork.
	flushedWork bool
}

type workbuf struct {
	workbufhdr
	// account for the above fields
	obj [0]uintptr
}

type workbufhdr struct {
	node lfnode // must be first
	nobj int
}

const (
	_WorkbufSize = 2048 // in bytes; larger values result in less contention

	// workbufAlloc is the number of bytes to allocate at a time
	// for new workbufs. This must be a multiple of pageSize and
	// should be a multiple of _WorkbufSize.
	//
	// Larger values reduce workbuf allocation overhead. Smaller
	// values reduce heap fragmentation.
	workbufAlloc = 32 << 10
)

//go:nowritebarrierrec
func (w *gcWork) dispose() {}

//go:nowritebarrierrec
func (w *gcWork) empty() bool { return false }
func prepareFreeWorkbufs()    {}

//go:nowritebarrier
func putempty(b *workbuf) {}

//go:nowritebarrierrec
func (w *gcWork) tryGetFast() uintptr     { return 0 }
func freeSomeWbufs(preemptible bool) bool { return false }

//go:nowritebarrierrec
func (w *gcWork) tryGet() uintptr { return 0 }

//go:nowritebarrierrec
func (w *gcWork) balance() {}

//go:nowritebarrier
func getempty() *workbuf { return nil }

//go:nowritebarrierrec
func (w *gcWork) putFast(obj uintptr) bool { return false }

//go:nowritebarrierrec
func (w *gcWork) put(obj uintptr) {}

//go:nowritebarrierrec
func (w *gcWork) putBatch(obj []uintptr) {}
