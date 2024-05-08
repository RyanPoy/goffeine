package goffeine

import (
	"container/list"
	"sync"
	"time"
)

func NewBuilder() *Builder {
	return &Builder{}
}

// A Builder is used to create a Goffeine instance
// e.g.	goffeine.NewBuilder().maximumSize(10).ExpireAfterWrite(time.Second, 5).Build()
type Builder struct {
	maximumSize         int
	expireMilliseconds  int64
	refreshMilliseconds int64
}

func (b *Builder) MaximumSize(size int) *Builder {
	if size < 1 {
		size = 3 // window: 1, probation: 1, protected: 1
	}
	b.maximumSize = size
	return b
}

func (b *Builder) ExpireAfterWrite(duration time.Duration, delay int) *Builder {
	b.expireMilliseconds = duration.Milliseconds() * int64(delay)
	return b
}

func (b *Builder) RefreshAfterWrite(duration time.Duration, delay int) *Builder {
	b.refreshMilliseconds = duration.Milliseconds() * int64(delay)
	return b
}

func (b *Builder) Build() *Goffeine {
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
		window:               list.List{},
		probationMaximumSize: probationMaxsize,
		probation:            list.List{},
		protectedMaximumSize: protectedMaxsize,
		protected:            list.List{},
		expireMilliseconds:   b.expireMilliseconds,
		refreshMilliseconds:  b.refreshMilliseconds,
		data:                 &sync.Map{},
		fsketch:              NewSketch(b.maximumSize),
	}
}
