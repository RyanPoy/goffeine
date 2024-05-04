package goffeine

import (
	"container/list"
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
	maximumSize int
	expireTime  TimeUnion
	refreshTime TimeUnion
	data        *sync.Map

	window            list.List
	windowMaximumSize int

	probation            map[string]any
	probationMaximumSize int

	protected            map[string]any
	protectedMaximumSize int
}

func (g *Goffeine) MaximumSize() int          { return g.maximumSize }
func (g *Goffeine) ExpireTime() TimeUnion     { return g.expireTime }
func (g *Goffeine) RefreshTime() TimeUnion    { return g.refreshTime }
func (g *Goffeine) WindowMaximumSize() int    { return g.windowMaximumSize }
func (g *Goffeine) ProbationMaximumSize() int { return g.probationMaximumSize }
func (g *Goffeine) ProtectedMaximumSize() int { return g.protectedMaximumSize }

func (g *Goffeine) Get(key string) (any, bool) {
	return g.data.Load(key)
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
	g.data.Store(key, value)
}

//func (g *Goffeine) set(key string, value any, expireTime TimeUnion) *list.Element {
//	v, ok := g.data.Load(key)
//	if !ok {
//		ele := g.window.PushFront(key)     // 不存在，则添加到头部
//		if g.window.Len() > g.windowSize && g.probation.Len() > g.probationSize { // window 满了和probation都满了，
//			ele1 := g.window.Back()
//			ele2 := g.probation.Back()
//			g.evict(ele1, ele2)
//		} else if g.window.Len() > g.windowSize
//		if g.window.Len() > g.windowSize { // window 满了，则要考虑和probation一起进行处理
//			ele1 := g.window.Back()
//			ele2 := g.probation.Back()
//			g.compete(ele1, ele2)
//		}
//		g.data.Store(key, ele)
//		return ele
//	}
//	ele := v.(*list.Element)
//	ele.Value = value
//	//ele
//}
//
//func (g *Goffeine) evict(key1 string, key2 string, ele1, ele2 *list.Element) {
//	if g.
//	if sketch.Frequency(string)
//}
