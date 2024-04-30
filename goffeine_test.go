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

func TestBuilder(t *testing.T) {
	cache := goffeine.NewBuilder().MaximumSize(10_1000).ExpireAfterWrite(time.Minute, 5).RefreshAfterWrite(time.Hour, 1).Build()
	assert.Equal(t, 10_1000, cache.MaximumSize())

	assert.Equal(t, 5, cache.ExpireTime().Delay)
	assert.Equal(t, time.Minute, cache.ExpireTime().Duration)

	assert.Equal(t, 1, cache.RefreshTime().Delay)
	assert.Equal(t, time.Hour, cache.RefreshTime().Duration)
}

func NewCache() *goffeine.Goffeine {
	return goffeine.NewBuilder().MaximumSize(10_1000).ExpireAfterWrite(time.Minute, 5).RefreshAfterWrite(time.Hour, 1).Build()
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
