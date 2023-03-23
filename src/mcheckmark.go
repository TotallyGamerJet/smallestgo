package runtime

import "runtime/internal/sys"

type checkmarksMap struct {
	_ sys.NotInHeap
	b [heapArenaBytes / 8 / 8]uint8
}

func startCheckmarks() {}
func endCheckmarks()   {}
