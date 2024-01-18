package queue

import (
	"github.com/stretchr/testify/assert"
	"goffeine/internal/node"
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
	q.LinkLast(node.New("id_123", 123))
	assert.Equal(1, q.Weight())
}

func TestAddTwice(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue()
	q.LinkLast(node.New("id_123", 123))
	q.LinkLast(node.New("id_123", 123))

	assert.Equal(1, q.Weight())
}

func TestAddMany(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue()
	q.LinkLast(node.New("id_123", 123))
	q.LinkLast(node.New("id_123", 123))
	q.LinkLast(node.New("id_456", 456))
	q.LinkLast(node.New("id_789", 789))
	assert.Equal(3, q.Weight())
}

func TestGetWhenNotExist(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue()
	v, ok := q.UnlinkFirst()
	assert.Equal(true, v == nil)
	assert.Equal(false, ok)
}

func TestGet(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue()
	pNode := node.New("id_123", 123)
	q.LinkLast(pNode)
	v, ok := q.UnlinkFirst()
	assert.Equal(pNode, v)
	assert.Equal(true, ok)
}

func TestRemove(t *testing.T) {
	assert := assert.New(t)
	q := newAccessOrderQueue()
	pNode1 := node.New("id_123", 123)
	q.LinkLast(pNode1)
	pNode2 := node.New("id_456", 456)
	q.LinkLast(pNode2)
	q.Remove(pNode1)
	assert.Equal(false, q.Contains(pNode1))
	assert.Equal(true, q.queue.Front().Value == pNode2)
}
