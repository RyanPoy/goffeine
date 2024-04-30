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
	return &Goffeine{
		maximumSize: b.maximumSize,
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

	window list.List
}

func (g *Goffeine) MaximumSize() int       { return g.maximumSize }
func (g *Goffeine) ExpireTime() TimeUnion  { return g.expireTime }
func (g *Goffeine) RefreshTime() TimeUnion { return g.refreshTime }
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
