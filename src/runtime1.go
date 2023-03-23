package runtime

//go:nosplit
func acquirem() *m {
	return nil
}

//go:nosplit
func releasem(mp *m) {
}

var debug struct {
	cgocheck           int32
	clobberfree        int32
	efence             int32
	gccheckmark        int32
	gcpacertrace       int32
	gcshrinkstackoff   int32
	gcstoptheworld     int32
	gctrace            int32
	invalidptr         int32
	madvdontneed       int32 // for Linux; issue 28466
	scavtrace          int32
	scheddetail        int32
	schedtrace         int32
	tracebackancestors int32
	asyncpreemptoff    int32
	harddecommit       int32
	adaptivestackstart int32

	// debug.malloc is used as a combined debug check
	// in the malloc function and should be set
	// if any of the below debug options is != 0.
	malloc         bool
	allocfreetrace int32
	inittrace      int32
	sbrk           int32
}

func environ() []string {
	return nil
}

var (
	argc int32
	argv **byte
)

//go:nosplit
func timediv(v int64, div int32, rem *int32) int32 {
	return 0
}

//go:nosplit
func argv_index(argv **byte, i int32) *byte {
	return nil
}
func goargs() {

}

func args(c int32, v **byte) {

}

func check() {

}
func parsedebugvars() {

}
