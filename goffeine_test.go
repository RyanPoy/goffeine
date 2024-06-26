package goffeine_test

import (
	"github.com/stretchr/testify/assert"
	"goffeine"
	"testing"
	"time"
)

type CacheItem struct {
	Foo int
	Bar string
}

func NewCache() *goffeine.Goffeine {
	return goffeine.NewBuilder().MaximumSize(10_1000).ExpireAfterWrite(time.Minute, 5).RefreshAfterWrite(time.Hour, 1).Build()
}
func NewCacheWithMaximumSize(maxSize int) *goffeine.Goffeine {
	return goffeine.NewBuilder().MaximumSize(maxSize).ExpireAfterWrite(time.Minute, 5).RefreshAfterWrite(time.Hour, 1).Build()
}

func TestBuilder(t *testing.T) {
	cache := NewCache()
	assert.Equal(t, 10_1000, cache.MaximumSize())
	assert.Equal(t, (5 * time.Minute).Milliseconds(), cache.ExpireMilliseconds())
	assert.Equal(t, (1 * time.Hour).Milliseconds(), cache.RefreshMilliseconds())
}

func TestGoffeineSize(t *testing.T) {
	cache := NewCacheWithMaximumSize(100)
	assert.Equal(t, 100, cache.MaximumSize())
	assert.Equal(t, 1, cache.WindowMaximumSize())
	assert.Equal(t, 19, cache.ProbationMaximumSize())
	assert.Equal(t, 80, cache.ProtectedMaximumSize())

	cache = NewCacheWithMaximumSize(0)
	assert.Equal(t, 1, cache.WindowMaximumSize())
	assert.Equal(t, 1, cache.ProbationMaximumSize())
	assert.Equal(t, 1, cache.ProtectedMaximumSize())
}

func TestBasicPutAndGet(t *testing.T) {
	cache := NewCache()
	cache.Put("a", 1)
	cache.Put("b", CacheItem{2, "b"})

	v, _ := cache.Get("a")
	assert.Equal(t, 1, v)

	v, _ = cache.Get("b")
	assert.Equal(t, 2, v.(CacheItem).Foo)
	assert.Equal(t, "b", v.(CacheItem).Bar)

	_, ok := cache.Get("c")
	assert.False(t, ok)
}

func TestPutToAFullWindowAndGet(t *testing.T) {
	// window maximum Size is 3,
	// so we must put 3/0.01 maximum Size to cache
	// 3 / 0.01 = 300
	cache := NewCacheWithMaximumSize(300)
	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)
	cache.Put("d", CacheItem{4, "d"})

	v, ok := cache.Get("a")
	assert.False(t, ok)

	v, _ = cache.Get("b")
	assert.Equal(t, 2, v)

	v, _ = cache.Get("c")
	assert.Equal(t, 3, v)

	v, _ = cache.Get("d")
	assert.Equal(t, 4, v.(CacheItem).Foo)
	assert.Equal(t, "d", v.(CacheItem).Bar)
}
