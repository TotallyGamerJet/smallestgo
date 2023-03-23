package runtime

const pageCachePages = 8 * 8

type pageCache struct {
	base  uintptr // base address of the chunk
	cache uint64  // 64-bit bitmap representing free pages (1 means free)
	scav  uint64  // 64-bit bitmap representing scavenged pages (1 means scavenged)
}

//go:systemstack
func (p *pageAlloc) allocToCache() pageCache { return pageCache{} }

func (c *pageCache) empty() bool { return false }

func (c *pageCache) alloc(npages uintptr) (uintptr, uintptr) { return 0, 0 }

//go:systemstack
func (c *pageCache) flush(p *pageAlloc) {}
