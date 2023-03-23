package runtime

import "unsafe"

// builtin
func mapaccess2_fast32(t *maptype, h *hmap, key uint32) (unsafe.Pointer, bool) {
	return nil, false
}

// builtin
func mapaccess1_fast32(t *maptype, h *hmap, key uint32) unsafe.Pointer { return nil }

// builtin
func mapassign_fast32(t *maptype, h *hmap, key uint32) unsafe.Pointer { return nil }
