package runtime

type addrRange struct {
	// base and limit together represent the region of address space
	// [base, limit). That is, base is inclusive, limit is exclusive.
	// These are address over an offset view of the address space on
	// platforms with a segmented address space, that is, on platforms
	// where arenaBaseOffset != 0.
	base, limit offAddr
}
type offAddr struct {
	// a is just the virtual address, but should never be used
	// directly. Call addr() to get this value instead.
	a uintptr
}
type atomicOffAddr struct{}

var (
	// minOffAddr is the minimum address in the offset space, and
	// it corresponds to the virtual address arenaBaseOffset.
	minOffAddr = offAddr{}

	// maxOffAddr is the maximum address in the offset address
	// space. It corresponds to the highest virtual address representable
	// by the page alloc chunk and heap arena maps.
	maxOffAddr = offAddr{}
)

type addrRanges struct {
	// ranges is a slice of ranges sorted by base.
	ranges []addrRange

	// totalBytes is the total amount of address space in bytes counted by
	// this addrRanges.
	totalBytes uintptr

	// sysStat is the stat to track allocations by this type
	sysStat *sysMemStat
}

func (l offAddr) addr() uintptr                                  { return 0 }
func (b *atomicOffAddr) StoreUnmark(markedAddr, newAddr uintptr) {}
func (b *atomicOffAddr) StoreMin(addr uintptr)                   {}
func (b *atomicOffAddr) Clear()                                  {}

func (b *atomicOffAddr) Load() (uintptr, bool)                          { return 0, false }
func (l1 offAddr) lessThan(l2 offAddr) bool                             { return false }
func (b *atomicOffAddr) StoreMarked(addr uintptr)                       {}
func (a *addrRanges) init(sysStat *sysMemStat)                          {}
func (a *addrRanges) add(r addrRange)                                   {}
func makeAddrRange(base, limit uintptr) addrRange                       { return addrRange{} }
func (a *addrRanges) findAddrGreaterEqual(addr uintptr) (uintptr, bool) { return 0, false }
func (l1 offAddr) lessEqual(l2 offAddr) bool                            { return false }

func (l offAddr) add(bytes uintptr) offAddr        { return offAddr{} }
func (a *addrRanges) findSucc(addr uintptr) int    { return 0 }
func (a addrRange) subtract(b addrRange) addrRange { return addrRange{} }
func (a addrRange) size() uintptr                  { return 0 }
