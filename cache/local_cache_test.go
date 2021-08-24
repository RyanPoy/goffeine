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

	cache.putToWindowQueue(node.NewWithWeight("key_1", 1, 10)) // 到这里，windowQ里面有内容了
	cache.putToWindowQueue(node.NewWithWeight("key_2", 2, 5))

	pNode, _ = cache.windowQ.Last()
	assert.Equal("key_2", pNode.Key)
}

func TestPutANodeInWindowQueueFirstWhenTheNodeDoesNotExistAndNodeWeightEqualsWindowQueueMaxWeight(t *testing.T) {
	assert := assert.New(t)
	cache := newLocalCache(100, 20, 60)

	pNode, ok := cache.windowQ.First()
	assert.Equal(false, ok)

	cache.putToWindowQueue(node.NewWithWeight("key_1", 1, 10)) // 到这里，windowQ里面有内容了
	cache.putToWindowQueue(node.NewWithWeight("key_2", 2, 20))

	pNode, _ = cache.windowQ.First()
	assert.Equal("key_2", pNode.Key)
}

func TestPutANodeInWindowQueueFirstWhenTheNodeDoesNotExistAndNodeWeightGreatWindowQueueMaxWeight(t *testing.T) {
	assert := assert.New(t)
	cache := newLocalCache(100, 20, 60)

	pNode, ok := cache.windowQ.First()
	assert.Equal(false, ok)

	cache.putToWindowQueue(node.NewWithWeight("key_1", 1, 10)) // 到这里，windowQ里面有内容了
	cache.putToWindowQueue(node.NewWithWeight("key_2", 2, 20)) //

	pNode, _ = cache.windowQ.First()
	assert.Equal("key_2", pNode.Key)
}

func TestEvictFromWindowShouldDoNothingWhenWindowQueueWeightLessThanMaxWeight(t *testing.T) {
	assert := assert.New(t)
	cache := newLocalCache(100, 20, 60)
	cache.putToWindowQueue(node.NewWithWeight("key_1", 1, 10)) //
	cache.putToWindowQueue(node.NewWithWeight("key_2", 2, 10)) // 到这里，windowQ的Weight超过了MaxWeight了

	_, ok := cache.probationQ.First()
	assert.Equal(false, ok)

	cache.evictFromWindow() // 不会做任何事情
	_, ok = cache.probationQ.First()
	assert.Equal(false, ok)
}

func TestEvictFromWindow(t *testing.T) {
	assert := assert.New(t)
	cache := newLocalCache(100, 20, 60)
	cache.putToWindowQueue(node.NewWithWeight("key_1", 1, 10)) //
	cache.putToWindowQueue(node.NewWithWeight("key_2", 2, 80)) // 到这里，windowQ的Weight超过了MaxWeight了
	cache.putToWindowQueue(node.NewWithWeight("key_3", 3, 10)) //
	cache.putToWindowQueue(node.NewWithWeight("key_4", 4, 15)) //

	pNode, ok := cache.probationQ.First()
	assert.Equal(false, ok)

	cache.evictFromWindow() // 这里应该依次淘汰把key_2，key_1, key_3 淘汰

	assert.Equal(15, cache.windowQ.Weight())

	pNode, ok = cache.probationQ.UnlinkFirst()
	assert.Equal("key_2", pNode.Key)

	pNode, ok = cache.probationQ.UnlinkFirst()
	assert.Equal("key_1", pNode.Key)

	pNode, ok = cache.probationQ.UnlinkFirst()
	assert.Equal("key_3", pNode.Key)
}

func TestEvictFromProbation(t *testing.T) {
	assert := assert.New(t)
	cache := newLocalCache(100, 20, 60)
	cache.putToWindowQueue(node.NewWithWeight("key_1", 1, 10)) //
	cache.putToWindowQueue(node.NewWithWeight("key_2", 2, 80)) // 到这里，windowQ的Weight超过了MaxWeight了
	cache.putToWindowQueue(node.NewWithWeight("key_3", 3, 10)) //
	cache.putToWindowQueue(node.NewWithWeight("key_4", 4, 15)) //

	cache.evictFromWindow() // 这里应该依次淘汰key_2，key_1, key_3 到 probation
	cache.evictFromProbation() // 这里应该淘汰key_3, key_1
	pNode, ok := cache.probationQ.UnlinkFirst()
	assert.Equal("key_2", pNode.Key)

	pNode, ok = cache.probationQ.UnlinkFirst()
	assert.Equal(false, ok)
}

func TestEvictFromProbationShouldDoNothingWhenWindowQueueWeightLessThanMaxWeight(t *testing.T) {
	assert := assert.New(t)
	cache := newLocalCache(100, 20, 60)
	cache.putToWindowQueue(node.NewWithWeight("key_1", 1, 10)) //
	cache.putToWindowQueue(node.NewWithWeight("key_2", 2, 30)) // 到这里，windowQ的Weight超过了MaxWeight了
	cache.putToWindowQueue(node.NewWithWeight("key_3", 3, 10)) //
	cache.putToWindowQueue(node.NewWithWeight("key_4", 4, 15)) //

	cache.evictFromWindow() // 这里应该依次淘汰key_2，key_1, key_3 到 probation
	cache.evictFromProbation() // 这里不会淘汰任何node
	assert.Equal(50, cache.probationQ.Weight())
}
