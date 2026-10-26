package invezgo

import (
	"container/list"
	"sync"
	"time"
)

// lruCache stores raw response bodies (not decoded values): smallest footprint, and every hit
// decodes a fresh value so callers can never mutate cached data.
type lruCache struct {
	mu         sync.Mutex
	ll         *list.List // front = most recently used
	items      map[string]*list.Element
	bytes      int64
	maxEntries int
	maxBytes   int64
}

type cacheEntry struct {
	key     string
	body    []byte
	expires time.Time
}

func newLRUCache(maxEntries int, maxBytes int64) *lruCache {
	return &lruCache{ll: list.New(), items: map[string]*list.Element{}, maxEntries: maxEntries, maxBytes: maxBytes}
}

func (c *lruCache) get(key string, now time.Time) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	el, ok := c.items[key]
	if !ok {
		return nil, false
	}
	e := el.Value.(*cacheEntry)
	if !now.Before(e.expires) {
		c.remove(el)
		return nil, false
	}
	c.ll.MoveToFront(el)
	return e.body, true
}

func (c *lruCache) set(key string, body []byte, ttl time.Duration, now time.Time) {
	size := int64(len(body))
	if ttl <= 0 || size > c.maxBytes {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if el, ok := c.items[key]; ok {
		c.remove(el)
	}
	c.items[key] = c.ll.PushFront(&cacheEntry{key: key, body: body, expires: now.Add(ttl)})
	c.bytes += size
	for c.ll.Len() > c.maxEntries || c.bytes > c.maxBytes {
		c.remove(c.ll.Back())
	}
}

// purge removes every entry whose key matches.
func (c *lruCache) purge(match func(key string) bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for el := c.ll.Front(); el != nil; {
		next := el.Next()
		if match(el.Value.(*cacheEntry).key) {
			c.remove(el)
		}
		el = next
	}
}

func (c *lruCache) clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ll.Init()
	c.items = map[string]*list.Element{}
	c.bytes = 0
}

func (c *lruCache) len() (int, int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ll.Len(), c.bytes
}

// remove must be called with mu held.
func (c *lruCache) remove(el *list.Element) {
	e := c.ll.Remove(el).(*cacheEntry)
	delete(c.items, e.key)
	c.bytes -= int64(len(e.body))
}
