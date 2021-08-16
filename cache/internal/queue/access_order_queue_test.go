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
	v, err := q.Pop()
	assert.Equal(nil, v)
	assert.Equal(EmptyError, err)
}

func TestGet(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue(10)
	q.Push("id_123")
	v, err := q.Pop()
	assert.Equal("id_123", v)
	assert.Equal(nil, err)
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

	v, err := q.Pop()
	assert.Equal("id_456", v)

	v, _ = q.Pop()
	assert.Equal("id_789", v)

	v, _ = q.Pop()
	assert.Equal("id_abc", v)

	v, err = q.Pop()
	assert.Equal(nil, v)
	assert.Equal(EmptyError, err)
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

	v, _ := q.Pop()
	assert.Equal("id_789", v)

	v, _ = q.Pop()
	assert.Equal("id_123", v)

	v, _ = q.Pop()
	assert.Equal("id_abc", v)

	assert.Equal(true, q.IsEmpty())
}
