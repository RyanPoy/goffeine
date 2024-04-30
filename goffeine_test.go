package goffeine_test

import (
	"github.com/stretchr/testify/assert"
	"goffeine"
	"testing"
	"time"
)

func TestBuilder(t *testing.T) {
	cache := goffeine.NewBuilder().MaximumSize(10_1000).ExpireAfterWrite(time.Minute, 5).RefreshAfterWrite(time.Hour, 1).Build()
	assert.Equal(t, 10_1000, cache.MaximumSize)

	assert.Equal(t, 5, cache.ExpireTime.Delay)
	assert.Equal(t, time.Minute, cache.ExpireTime.Duration)

	assert.Equal(t, 1, cache.RefreshTime.Delay)
	assert.Equal(t, time.Hour, cache.RefreshTime.Duration)
}
