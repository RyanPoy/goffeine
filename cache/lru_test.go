package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitial(t *testing.T) {
	assert := assert.New(t)
	lru := NewLRU(10)
	assert.Equal(0, lru.Len())
	assert.Equal(10, lru.Capacity())
}


func TestAdd(t *testing.T) {
	assert := assert.New(t)
	lru := NewLRU(10)
	lru.Add("id_123", 123)
	assert.Equal(1, lru.Len())
	assert.Equal(10, lru.Capacity())
}

func TestGetWhenNotExist(t *testing.T) {
	assert := assert.New(t)
	lru := NewLRU(10)
	_, err := lru.Get("id_123")
	assert.Equal("『id_123』does not exist", err.Error())
}

func TestGet(t *testing.T) {
	assert := assert.New(t)
	lru := NewLRU(10)
	lru.Add("id_123", 123)
	value, _ := lru.Get("id_123")
	assert.Equal(123, value)
}

