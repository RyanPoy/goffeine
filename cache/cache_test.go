package cache

func newCache(n int) Cache {
	return New(n)
}

//func TestInitial(t *testing.T) {
//	assert := assert.New(t)
//	cache := newCache(10)
//	assert.Equal(10, cache.Capacity())
//	assert.Equal(0, cache.Len())
//	assert.Equal(false, cache.Contains(node.New("abc", 1)))
//}
