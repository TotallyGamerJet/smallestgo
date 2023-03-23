package runtime

//go:nosplit
func debugCallCheck(pc uintptr) string { return "" }

//go:nosplit
func debugCallWrap(dispatch uintptr) {}
