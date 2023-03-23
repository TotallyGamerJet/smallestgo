package runtime

const (
	// The number of levels in the radix tree.
	summaryLevels = 5

	// Constants for testing.
	pageAlloc32Bit = 0
	pageAlloc64Bit = 1

	// Number of bits needed to represent all indices into the L1 of the
	// chunks map.
	//
	// See (*pageAlloc).chunks for more details. Update the documentation
	// there should this number change.
	pallocChunksL1Bits = 13
)
