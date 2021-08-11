package long

import "goffeine/utils/integer"

// Returns the smallest power of two greater than or equal to {@code x}.
func CeilingPowerOfTwo(x int64) int64 {
    // From Hacker's Delight, Chapter 3, Harry S. Warren Jr.
	var n int64 = 1
	tmp := -NumberOfLeadingZeros(x - 1)
	if tmp < 0 {
		tmp = 64 + tmp
	}
    return n << tmp
}

func NumberOfLeadingZeros(i int64) int {
	x := int(uint64(i) >> 32)
	if x == 0 {
		return 32 + integer.NumberOfLeadingZeros(int(i))
	}
	return integer.NumberOfLeadingZeros(x)
}
