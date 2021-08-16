package cache

import (
	"sync"
	"testing"
)

func newCache(n int) Cache {
	return New(n)
}

func TestInitial(t *testing.T) {
	//assert := assert.New(t)
	//cache := newCache(10)
	//assert.Equal(10, cache.Capacity())
	//assert.Equal(0, cache.Len())
	//assert.Equal(false, cache.Contains("abc"))

	m := sync.Map {}
	m.Store("A", 1)
}
