package runtime

type pageTraceBuf struct {
}

func finishPageTrace() {
}

//go:systemstack
func pageTraceAlloc(pp *p, now int64, base, npages uintptr) {
}

//go:systemstack
func pageTraceFree(pp *p, now int64, base, npages uintptr) {
}

func initPageTrace(env string) {
}

//go:systemstack
func pageTraceScav(pp *p, now int64, base, npages uintptr) {
}
