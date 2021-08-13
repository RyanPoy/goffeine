package internal

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitial(t *testing.T) {
	assert := assert.New(t)
	m := NewLRU(10)
	assert.Equal(0, m.Len())
	assert.Equal(10, m.Capacity())
}

func TestAddOnce(t *testing.T) {
	assert := assert.New(t)
	m := NewLRU(10)
	m.Add("id_123", 123)
	assert.Equal(1, m.Len())
	assert.Equal(10, m.Capacity())
}

func TestAddTwice(t *testing.T) {
	assert := assert.New(t)
	m := NewLRU(10)
	m.Add("id_123", 123)
	m.Add("id_123", 123)
	assert.Equal(1, m.Len())
	assert.Equal(10, m.Capacity())
}

func TestAddMany(t *testing.T) {
	assert := assert.New(t)
	m := NewLRU(10)
	m.Add("id_123", 123)
	m.Add("id_123", 123)
	m.Add("id_456", 456)
	m.Add("id_789", 789)
	assert.Equal(3, m.Len())
	assert.Equal(10, m.Capacity())
}

func TestGetWhenNotExist(t *testing.T) {
	assert := assert.New(t)
	m := NewLRU(10)
	_, err := m.Get("id_123")
	assert.Equal("『id_123』does not exist", err.Error())
}

func TestGet(t *testing.T) {
	assert := assert.New(t)
	m := NewLRU(10)
	m.Add("id_123", 123)
	value, _ := m.Get("id_123")
	assert.Equal(123, value)
}

func TestUpdateWhenAndAExistKeyButDifferentValue(t *testing.T) {
	assert := assert.New(t)
	m := NewLRU(10)
	m.Add("id_123", 123)
	m.Add("id_123", 456)
	value, _ := m.Get("id_123")
	assert.Equal(456, value)
}

func TestRemoveAndGetNilIfDoesNotExist(t *testing.T) {
	assert := assert.New(t)
	m := NewLRU(10)
	assert.Equal(nil, m.Remove("id_123"))
}

func TestRemove(t *testing.T) {
	assert := assert.New(t)
	m := NewLRU(10)
	m.Add("id_123", 123)
	assert.Equal(123, m.Remove("id_123"))

	_, err := m.Get("id_123")
	assert.Equal(fmt.Sprintf("『%s』does not exist", "id_123"), err.Error())
}


func TestFull(t *testing.T) {
	assert := assert.New(t)
	m := NewLRU(3)
	m.Add("id_123", 123)
	m.Add("id_456", 456)
	m.Add("id_789", 789)
	assert.Equal(true, m.IsFull())
}

func TestAddWillBeEliminatedAutomaticWhenCapacityIsFull(t *testing.T) {
	assert := assert.New(t)
	m := NewLRU(3)
	m.Add("id_123", 123)
	m.Add("id_456", 456)
	m.Add("id_789", 789)
	m.Add("id_abc", "abc")
	assert.Equal(true, m.IsFull())
	assert.Equal(3, m.Len())

	v, _ := m.Get("id_456")
	assert.Equal(456, v)

	v, _ = m.Get("id_789")
	assert.Equal(789, v)

	v, _ = m.Get("id_abc")
	assert.Equal("abc", v)

	_, v = m.Get("id_123")
	assert.Equal(nil, v)
}


