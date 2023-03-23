package runtime

import "unsafe"

// builtin
func memequal0(p, q unsafe.Pointer) bool {
	return true
}

// builtin
func memequal8(p, q unsafe.Pointer) bool {
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
func strequal(p, q unsafe.Pointer) bool {
	return false
}
