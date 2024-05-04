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

	assert.Equal(t, 5, cache.ExpireTime().Delay)
	assert.Equal(t, time.Minute, cache.ExpireTime().Duration)

	assert.Equal(t, 1, cache.RefreshTime().Delay)
	assert.Equal(t, time.Hour, cache.RefreshTime().Duration)
}

func TestGoffeineSize(t *testing.T) {
	cache := NewCacheWithMaximumSize(100)
	assert.Equal(t, 100, cache.MaximumSize())
	assert.Equal(t, 1, cache.WindowMaximumSize())
	assert.Equal(t, 20, cache.ProbationMaximumSize())
	assert.Equal(t, 80, cache.ProtectedMaximumSize())

	cache = NewCacheWithMaximumSize(0)
	assert.Equal(t, 1, cache.WindowMaximumSize())
	assert.Equal(t, 1, cache.ProbationMaximumSize())
	assert.Equal(t, 1, cache.ProtectedMaximumSize())
}

func TestSetAndGet(t *testing.T) {
	cache := NewCache()
	cache.Set("a", 1)
	cache.Set("b", CacheItem{2, "b"})

	v, _ := cache.Get("a")
	assert.Equal(t, 1, v)

	v, _ = cache.Get("b")
	assert.Equal(t, 2, v.(*CacheItem).Foo)
	assert.Equal(t, "b", v.(*CacheItem).Bar)

	_, ok := cache.Get("c")
	assert.False(t, ok)
}
