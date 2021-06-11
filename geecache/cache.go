package geecache

import (
	lru2 "GeeCache/geecache/LRU"
	"sync"
)

type cache struct {
	mu 		sync.Mutex
	lru 	*lru2.Cache
	cacheBytes int64
}

func (c *cache)lockadd(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru2.CacheNew(c.cacheBytes, nil)
	}
	c.lru.CacheAdd(key, value)
}

//增加锁控制的get方法
func (c *cache)lockget(key string)(value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.CacheGet(key);ok{
		return v.(ByteView), ok
	}
	return
}