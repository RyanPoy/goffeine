package integer

import "math/bits"

// Returns the smallest power of two greater than or equal to {@code x}.
func CeilingPowerOfTwo(x int) int {
    // From Hacker's Delight, Chapter 3, Harry S. Warren Jr.
	n := 1
	tmp := -NumberOfLeadingZeros(x - 1)
	if tmp < 0 {
		tmp = 32 + tmp
	}
    return n << tmp
}

func NumberOfLeadingZeros(i int) int {
	return bits.LeadingZeros32(uint32(i))
}
