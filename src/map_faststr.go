package runtime

import "unsafe"

func mapaccess2_faststr(t *maptype, h *hmap, ky string) (unsafe.Pointer, bool) { return nil, false }

func mapassign_faststr(t *maptype, h *hmap, s string) unsafe.Pointer { return nil }

func mapdelete_faststr(t *maptype, h *hmap, ky string) {}
