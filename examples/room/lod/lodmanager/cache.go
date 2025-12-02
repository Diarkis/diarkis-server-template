// © 2019-2025 Diarkis Inc. All rights reserved.

package lodmanager

import (
	"sync"
)

type cache struct {
	m map[cacheKey][]byte
	k map[string]map[cacheKey]struct{}
	sync.RWMutex
}

type cacheKey struct {
	a string
	b string
}

func newCache() *cache {
	c := &cache{m: make(map[cacheKey][]byte), k: make(map[string]map[cacheKey]struct{})}
	return c
}

func (c *cache) get(key cacheKey) ([]byte, bool) {
	c.RLock()
	defer c.RUnlock()

	v, ok := c.m[key]

	return v, ok
}

func (c *cache) set(key cacheKey, data []byte) bool {
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

	return true
}

// clears all cache by cacheKey's a or b
func (c *cache) clearBy(cacheKeyFragment string) bool {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.k[cacheKeyFragment]; !ok {
		return false
	}

	keys := c.k[cacheKeyFragment]

	for key := range keys {
		delete(c.m, key)
	}
	delete(c.k, cacheKeyFragment)

	return true
}

func (c *cache) clear() {
	c.Lock()
	defer c.Unlock()

	c.m = make(map[cacheKey][]byte)
	c.k = make(map[string]map[cacheKey]struct{})
}
