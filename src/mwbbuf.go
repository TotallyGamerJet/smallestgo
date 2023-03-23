package runtime

//go:nowritebarrierrec
//go:nosplit
func wbBufFlush(dst *uintptr, src uintptr) {}

type wbBuf struct {
	// next points to the next slot in buf. It must not be a
	// pointer type because it can point past the end of buf and
	// must be updated without write barriers.
	//
	// This is a pointer rather than an index to optimize the
	// write barrier assembly.
	next uintptr

	// end points to just past the end of buf. It must not be a
	// pointer type because it points past the end of buf and must
	// be updated without write barriers.
	end uintptr

	// buf stores a series of pointers to execute write barriers
	// on. This must be a multiple of wbBufEntryPointers because
	// the write barrier only checks for overflow once per entry.
	buf [1]uintptr
}

//go:nowritebarrierrec
//go:nosplit
func (b *wbBuf) putFast(old, new uintptr) bool { return false }

//go:nowritebarrierrec
//go:systemstack
func wbBufFlush1(pp *p) {}

func (b *wbBuf) reset() {}
