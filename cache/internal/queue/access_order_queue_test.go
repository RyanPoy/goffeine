package queue

import (
	"github.com/stretchr/testify/assert"
	"goffeine/cache/internal/node"
	"testing"
)

func newAccessOrderQueue() *AccessOrderQueue {
	return New()
}

func TestInitial(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue()
	assert.Equal(0, q.Weight())
}

func TestAddOnce(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue()
	q.Push(node.New("id_123", 123))
	assert.Equal(1, q.Weight())
}

func TestAddTwice(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue()
	q.Push(node.New("id_123", 123))
	q.Push(node.New("id_123", 123))

	assert.Equal(1, q.Weight())
}

func TestAddMany(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue()
	q.Push(node.New("id_123", 123))
	q.Push(node.New("id_123", 123))
	q.Push(node.New("id_456", 456))
	q.Push(node.New("id_789", 789))
	assert.Equal(3, q.Weight())
}

func TestGetWhenNotExist(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue()
	v, err := q.Pop()
	assert.Equal(true, v == nil)
	assert.Equal(EmptyError, err)
}

func TestGet(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue()
	pNode := node.New("id_123", 123)
	q.Push(pNode)
	v, err := q.Pop()
	assert.Equal(pNode, v)
	assert.Equal(nil, err)
}

func TestAccessOrderQueue_Remove(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue()
	pNode1 := node.New("id_123", 123)
	q.Push(pNode1)
	pNode2 := node.New("id_456", 456)
	q.Push(pNode2)
	q.Remove(pNode1)
	assert.Equal(false,q.Contains(pNode1))
	assert.Equal(true,q.queue.Front().Value==pNode2)
}

func TestAddWillBeEliminatedAutomaticWhenCapacityIsFull2(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue()

	pNode1 := node.New("id_123", 123)
	q.Push(pNode1)

	pNode2 := node.New("id_456", 456)
	q.Push(pNode2)

	pNode3 := node.New("id_789", 789)
	q.Push(pNode3)

	q.Push(pNode1)

	pNode4 := node.New("id_abc", "abc")
	q.Push(pNode4)

	assert.Equal(4, q.Weight())

	v, _ := q.Pop()
	assert.Equal(pNode2, v)

	v, _ = q.Pop()
	assert.Equal(pNode3, v)

	v, _ = q.Pop()
	assert.Equal(pNode1, v)

	v, _ = q.Pop()
	assert.Equal(pNode4, v)

	assert.Equal(true, q.IsEmpty())
}
