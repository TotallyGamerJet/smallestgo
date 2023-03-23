package runtime

type sudog struct {
}

type itab struct {
}

type mutex struct {
	// Empty struct if lock ranking is disabled, otherwise includes the lock rank
	// Futex-based impl treats it as uint32 key,
	// while sema-based impl as M* waitm.
	// Used to be a union, but unions break precise GC.
	key uintptr
}

type g struct {
}

type puintptr uintptr

type lfnode struct {
}

type muintptr uintptr

type m struct {
}

type guintptr uintptr

func (_ guintptr) set(_ *g) {}

type _func struct {
}

type iface struct {
}

type eface struct {
}
