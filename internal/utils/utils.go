package utils

import (
	"hash/fnv"
	"math"
	"math/bits"
)

type CmpType interface {
	uint | uint8 | uint16 | uint32 | uint64 |
		int | int8 | int16 | int32 | int64 |
		float32 | float64
}

func Max[T1, T2 CmpType](a T1, b T2) float64 {
	return math.Max(float64(a), float64(b))
}

func Min[T1, T2 CmpType](a T1, b T2) float64 {
	return math.Min(float64(a), float64(b))
}

func HashCode(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func CeilingPowerOfTwo32(x int) int {
	// From Hacker's Delight, Chapter 3, Harry S. Warren Jr.
	tmp := -1 * bits.LeadingZeros32(uint32(x-1))
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
