package runtime

import "unsafe"

type hchan struct {
	qcount   uint           // total data in the queue
	dataqsiz uint           // size of the circular queue
	buf      unsafe.Pointer // points to an array of dataqsiz elements
	elemsize uint16
	closed   uint32
	elemtype *_type // element type
	sendx    uint   // send index
	recvx    uint   // receive index
	recvq    waitq  // list of recv waiters
	sendq    waitq  // list of send waiters

	// lock protects all fields in hchan, as well as several
	// fields in sudogs blocked on this channel.
	//
	// Do not change another G's status while holding this lock
	// (in particular, do not ready a G), as this can deadlock
	// with stack shrinking.
	lock mutex
}

type waitq struct {
	first *sudog
	last  *sudog
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
