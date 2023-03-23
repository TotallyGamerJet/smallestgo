package runtime

import "unsafe"

type hchan struct {
}

type waitq struct {
}

func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool    { return false }
func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) { return }

// builtin
//
//go:nosplit
func chanrecv1(c *hchan, elem unsafe.Pointer) {}

// builtin
//
//go:nosplit
func chansend1(c *hchan, elem unsafe.Pointer) {}

// builtin
func makechan(t *chantype, size int) *hchan { return nil }

// builtin
func closechan(c *hchan) {}
