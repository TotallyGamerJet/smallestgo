package runtime

import "unsafe"

func gogetenv(key string) string { return "" }

var _cgo_setenv unsafe.Pointer   // pointer to C function
var _cgo_unsetenv unsafe.Pointer // pointer to C function

func setenv_c(k string, v string) {}
func unsetenv_c(k string)         {}
