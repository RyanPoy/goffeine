package goffeine

import (
	"sync"
	"time"
)

// A GoffeineBuilder is used to create a Goffeine instance
// e.g.	goffeine.NewBuilder().maximumSize(10).ExpireAfterWrite(time.Second, 5).Build()
type GoffeineBuilder struct {
	maximumSize         int
	expireMilliseconds  int64
	refreshMilliseconds int64
}

func (b *GoffeineBuilder) MaximumSize(size int) *GoffeineBuilder {
	if size < 1 {
		size = 3 // window: 1, probation: 1, protected: 1
	}
	b.maximumSize = size
	return b
}

func (b *GoffeineBuilder) ExpireAfterWrite(duration time.Duration, delay int) *GoffeineBuilder {
	b.expireMilliseconds = duration.Milliseconds() * int64(delay)
	return b
}

func (b *GoffeineBuilder) RefreshAfterWrite(duration time.Duration, delay int) *GoffeineBuilder {
	b.refreshMilliseconds = duration.Milliseconds() * int64(delay)
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

		expireMilliseconds:  b.expireMilliseconds,
		refreshMilliseconds: b.refreshMilliseconds,
		data:                &sync.Map{},
	}
}
