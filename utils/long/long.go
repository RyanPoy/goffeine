package long

import (
	"math/bits"
)

// Returns the smallest power of two greater than or equal to {@code x}.
func CeilingPowerOfTwo(x int64) int64 {
	// From Hacker's Delight, Chapter 3, Harry S. Warren Jr.
	var n int64 = 1
	tmp := -1 * bits.LeadingZeros64(uint64(x-1))
	if tmp < 0 {
		tmp = 64 + tmp
	}
	return n << tmp
}
