// © 2019-2025 Diarkis Inc. All rights reserved.

package lodmanager

import (
	"sync"
	"time"
)

type cache struct {
	m map[cacheKey]time.Time
	k map[string]map[cacheKey]struct{}
	sync.RWMutex
}

type cacheKey struct {
	a string
	b string
}

func newCache() *cache {
	c := &cache{m: make(map[cacheKey]time.Time), k: make(map[string]map[cacheKey]struct{})}
	return c
}

func (c *cache) get(key cacheKey) time.Time {
	c.RLock()
	defer c.RUnlock()

	v := c.m[key]

	return v
}

func (c *cache) set(key cacheKey, data time.Time) {
	c.Lock()
	defer c.Unlock()

	c.m[key] = data

	if _, ok := c.k[key.a]; !ok {
		c.k[key.a] = make(map[cacheKey]struct{})
	}

	if _, ok := c.k[key.b]; !ok {
		c.k[key.b] = make(map[cacheKey]struct{})
	}

	c.k[key.a][key] = struct{}{}
	c.k[key.b][key] = struct{}{}
}

// clears all cache by cacheKey's a or b
func (c *cache) clearBy(cacheKeyFragment string) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.k[cacheKeyFragment]; !ok {
		return
	}

	keys := c.k[cacheKeyFragment]

	for key := range keys {
		delete(c.m, key)
	}
	delete(c.k, cacheKeyFragment)

	return
}

func (c *cache) clear() {
	c.Lock()
	defer c.Unlock()

	c.m = make(map[cacheKey]time.Time)
	c.k = make(map[string]map[cacheKey]struct{})
}
