package utils

import (
	"hash/fnv"
)

func HashCode(s string) int32 {
	algorithm := fnv.New32()
	algorithm.Write([]byte(s))
	return int32(algorithm.Sum32())
}

