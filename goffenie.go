package goffeine

import (
	"container/list"
	"goffeine/internal/node"
	"sync"
	"time"
)

type TimeUnion struct {
	Duration time.Duration
	Delay    int
}

func NewBuilder() *GoffeineBuilder {
	return &GoffeineBuilder{}
}

// A GoffeineBuilder is used to create a Goffeine instance
// e.g.	goffeine.NewBuilder().maximumSize(10).ExpireAfterWrite(time.Second, 5).Build()
type GoffeineBuilder struct {
	maximumSize      int
	expireTimeUnion  TimeUnion
	refreshTimeUnion TimeUnion
}

func (b *GoffeineBuilder) MaximumSize(size int) *GoffeineBuilder {
	if size < 1 {
		size = 3 // window: 1, probation: 1, protected: 1
	}
	b.maximumSize = size
	return b
}

func (b *GoffeineBuilder) ExpireAfterWrite(duration time.Duration, delay int) *GoffeineBuilder {
	b.expireTimeUnion = TimeUnion{Duration: duration, Delay: delay}
	return b
}

func (b *GoffeineBuilder) RefreshAfterWrite(duration time.Duration, delay int) *GoffeineBuilder {
	b.refreshTimeUnion = TimeUnion{Duration: duration, Delay: delay}
	return b
}

func (b *GoffeineBuilder) Build() *Goffeine {
	windowMaxsize := b.maximumSize / 100
	if windowMaxsize < 1 {
		windowMaxsize = 1
	}

	probationMaxsize := (b.maximumSize - windowMaxsize) * 20 / 100
	if probationMaxsize < 1 {
		probationMaxsize = 1
	}

	protectedMaxsize := b.maximumSize - windowMaxsize - probationMaxsize
	if protectedMaxsize < 1 {
		protectedMaxsize = 1
	}

	return &Goffeine{
		maximumSize:          b.maximumSize,
		windowMaximumSize:    windowMaxsize,
		probationMaximumSize: probationMaxsize,
		protectedMaximumSize: protectedMaxsize,

		expireTime:  b.expireTimeUnion,
		refreshTime: b.refreshTimeUnion,
		data:        &sync.Map{},
	}
}

// A Goffeine represents a cache
// It is implemented with Window-TinyLFU algorithm
type Goffeine struct {
	data                 *sync.Map
	maximumSize          int
	window               list.List
	windowMaximumSize    int
	probation            list.List
	probationMaximumSize int
	protected            list.List
	protectedMaximumSize int
	counter              map[string]int
	expireTime           TimeUnion
	refreshTime          TimeUnion
}

func (g *Goffeine) MaximumSize() int          { return g.maximumSize }
func (g *Goffeine) ExpireTime() TimeUnion     { return g.expireTime }
func (g *Goffeine) RefreshTime() TimeUnion    { return g.refreshTime }
func (g *Goffeine) WindowMaximumSize() int    { return g.windowMaximumSize }
func (g *Goffeine) ProbationMaximumSize() int { return g.probationMaximumSize }
func (g *Goffeine) ProtectedMaximumSize() int { return g.protectedMaximumSize }

func (g *Goffeine) windowIsFull() bool    { return g.window.Len() == g.windowMaximumSize }
func (g *Goffeine) probationIsFull() bool { return g.probation.Len() == g.probationMaximumSize }
func (g *Goffeine) protectedIsFull() bool { return g.protected.Len() == g.protectedMaximumSize }

func (g *Goffeine) Get(key string) (any, bool) {
	ele, ok := g.data.Load(key)
	if !ok {
		return nil, false
	}
	gnode := ele.(*list.Element).Value.(*node.GoffeineNode)
	return gnode.Value, true
}

func (g *Goffeine) Set(key string, value any) {
	g.set(key, value, g.expireTime)
}

func (g *Goffeine) SetWithDelay(key string, value any, delay int) {
	g.set(key, value, TimeUnion{Delay: delay, Duration: g.expireTime.Duration})
}

func (g *Goffeine) SetWithDelayAndDuration(key string, value any, delay int, duration time.Duration) {
	g.set(key, value, TimeUnion{Delay: delay, Duration: duration})
}

func (g *Goffeine) set(key string, value any, expireTime TimeUnion) {
	gnode := &node.GoffeineNode{Key: key, Value: value}
	v, ok := g.data.Load(key)
	if ok {
		oldEle := v.(*list.Element)
		oldEle.Value = gnode
		g.window.MoveToFront(oldEle)
		return
	}
	var ele *list.Element

	if g.windowIsFull() {
		// remove the key data
		ele = g.window.Back()
		oldKey := ele.Value.(*node.GoffeineNode).Key
		g.data.Delete(oldKey)

		// set value to back element and move it to front
		ele.Value = gnode
		g.window.MoveToFront(ele)
	} else {
		ele = g.window.PushFront(gnode)
	}
	g.data.Store(key, ele)

	//if g.windowIsFull() && g.probationIsFull() { // window 满了和probation都满了
	//	ele1 := g.window.Back()
	//	ele2 := g.probation.Back()
	//	g.evict(key, ele1, ele2)
	//}
	////else if g.window.Len() > g.windowSize
	////if g.window.Len() > g.windowSize { // window 满了，则要考虑和probation一起进行处理
	////	ele1 := g.window.Back()
	////	ele2 := g.probation.Back()
	////	g.compete(ele1, ele2)
	////}
	////g.data.Store(key, ele)
	//return ele
	//ele := v.(*list.Element)
	//ele.Value = value
}

//
//func (g *Goffeine) evict(key string, wEle *list.Element, pbEle *list.Element) {
//	cnt = g.counter[key]
//}
