package goffeine

import (
	"sync"
	"time"
)

// A GoffeineBuilder is used to create a Goffeine instance
// e.g.	goffeine.NewBuilder().maximumSize(10).ExpireAfterWrite(time.Second, 5).Build()
type GoffeineBuilder struct {
	maximumSize         int
	expireTimeDuration  time.Duration
	expireTimeDelay     int
	refreshTimeDuration time.Duration
	refreshTimeDelay    int
}

func (b *GoffeineBuilder) MaximumSize(size int) *GoffeineBuilder {
	if size < 1 {
		size = 3 // window: 1, probation: 1, protected: 1
	}
	b.maximumSize = size
	return b
}

func (b *GoffeineBuilder) ExpireAfterWrite(duration time.Duration, delay int) *GoffeineBuilder {
	b.expireTimeDuration = duration
	b.expireTimeDelay = delay
	return b
}

func (b *GoffeineBuilder) RefreshAfterWrite(duration time.Duration, delay int) *GoffeineBuilder {
	b.refreshTimeDuration = duration
	b.refreshTimeDelay = delay
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

		expireMilliseconds:  b.expireTimeDuration.Milliseconds() * int64(b.expireTimeDelay),
		refreshMilliseconds: b.refreshTimeDuration.Milliseconds() * int64(b.refreshTimeDelay),
		data:                &sync.Map{},
	}
}
