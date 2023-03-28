package runtime

import "unsafe"

// builtin
func memequal0(p, q unsafe.Pointer) bool {
	return true
}

// builtin
func memequal8(p, q unsafe.Pointer) bool {
	return false
}

// builtin
func memequal32(p, q unsafe.Pointer) bool {
	return false
}

// builtin
func memequal64(p, q unsafe.Pointer) bool {
	return false
}

// builtin
func strequal(p, q unsafe.Pointer) bool {
	return false
}

// dwarf: missing type: type:runtime.hchan
type hchan struct{}

// dwarf: missing type: type:runtime.waitq
type waitq struct{}

var buildVersion string

type hmap struct {
	buckets    uintptr // array of 2^B Buckets. may be nil if count==0.
	oldbuckets uintptr // previous bucket array of half the size, non-nil only when growing
}
type bmap struct{}

var writeBarrier struct{}

type ptabEntry struct{}

// set using cmd/go/internal/modload.ModInfoProg
var modinfo string

//go:linkname runtime_inittask runtime..inittask
var runtime_inittask initTask

//go:linkname main_inittask main..inittask
var main_inittask initTask

//go:linkname main_main main.main
func main_main()

func schedinit() {
	modules := new([]*moduledata)
	for md := &firstmoduledata; md != nil; md = md.next {
		*modules = append(*modules, md)
	}

	if buildVersion == "" {
		// Condition should never trigger. This code just serves
		// to ensure runtime·buildVersion is kept in the resulting binary.
		buildVersion = "unknown"
	}
	if len(modinfo) == 1 {
		// Condition should never trigger. This code just serves
		// to ensure runtime·modinfo is kept in the resulting binary.
		modinfo = ""
	}
	main := main_main
	main()
}

func (mp *m) becomeSpinning() {
}

// A gQueue is a dequeue of Gs linked through g.schedlink. A G can only
// be on one gQueue or gList at a time.
type gQueue struct {
	head guintptr
	tail guintptr
}

// empty reports whether q is empty.
func (q *gQueue) empty() bool {
	return false
}

// push adds gp to the head of q.
func (q *gQueue) push(gp *g) {
}

// pushBack adds gp to the tail of q.
func (q *gQueue) pushBack(gp *g) {
}

// pushBackAll adds all Gs in q2 to the tail of q. After this q2 must
// not be used.
func (q *gQueue) pushBackAll(q2 gQueue) {
}

// pop removes and returns the head of queue q. It returns nil if
// q is empty.
func (q *gQueue) pop() *g {
	return nil
}

// popList takes all Gs in q and returns them as a gList.
func (q *gQueue) popList() gList {
	return gList{}
}

// A gList is a list of Gs linked through g.schedlink. A G can only be
// on one gQueue or gList at a time.
type gList struct{}

// An initTask represents the set of initializations that need to be done for a package.
// Keep in sync with ../../test/initempty.go:initTask
type initTask struct{}

type sudog struct{}

type itab struct{}

type mutex struct {
	// Empty struct if lock ranking is disabled, otherwise includes the lock rank
	// Futex-based impl treats it as uint32 key,
	// while sema-based impl as M* waitm.
	// Used to be a union, but unions break precise GC.
	key uintptr
}

type g struct{}

type puintptr uintptr

type lfnode struct{}

type muintptr uintptr

type m struct{}

type guintptr uintptr

func (_ guintptr) set(_ *g) {}

type _func struct{}

type iface struct{}

type eface struct{}

type slice struct {
	array uintptr
}

// builtin
func growslice(oldPtr unsafe.Pointer, newLen, oldCap, num int, et *_type) slice {
	return slice{}
}

// Frames may be used to get function/file/line information for a
// slice of PC values returned by Callers.
type Frames struct{}

// Frame is the information returned by Frames for each call frame.
type Frame struct{}

// A Func represents a Go function in the running binary.
type Func struct{}

// pcHeader holds data used by the pclntab lookups.
type pcHeader struct{}

// moduledata records information about the layout of the executable
// image. It is written by the linker. Any changes here must be
// matched changes to the code in cmd/link/internal/ld/symtab.go:symtab.
// moduledata is stored in statically allocated non-pointer memory;
// none of the pointers here are visible to the garbage collector.
type moduledata struct {
	pcHeader     *pcHeader
	funcnametab  []byte
	cutab        []uint32
	filetab      []byte
	pctab        []byte
	pclntable    []byte
	ftab         []functab
	findfunctab  uintptr
	minpc, maxpc uintptr

	text, etext           uintptr
	noptrdata, enoptrdata uintptr
	data, edata           uintptr
	bss, ebss             uintptr
	noptrbss, enoptrbss   uintptr
	covctrs, ecovctrs     uintptr
	end, gcdata, gcbss    uintptr
	types, etypes         uintptr
	rodata                uintptr
	gofunc                uintptr // go.func.*

	textsectmap []textsect
	typelinks   []int32 // offsets from types
	itablinks   []*itab

	ptab []ptabEntry

	pluginpath string
	pkghashes  []modulehash

	modulename   string
	modulehashes []modulehash

	hasmain uint8 // 1 if module contains the main function, 0 otherwise

	gcdatamask, gcbssmask bitvector

	typemap map[typeOff]*_type // offset to *_rtype in previous module

	bad bool // module failed to load and should be ignored

	next *moduledata
}

type modulehash struct{}

var firstmoduledata moduledata // linker symbol

type functab struct{}

// Mapping information for secondary text sections

type textsect struct{}

type (
	bitvector         struct{}
	stackObjectRecord struct{}
)

// Variant with *byte pointer type for DWARF debugging.
type stringStructDWARF struct {
	str *byte
	len int
}

func gcWriteBarrier() {}

func memequal_varlen(a, b unsafe.Pointer) bool { return false }

type _type struct{}

type typeOff int32

// these must stick around bc the linker will complain:
// dwarf: missing type: type:runtime.imethod
type imethod struct{}

type interfacetype struct{}

type maptype struct{}

type arraytype struct{}

type chantype struct{}

type slicetype struct{}

type functype struct{}

type ptrtype struct{}

type structfield struct{}

type structtype struct{}
