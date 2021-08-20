package cache

import (
	"goffeine/cache/internal/node"
	"goffeine/cache/internal/queue"
	"goffeine/cache/internal/sketch"
	"math/rand"
	"time"
)

type Cache struct {
	cap        int
	fsketch    *sketch.FrequencySketch
	windowQ    *queue.AccessOrderQueue
	probationQ *queue.AccessOrderQueue
	protectedQ *queue.AccessOrderQueue
}

func New(cap int) Cache {
	return Cache{
		cap:        cap,
		fsketch:    sketch.New(cap),
		windowQ:    queue.New(cap),
		probationQ: queue.New(cap),
		protectedQ: queue.New(cap),
	}
}

func (c *Cache) Capacity() int {
	return c.cap
}

func (c *Cache) Len() int {
	return c.windowQ.Len() + c.probationQ.Len() + c.protectedQ.Len()
}

func (c *Cache) Contains(pNode *node.Node) bool {
	return c.windowQ.Contains(pNode) || c.probationQ.Contains(pNode) || c.protectedQ.Contains(pNode)
}

func (c *Cache) random() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Int() >= 50
}

func (c *Cache) addNodeWhenDoesNotExist(pNode *node.Node) {
	c.fsketch.Increment(pNode)
	if !c.windowQ.IsFull() { // 没满，则添加
		c.windowQ.Push(pNode)
		return
	}

	pNodeC, err := c.windowQ.Pop()
	c.windowQ.Push(pNode)
	if err != nil {
		return
	}
	c.addNodeToProbation(pNodeC)
}

func (c *Cache) addNodeToProbation(pNodeC *node.Node) {
	freqOfNodeC := c.fsketch.Frequency(pNodeC)
	if freqOfNodeC < 5 { // 被淘汰
		return
	}
	if !c.probationQ.IsFull() { // protationQ没满直接添加，@todo 这个需要确定一下
		c.probationQ.Push(pNodeC)
		return
	}
	pNodeV, err := c.probationQ.First()
	if err != nil {
		return
	}
	freqOfNodeV := c.fsketch.Frequency(pNodeV)
	if freqOfNodeC >= 5 && freqOfNodeC < freqOfNodeV {
		if !c.random() { // 随机淘汰c
			return
		}
		c.probationQ.RemoveFirst() // 把第1个淘汰掉，本质上就是把nodeV淘汰掉
		c.probationQ.Push(pNodeC)
	} else { // freqOfNodeC > freqOfNodeV
		c.probationQ.RemoveFirst() // 把第1个淘汰掉，本质上就是把nodeV淘汰掉
		c.probationQ.Push(pNodeC)
	}
}

// 往cache里面添加内容
func (c *Cache) Add(key string, value interface{}) {
	pNode := node.New(key, value)

	// 如果不在cache里面，先添加到admission
	if !c.Contains(pNode) {
		c.addNodeWhenDoesNotExist(pNode)
		return
	}

	// 如果在window或这protected存在，不处理
	if c.windowQ.Contains(pNode) {
		c.fsketch.Increment(pNode)
		c.windowQ.MoveToBack(pNode)
		return
	}
	if c.protectedQ.Contains(pNode) {
		c.fsketch.Increment(pNode)
		c.protectedQ.MoveToBack(pNode)
		return
	}
	// 如果在probation存在，需要移动到protected
	if c.probationQ.Contains(pNode) {
		element := c.probationQ.GetQueueElementBy(pNode)
		pNode := element.Value.(*node.Node)
		if !c.protectedQ.IsFull() {
			c.protectedQ.Push(pNode)
			c.fsketch.Increment(pNode)
			return
		}
		if pNodeC, err := c.protectedQ.Pop(); err == nil {
			c.protectedQ.Push(pNode)
			c.addNodeToProbation(pNodeC)
			c.fsketch.Increment(pNode)
		}
		c.fsketch.Increment(pNode)
	}
}
