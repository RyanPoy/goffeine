package utils

import (
	"hash/fnv"
)

func HashCode(s string) uint32 {
	algorithm := fnv.New32()
	algorithm.Write([]byte(s))
	return algorithm.Sum32()
}

// Returns the smallest power of two greater than or equal to {@code x}.
func CeilingPowerOfTwoForLong(x int64) int64 {
    // From Hacker's Delight, Chapter 3, Harry S. Warren Jr.
	var n int64 = 1
	tmp := -NumberOfLeadingZerosForLong(x - 1)
	if tmp < 0 {
		tmp = 64 + tmp
	}
    return n << tmp
}

func NumberOfLeadingZerosForLong(i int64) int {
	x := int(uint64(i) >> 32)
	if x == 0 {
		return 32 + NumberOfLeadingZerosForInt(int(i))
	}
	return NumberOfLeadingZerosForInt(x)
}

func CeilingPowerOfTwoForInt(x int) int {
	tmp := -NumberOfLeadingZerosForInt(x - 1)
	if tmp < 0 {
		tmp = 32 + tmp
	}
	return 1 << tmp
}

func NumberOfLeadingZerosForInt(i int) int {
	// HD, Count leading 0's
	if i < 0 {
		return 0
	}
	if i == 0 {
		return 32
	}
	n := 31
	if i >= 1<<16 {
		n -= 16
		i = int(uint(i) >> 16)
	}
	if i >= 1<<8 {
		n -= 8
		i = int(uint(i) >> 8)
	}
	if i >= 1<<4 {
		n -= 4
		i = int(uint(i) >> 4)
	}
	if i >= 1<<2 {
		n -= 2
		i = int(uint(i) >> 2)
	}
	v := int(uint(i) >> 1)
	return n - v
}
