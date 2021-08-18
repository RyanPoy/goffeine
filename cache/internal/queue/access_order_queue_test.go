package queue

import (
	"github.com/stretchr/testify/assert"
	"goffeine/cache/internal/node"
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
	q.Push(node.New("id_123", 123))
	assert.Equal(1, q.Len())
	assert.Equal(10, q.Capacity())
}

func TestAddTwice(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue(10)
	q.Push(node.New("id_123", 123))
	q.Push(node.New("id_123", 123))

	assert.Equal(1, q.Len())
	assert.Equal(10, q.Capacity())
}

func TestAddMany(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue(10)
	q.Push(node.New("id_123", 123))
	q.Push(node.New("id_123", 123))
	q.Push(node.New("id_456", 456))
	q.Push(node.New("id_789", 789))
	assert.Equal(3, q.Len())
	assert.Equal(10, q.Capacity())
}

func TestGetWhenNotExist(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue(10)
	v, err := q.Pop()
	assert.Equal(true, v == nil)
	assert.Equal(EmptyError, err)
}

func TestGet(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue(10)
	pNode := node.New("id_123", 123)
	q.Push(pNode)
	v, err := q.Pop()
	assert.Equal(pNode, v)
	assert.Equal(nil, err)
}

func TestFull(t *testing.T) {
	assert := assert.New(t)
	q := New(3)
	q.Push(node.New("id_123", 123))
	q.Push(node.New("id_456", 456))
	q.Push(node.New("id_789", 789))
	assert.Equal(true, q.IsFull())
}

func TestAddWillBeEliminatedAutomaticWhenCapacityIsFull(t *testing.T) {
	assert := assert.New(t)
	q := New(3)

	pNode1 := node.New("id_123", 123)
	q.Push(pNode1)

	pNode2 := node.New("id_456", 456)
	q.Push(pNode2)

	pNode3 := node.New("id_789", 789)
	q.Push(pNode3)

	pNode4 := node.New("id_abc", "abc")
	q.Push(pNode4)

	assert.Equal(true, q.IsFull())
	assert.Equal(3, q.Len())

	v, err := q.Pop()
	assert.Equal(pNode2, v)

	v, _ = q.Pop()
	assert.Equal(pNode3, v)

	v, _ = q.Pop()
	assert.Equal(pNode4, v)

	v, err = q.Pop()
	assert.Equal(true, v == nil)
	assert.Equal(EmptyError, err)
}

func TestAddWillBeEliminatedAutomaticWhenCapacityIsFull2(t *testing.T) {
	assert := assert.New(t)
	q := New(3)

	pNode1 := node.New("id_123", 123)
	q.Push(pNode1)

	pNode2 := node.New("id_456", 456)
	q.Push(pNode2)

	pNode3 := node.New("id_789", 789)
	q.Push(pNode3)

	q.Push(pNode1)

	pNode4 := node.New("id_abc", "abc")
	q.Push(pNode4)

	assert.Equal(true, q.IsFull())
	assert.Equal(3, q.Len())

	v, _ := q.Pop()
	assert.Equal(pNode3, v)

	v, _ = q.Pop()
	assert.Equal(pNode1, v)

	v, _ = q.Pop()
	assert.Equal(pNode4, v)

	assert.Equal(true, q.IsEmpty())
}
