package runtime

import (
	"unsafe"
)

type hex int

func gwrite(b []byte) {

}
func printlock() {
}

func printunlock() {
}

func hexdumpWords(p, end uintptr, mark func(uintptr) byte) {

}
func printsp() {
}

func printnl() {
}

func printbool(v bool) {

}

func printfloat(v float64) {
}

func printcomplex(c complex128) {
}

func printuint(v uint64) {
}

func printint(v int64) {
}

func printhex(v uint64) {
}

func printpointer(p unsafe.Pointer) {
}
func printuintptr(p uintptr) {
}

func printstring(s string) {
}

func printslice(s []byte) {
}

func printeface(e eface) {
}

func printiface(i iface) {
}
