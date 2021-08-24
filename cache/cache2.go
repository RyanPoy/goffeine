package cache

import (
	"goffeine/cache/internal/node"
	"goffeine/cache/internal/queue"
	"goffeine/cache/internal/sketch"
	"sync"
)

type LocalCache struct {
	maxWeight   int
	sketch      *sketch.FrequencySketch
	windowQ     *queue.AccessOrderQueue
	probationQ  *queue.AccessOrderQueue
	protectedQ  *queue.AccessOrderQueue
	hashmap     sync.Map
	Weight      int //集合当前权重，容量
	//wMaxWeight  int //window大小
	//ptMaxWeight int // protectedQ size
}

func NewLocalCache(maxWeight, windowQuqueMaxWeight, protectedQueueMaxWeight int) LocalCache {
	return LocalCache{
		maxWeight:   maxWeight,
		sketch:      sketch.New(maxWeight),
		windowQ:     queue.NewWith(windowQuqueMaxWeight),
		probationQ:  queue.New(),
		protectedQ:  queue.NewWith(protectedQueueMaxWeight),
		hashmap:     sync.Map{},
		Weight:      0,                       //集合当前权重，容量
		//wMaxWeight:  windowQuqueMaxWeight,    //window大小
		//ptMaxWeight: protectedQueueMaxWeight, // protectedQ size
	}
}

func (c *LocalCache) Put(key string, value interface{}, weight int) {
	c.PutWithWeight(key, value, weight)
}

func (c *LocalCache) PutWithWeight(key string, value interface{}, weight int) {
	pNode := node.NewWithWeight(key, value, weight)
	c.put(pNode)
}

func (c *LocalCache) put(pNewNode *node.Node) {
	_, ok := c.hashmap.Load(pNewNode.Key)
	if !ok { // * 添加一个新node
		c.hashmap.Store(pNewNode.Key, pNewNode)
		c.putToWindowQueue(pNewNode)
		c.sketch.Increment(pNewNode)
		// 如果node的权重大于windowq的最大权重，push到windowq的first，否则push到windowq的last
		// 如果window的当前权重大于window最大权重，挪动window的first，放到probation的last，直到window的当前权重小于等于window的最大权重。到此：window的当前权重已经收缩到合理值了。
		// loop：如果cache的当前权重超出最大权重，进行淘汰：
		//   如果probation的 victim(first) 和 candidate(last) 进行对比，按照FrequencyCandidate 和 FrequencyVictim 和 随机数 一起来判断淘汰 Victim 或者 Candidate。到此：Cache的当前权重已经收缩到合理值了。
	}
}

// 把一个node添加到windowQ，策略如下：
// 1、如果windowQ不存在这个node：
//  1.1、如果node的权重大于windowq的最大权重，push到windowq的first，否则push到windowq的last
func (c *LocalCache) putToWindowQueue(pNode *node.Node) {
	if !c.windowQ.Contains(pNode) {
		if pNode.Weight >= c.windowQ.MaxWeight {
			c.windowQ.LinkFirst(pNode)
		} else {
			c.windowQ.LinkLast(pNode)
		}
	}
}


////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// * 添加一个新node
// hashmap.put(node)
// 如果node的权重大于windowq的最大权重，push到windowq的first，否则push到windowq的last
// 如果window的当前权重大于window最大权重，挪动window的first，放到probation的last，直到window的当前权重小于等于window的最大权重。到此：window的当前权重已经收缩到合理值了。
// loop：如果cache的当前权重超出最大权重，进行淘汰：
//   如果probation的 victim(first) 和 candidate(last) 进行对比，按照FrequencyCandidate 和 FrequencyVictim 和 随机数 一起来判断淘汰 Victim 或者 Candidate。到此：Cache的当前权重已经收缩到合理值了。

// * 更新一个node
// hashmap.put(node)
// node.update Weight and value
// if node is belongs to windowq:
//   如果node的权重大于windowq的最大权重，移动到windowq的first，否则移动到windowq的last
//   如果window的当前权重大于window最大权重，挪动window的first，放到probation的last，直到window的当前权重小于等于window的最大权重。到此：window的当前权重已经收缩到合理值了。
// elif node is belongs to probationq:
//   挪动node到protected
//   如果protected的当前权重大于protected最大权重，挪动protected的first，放到probation的last，直到protected的当前权重小于等于protected的最大权重。到此：protected的当前权重已经收缩到合理值了。
//   loop：如果cache的当前权重超出最大权重，进行淘汰：
//     如果probation的 victim(first) 和 candidate(last) 进行对比，按照FrequencyCandidate 和 FrequencyVictim 和 随机数 一起来判断淘汰 Victim 或者 Candidate。到此：Cache的当前权重已经收缩到合理值了。
// elif node is belongs to protected:
//   挪动node到protected的last
//   如果protected的当前权重大于protected最大权重，挪动protected的first，放到probation的last，直到protected的当前权重小于等于protected的最大权重。到此：protected的当前权重已经收缩到合理值了。
//   loop：如果cache的当前权重超出最大权重，进行淘汰：
//     如果probation的 victim(first) 和 candidate(last) 进行对比，按照FrequencyCandidate 和 FrequencyVictim 和 随机数 一起来判断淘汰 Victim 或者 Candidate。到此：Cache的当前权重已经收缩到合理值了。

// * 获取一个key的value
// if hit in window:
//   挪动到window队尾，返回value
// elif hit in probation:
//   挪动node，从probation到protected的队尾
//   如果protected的当前权重大于protected最大权重，挪动protected的first，放到probation的last，直到protected的当前权重小于等于protected的最大权重。到此：protected的当前权重已经收缩到合理值了。
//   loop：如果cache的当前权重超出最大权重，进行淘汰：
//     如果probation的 victim(first) 和 candidate(last) 进行对比，按照FrequencyCandidate 和 FrequencyVictim 和 随机数 一起来判断淘汰 Victim 或者 Candidate。到此：Cache的当前权重已经收缩到合理值了。
// elif hit in protected:
//   挪动node到protected的队尾

// *删除一个node
// 直接从hashmap里面删除掉
// 从对应q里面删除掉
