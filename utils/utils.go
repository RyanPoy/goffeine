package utils

import (
	"hash/fnv"
)

func HashCode(s string) uint32 {
	algorithm := fnv.New32()
	algorithm.Write([]byte(s))
	return algorithm.Sum32()
}

func CeilingPowerOfTwo(x int) int {
	tmp := -NumberOfLeadingZeros(x - 1)
	if tmp < 0 {
		tmp = 32 + tmp
	}
	return 1 << tmp
}

func NumberOfLeadingZeros(i int) int {
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
