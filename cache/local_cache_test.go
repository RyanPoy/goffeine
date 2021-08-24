package cache

import (
	"github.com/stretchr/testify/assert"
	"goffeine/cache/internal/node"
	"testing"
)

func newLocalCache(maxWeight, windowQuqueMaxWeight, protectedQueueMaxWeight int) LocalCache {
	return NewLocalCache(maxWeight, windowQuqueMaxWeight, protectedQueueMaxWeight)
}

func TestPutANodeInWindowQueueLastWhenTheNodeDoesNotExistAndNodeWeightLessThanWindowQueueMaxWeight(t *testing.T) {
	assert := assert.New(t)
	cache := newLocalCache(100, 20, 60)

	pNode, ok := cache.windowQ.Last()
	assert.Equal(false, ok)

	cache.putToWindowQueue(node.NewWithWeight("key_1", 1, 10)) // 执行完后，windowQ里面有内容了
	cache.putToWindowQueue(node.NewWithWeight("key_2", 2, 5))

	pNode, _ = cache.windowQ.Last()
	assert.Equal("key_2", pNode.Key)
}

func TestPutANodeInWindowQueueFirstWhenTheNodeDoesNotExistAndNodeWeightEqualsWindowQueueMaxWeight(t *testing.T) {
	assert := assert.New(t)
	cache := newLocalCache(100, 20, 60)

	pNode, ok := cache.windowQ.First()
	assert.Equal(false, ok)

	cache.putToWindowQueue(node.NewWithWeight("key_1", 1, 10)) // 执行完后，windowQ里面有内容了
	cache.putToWindowQueue(node.NewWithWeight("key_2", 2, 20))

	pNode, _ = cache.windowQ.First()
	assert.Equal("key_2", pNode.Key)
}

func TestPutANodeInWindowQueueFirstWhenTheNodeDoesNotExistAndNodeWeightGreatWindowQueueMaxWeight(t *testing.T) {
	assert := assert.New(t)
	cache := newLocalCache(100, 20, 60)

	pNode, ok := cache.windowQ.First()
	assert.Equal(false, ok)

	cache.putToWindowQueue(node.NewWithWeight("key_1", 1, 10)) // 执行完后，windowQ里面有内容了
	cache.putToWindowQueue(node.NewWithWeight("key_2", 2, 20)) //

	pNode, _ = cache.windowQ.First()
	assert.Equal("key_2", pNode.Key)
}
