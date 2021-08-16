package queue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func newAccessOrderQueue(c int) *AccessOrderQueue {
	return New(c)
}

func TestInitial(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue(10)
	assert.Equal(0, q.Len())
	assert.Equal(10, q.Capacity())
}

func TestAddOnce(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue(10)
	q.Push("id_123")
	assert.Equal(1, q.Len())
	assert.Equal(10, q.Capacity())
}

func TestAddTwice(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue(10)
	q.Push("id_123")
	q.Push("id_123")

	assert.Equal(1, q.Len())
	assert.Equal(10, q.Capacity())
}

func TestAddMany(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue(10)
	q.Push("id_123")
	q.Push("id_123")
	q.Push("id_456")
	q.Push("id_789")
	assert.Equal(3, q.Len())
	assert.Equal(10, q.Capacity())
}

func TestGetWhenNotExist(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue(10)
	v := q.Pop()
	assert.Equal(nil, v)
}

func TestGet(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue(10)
	q.Push("id_123")
	v := q.Pop()
	assert.Equal("id_123", v)
}

func TestFull(t *testing.T) {
	assert := assert.New(t)
	q := New(3)
	q.Push("id_123")
	q.Push("id_456")
	q.Push("id_789")
	assert.Equal(true, q.IsFull())
}

func TestAddWillBeEliminatedAutomaticWhenCapacityIsFull(t *testing.T) {
	assert := assert.New(t)
	q := New(3)

	q.Push("id_123")
	q.Push("id_456")
	q.Push("id_789")
	q.Push("id_abc")
	assert.Equal(true, q.IsFull())
	assert.Equal(3, q.Len())

	assert.Equal("id_456",q.Pop())
	assert.Equal("id_789", q.Pop())
	assert.Equal("id_abc", q.Pop())
	assert.Equal(nil, q.Pop())
}

func TestAddWillBeEliminatedAutomaticWhenCapacityIsFull2(t *testing.T) {
	assert := assert.New(t)
	q := New(3)

	q.Push("id_123")
	q.Push("id_456")
	q.Push("id_789")
	q.Push("id_123")
	q.Push("id_abc")
	assert.Equal(true, q.IsFull())
	assert.Equal(3, q.Len())

	assert.Equal("id_789", q.Pop())
	assert.Equal("id_123", q.Pop())
	assert.Equal("id_abc", q.Pop())
	assert.Equal(true, q.IsEmpty())
}
