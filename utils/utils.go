package utils

import (
	"hash/fnv"
)

func HashCode(s string) int {
	algorithm := fnv.New32()
	algorithm.Write([]byte(s))
	return int(algorithm.Sum32())
}

