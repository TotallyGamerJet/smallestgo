package runtime

import "unsafe"

const (
	maxCPUProfStack = 64

	// profBufWordCount is the size of the CPU profile buffer's storage for the
	// header and stack of each sample, measured in 64-bit words. Every sample
	// has a required header of two words. With a small additional header (a
	// word or two) and stacks at the profiler's maximum length of 64 frames,
	// that capacity can support 1900 samples or 19 thread-seconds at a 100 Hz
	// sample rate, at a cost of 1 MiB.
	profBufWordCount = 1 << 17
	// profBufTagCount is the size of the CPU profile buffer's storage for the
	// goroutine tags associated with each sample. A capacity of 1<<14 means
	// room for 16k samples, or 160 thread-seconds at a 100 Hz sample rate.
	profBufTagCount = 1 << 14
)

var cpuprof cpuProfile

type cpuProfile struct {
	lock mutex
	on   bool     // profiling is on
	log  *profBuf // profile events written here

	// extra holds extra stacks accumulated in addNonGo
	// corresponding to profiling signals arriving on
	// non-Go-created threads. Those stacks are written
	// to log the next time a normal Go thread gets the
	// signal handler.
	// Assuming the stacks are 2 words each (we don't get
	// a full traceback from those threads), plus one word
	// size for framing, 100 Hz profiling would generate
	// 300 words per second.
	// Hopefully a normal Go thread will get the profiling
	// signal at least once every few seconds.
	extra      [1000]uintptr
	numExtra   int
	lostExtra  uint64 // count of frames lost because extra is full
	lostAtomic uint64 // count of frames lost because of being in atomic64 on mips/arm; updated racily
}

//go:nowritebarrierrec
func (p *cpuProfile) add(tagPtr *unsafe.Pointer, stk []uintptr) {}
