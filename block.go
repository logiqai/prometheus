package tsdb

import "sort"

const (
	magicIndex  = 0xCAFECAFE
	magicSeries = 0xAFFEAFFE
)

// Block handles reads against a block of time series data within a time window.
type Block interface {
	Querier(mint, maxt int64) Querier
}

const (
	flagNone = 0
	flagStd  = 1
)

// A skiplist maps offsets to values. The values found in the data at an
// offset are strictly greater than the indexed value.
type skiplist interface {
	// offset returns the offset to data containing values of x and lower.
	offset(x int64) (uint32, bool)
}

// simpleSkiplist is a slice of plain value/offset pairs.
type simpleSkiplist []skiplistPair

type skiplistPair struct {
	value  int64
	offset uint32
}

func (sl simpleSkiplist) offset(x int64) (uint32, bool) {
	// Search for the first offset that contains data greater than x.
	i := sort.Search(len(sl), func(i int) bool { return sl[i].value >= x })

	// If no element was found return false. If the first element is found,
	// there's no previous offset actually containing values that are x or lower.
	if i == len(sl) || i == 0 {
		return 0, false
	}
	return sl[i-1].offset, true
}
