package runtime

import (
	"runtime/internal/sys"
)

// Central list of free objects of a given size.
type mcentral struct {
	_         sys.NotInHeap
	spanclass spanClass

	// partial and full contain two mspan sets: one of swept in-use
	// spans, and one of unswept in-use spans. These two trade
	// roles on each GC cycle. The unswept set is drained either by
	// allocation or by the background sweeper in every GC cycle,
	// so only two roles are necessary.
	//
	// sweepgen is increased by 2 on each GC cycle, so the swept
	// spans are in partial[sweepgen/2%2] and the unswept spans are in
	// partial[1-sweepgen/2%2]. Sweeping pops spans from the
	// unswept set and pushes spans that are still in-use on the
	// swept set. Likewise, allocating an in-use span pushes it
	// on the swept set.
	//
	// Some parts of the sweeper can sweep arbitrary spans, and hence
	// can't remove them from the unswept set, but will add the span
	// to the appropriate swept list. As a result, the parts of the
	// sweeper and mcentral that do consume from the unswept list may
	// encounter swept spans, and these should be ignored.
	partial [2]spanSet // list of spans with a free object
	full    [2]spanSet // list of spans with no free objects
}

func (c *mcentral) uncacheSpan(s *mspan)               {}
func (c *mcentral) cacheSpan() *mspan                  { return nil }
func (c *mcentral) fullSwept(sweepgen uint32) *spanSet { return nil }
func (c *mcentral) init(spc spanClass)                 {}
