package goffeine

import (
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
// e.g.	goffeine.NewBuilder().MaximumSize(10).ExpireAfterWrite(time.Second, 5).Build()
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
		MaximumSize: b.maximumSize,
		ExpireTime:  b.expireTimeUnion,
		RefreshTime: b.refreshTimeUnion,
		data:        sync.Map{},
	}
}

// A Goffeine represents a cache
type Goffeine struct {
	MaximumSize int
	ExpireTime  TimeUnion
	RefreshTime TimeUnion
	data        sync.Map
}

func (g *Goffeine) Get(key string) (any, bool) {
	return g.data.Load(key)
}

func (g *Goffeine) Set(key string, value any) {
	g.set(key, value, g.ExpireTime)
}

func (g *Goffeine) SetWithDelay(key string, value any, delay int) {
	g.set(key, value, TimeUnion{Delay: delay, Duration: g.ExpireTime.Duration})
}

func (g *Goffeine) SetWithDelayAndDuration(key string, value any, delay int, duration time.Duration) {
	g.set(key, value, TimeUnion{Delay: delay, Duration: duration})
}
func (g *Goffeine) set(key string, value any, expireTime TimeUnion) {
	g.data.Store(key, value)
}
