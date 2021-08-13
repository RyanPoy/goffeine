package cache

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitial(t *testing.T) {
	assert := assert.New(t)
	lru := NewLRU(10)
	assert.Equal(0, lru.Len())
	assert.Equal(10, lru.Capacity())
}

func TestAddOnce(t *testing.T) {
	assert := assert.New(t)
	lru := NewLRU(10)
	lru.Add("id_123", 123)
	assert.Equal(1, lru.Len())
	assert.Equal(10, lru.Capacity())
}

func TestAddTwice(t *testing.T) {
	assert := assert.New(t)
	lru := NewLRU(10)
	lru.Add("id_123", 123)
	lru.Add("id_123", 123)
	assert.Equal(1, lru.Len())
	assert.Equal(10, lru.Capacity())
}

func TestAddMany(t *testing.T) {
	assert := assert.New(t)
	lru := NewLRU(10)
	lru.Add("id_123", 123)
	lru.Add("id_123", 123)
	lru.Add("id_456", 456)
	lru.Add("id_789", 789)
	assert.Equal(3, lru.Len())
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

func TestUpdateWhenAndAExistKeyButDifferentValue(t *testing.T) {
	assert := assert.New(t)
	lru := NewLRU(10)
	lru.Add("id_123", 123)
	lru.Add("id_123", 456)
	value, _ := lru.Get("id_123")
	assert.Equal(456, value)
}

func TestRemoveAndGetNilIfDoesNotExist(t *testing.T) {
	assert := assert.New(t)
	lru := NewLRU(10)
	assert.Equal(nil, lru.Remove("id_123"))
}

func TestRemove(t *testing.T) {
	assert := assert.New(t)
	lru := NewLRU(10)
	lru.Add("id_123", 123)
	assert.Equal(123, lru.Remove("id_123"))

	_, err := lru.Get("id_123")
	assert.Equal(fmt.Sprintf("『%s』does not exist", "id_123"), err.Error())
}


func TestFull(t *testing.T) {
	assert := assert.New(t)
	lru := NewLRU(3)
	lru.Add("id_123", 123)
	lru.Add("id_456", 456)
	lru.Add("id_789", 789)
	assert.Equal(true, lru.IsFull())
}

func TestAddWillBeEliminatedAutomaticWhenCapacityIsFull(t *testing.T) {
	assert := assert.New(t)
	lru := NewLRU(3)

	lru.Add("id_123", 123)
	lru.Add("id_456", 456)
	lru.Add("id_789", 789)
	lru.Add("id_abc", "abc")
	assert.Equal(true, lru.IsFull())
	assert.Equal(3, lru.Len())

	v, _ := lru.Get("id_456")
	assert.Equal(456, v)

	v, _ = lru.Get("id_789")
	assert.Equal(789, v)

	v, _ = lru.Get("id_abc")
	assert.Equal("abc", v)

	_, v = lru.Get("id_123")
	assert.Equal(nil, v)
}


