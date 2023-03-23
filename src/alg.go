package runtime

import "unsafe"

var hashkey [4]uintptr

var useAeshash bool

const hashRandomBytes = 128

// used in asm_{386,amd64,arm64}.s to seed the hash function
var aeskeysched [hashRandomBytes]byte

func readUnaligned32(p unsafe.Pointer) uint32 { return 0 }

func alginit() {}

func int64Hash(i uint64, seed uintptr) uintptr { return 0 }

func memhash(p unsafe.Pointer, h, s uintptr) uintptr

func readUnaligned64(p unsafe.Pointer) uint64 { return 0 }

// builtin
func memhash0(p unsafe.Pointer, h uintptr) uintptr {
	return h
}

// builtin
func memhash8(p unsafe.Pointer, h uintptr) uintptr {
	return 0
}

// builtin
func memhash16(p unsafe.Pointer, h uintptr) uintptr {
	return 0
}

// builtin
func memhash128(p unsafe.Pointer, h uintptr) uintptr {
	return 0
}

// builtin
func memequal0(p, q unsafe.Pointer) bool {
	return true
}

// builtin
func memequal8(p, q unsafe.Pointer) bool {
	return false
}

// builtin
func memequal16(p, q unsafe.Pointer) bool {
	return false
}

// builtin
func memequal32(p, q unsafe.Pointer) bool {
	return false
}

// builtin
func memequal64(p, q unsafe.Pointer) bool {
	return false
}

// builtin
func memequal128(p, q unsafe.Pointer) bool {
	return false
}

// builtin
func f32equal(p, q unsafe.Pointer) bool {
	return false
}

// builtin
func f64equal(p, q unsafe.Pointer) bool {
	return false
}

// builtin
func c64equal(p, q unsafe.Pointer) bool {
	return false
}

// builtin
func c128equal(p, q unsafe.Pointer) bool {
	return false
}

// builtin
func strequal(p, q unsafe.Pointer) bool {
	return false
}

// builtin
func nilinterequal(p, q unsafe.Pointer) bool {
	return false
}

func efaceeq(t *_type, x, y unsafe.Pointer) bool { return false }
