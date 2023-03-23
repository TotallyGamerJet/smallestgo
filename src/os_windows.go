package runtime

import "unsafe"

type stdFunction unsafe.Pointer

func usleep2(dt int32)

type sigset struct{}

//go:nosplit
func exit(code int32) {}

//go:nosplit
func write1(fd uintptr, buf unsafe.Pointer, n int32) int32 { return 0 }

//go:nosplit
//go:cgo_unsafe_args
func stdcall1(fn stdFunction, a0 uintptr) uintptr { return 0 }

func usleep(us uint32) {}
func osinit()          {}
