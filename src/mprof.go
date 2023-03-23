package runtime

import (
	"runtime/internal/atomic"
	"runtime/internal/sys"
	"unsafe"
)

type bucket struct {
	_       sys.NotInHeap
	next    *bucket
	allnext *bucket
	typ     bucketType // memBucket or blockBucket (includes mutexProfile)
	hash    uintptr
	size    uintptr
	nstk    uintptr
}

var goroutineProfile = struct {
	sema    uint32
	active  bool
	offset  atomic.Int64
	records []StackRecord
	labels  []unsafe.Pointer
}{
	sema: 1,
}

func tryRecordGoroutineProfile(gp1 *g, yield func()) {}

var disableMemoryProfiling bool
var MemProfileRate int = 512 * 1024

type goroutineProfileStateHolder atomic.Uint32

//go:yeswritebarrierrec
func tryRecordGoroutineProfileWB(gp1 *g) {}

type goroutineProfileState uint32

const (
	goroutineProfileAbsent goroutineProfileState = iota
	goroutineProfileInProgress
	goroutineProfileSatisfied
)

func (p *goroutineProfileStateHolder) Store(value goroutineProfileState) {}

type StackRecord struct {
	Stack0 [32]uintptr // stack trace for this record; ends at first 0 entry
}

const (
	// profile types
	memProfile bucketType = 1 + iota
	blockProfile
	mutexProfile

	// size of bucket hash table
	buckHashSize = 179999

	// max depth of stack to record in bucket
	maxStack = 32
)

type bucketType int
