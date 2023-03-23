// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Runtime type representation.

package runtime

type _type struct {
}

type typeOff int32

// these must stick around bc the linker will complain:
//dwarf: missing type: type:runtime.imethod

type imethod struct {
}

type interfacetype struct {
}

type maptype struct {
}

type arraytype struct {
}

type chantype struct {
}

type slicetype struct {
}

type functype struct {
}

type ptrtype struct {
}

type structfield struct {
}

type structtype struct {
}
