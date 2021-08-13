package cache

import (
	"goffeine/cache/internal/fsketch"
	"goffeine/cache/internal/lru"
)

type Cache struct {
	cap       int
	sketch    *fsketch.FSketch
	admission *lru.LRU
	probation *lru.LRU
	protected *lru.LRU
}

func New(cap int) Cache {
	return Cache{
		cap:       cap,
		sketch:    fsketch.New(cap),
		admission: lru.New(cap),
		probation: lru.New(cap),
		protected: lru.New(cap),
	}
}

func (c *Cache) Capacity() int {
	return c.cap
}

func (c *Cache) Len() int {
	return c.admission.Len() + c.probation.Len() + c.protected.Len()
}

