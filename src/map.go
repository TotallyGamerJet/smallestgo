package runtime

type hmap struct {
	buckets    uintptr // array of 2^B Buckets. may be nil if count==0.
	oldbuckets uintptr // previous bucket array of half the size, non-nil only when growing
}
type bmap struct {
}
