package runtime

type pallocData struct {
	pallocBits
	scavenged pageBits
}

type pageBits [128 / 64]uint64

type pallocBits pageBits

func (b *pallocBits) summarize() pallocSum { return 0 }

func (b *pageBits) setRange(i, n uint) {}

func (b *pageBits) popcntRange(i, n uint) (s uint) { return }

func (b *pallocBits) allocRange(i, n uint)                             {}
func (b *pallocBits) allocAll()                                        {}
func (b *pallocBits) find(npages uintptr, searchIdx uint) (uint, uint) { return 0, 0 }
func (b *pallocBits) free1(i uint)                                     {}
func (b *pallocBits) free(i, n uint)                                   {}
func (b *pallocBits) freeAll()                                         {}
