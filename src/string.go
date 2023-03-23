package runtime

import "unsafe"

const tmpStringBufSize = 32

type tmpBuf [tmpStringBufSize]byte

const (
	maxUint64 = ^uint64(0)
	maxInt64  = int64(maxUint64 >> 1)
)

type stringStruct struct {
	str unsafe.Pointer
	len int
}

func stringStructOf(sp *string) *stringStruct {
	return (*stringStruct)(unsafe.Pointer(sp))
}

func rawstring(size int) (s string, b []byte) {
	return "", nil
}

func slicebytetostringtmp(ptr *byte, n int) string { return "" }
func atoi32(s string) (int32, bool)                { return 0, false }
func hasPrefix(s, prefix string) bool {
	return false
}

//go:nosplit
func gostringnocopy(str *byte) string { return "" }
func gostringw(strw *uint16) string   { return "" }

//go:nosplit
func findnull(s *byte) int { return 0 }

func parseByteCount(s string) (int64, bool) { return 0, false }

//go:linkname gostring
func gostring(p *byte) string { return "" }

// builtin
func slicebytetostring(buf *tmpBuf, ptr *byte, n int) string { return "" }

// Variant with *byte pointer type for DWARF debugging.
type stringStructDWARF struct {
	str *byte
	len int
}
