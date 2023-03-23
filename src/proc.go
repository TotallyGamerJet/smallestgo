// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import _ "unsafe"

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
type gList struct {
}

// An initTask represents the set of initializations that need to be done for a package.
// Keep in sync with ../../test/initempty.go:initTask
type initTask struct {
}
