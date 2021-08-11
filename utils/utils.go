package utils

import (
	"hash/fnv"
	"math/bits"
)

func HashCode(s string) int32 {
	algorithm := fnv.New32()
	algorithm.Write([]byte(s))
	return int32(algorithm.Sum32())
}

func CeilingPowerOfTwo32(x int) int {
    // From Hacker's Delight, Chapter 3, Harry S. Warren Jr.
	tmp := -1 * bits.LeadingZeros32(uint32(x - 1))
	if tmp < 0 {
		tmp = 32 + tmp
	}
    return 1 << tmp
}

func CeilingPowerOfTwo64(x int64) int64 {
	// From Hacker's Delight, Chapter 3, Harry S. Warren Jr.
	var n int64 = 1
	tmp := -1 * bits.LeadingZeros64(uint64(x-1))
	if tmp < 0 {
		tmp = 64 + tmp
	}
	return n << tmp
}
