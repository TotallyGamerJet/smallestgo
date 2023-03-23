// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"unsafe"
)

func unsafestring(ptr unsafe.Pointer, len int) {
}

// Keep this code in sync with cmd/compile/internal/walk/builtin.go:walkUnsafeString
func unsafestring64(ptr unsafe.Pointer, len64 int64) {
}

func unsafestringcheckptr(ptr unsafe.Pointer, len64 int64) {

}

func panicunsafestringlen() {

}

func panicunsafestringnilptr() {

}

// Keep this code in sync with cmd/compile/internal/walk/builtin.go:walkUnsafeSlice
func unsafeslice(et *_type, ptr unsafe.Pointer, len int) {

}

// Keep this code in sync with cmd/compile/internal/walk/builtin.go:walkUnsafeSlice
func unsafeslice64(et *_type, ptr unsafe.Pointer, len64 int64) {

}

func unsafeslicecheckptr(et *_type, ptr unsafe.Pointer, len64 int64) {

}

func panicunsafeslicelen() {

}

func panicunsafeslicenilptr() {

}
