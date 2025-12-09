package core

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value      interface{}
	Expiration int64
}

type CacheService struct {
	items sync.Map
}

var (
	instance *CacheService
	once     sync.Once
)

func GetCacheService() *CacheService {
	once.Do(func() {
		instance = &CacheService{}
		go instance.startCleanup()
	})
	return instance
}

func (c *CacheService) Set(key string, value interface{}, duration time.Duration) {
	expiration := time.Now().Add(duration).UnixNano()
	c.items.Store(key, CacheItem{
		Value:      value,
		Expiration: expiration,
	})
}

func (c *CacheService) Get(key string) (interface{}, bool) {
	item, ok := c.items.Load(key)
	if !ok {
		return nil, false
	}

	cacheItem := item.(CacheItem)
	if time.Now().UnixNano() > cacheItem.Expiration {
		c.items.Delete(key)
		return nil, false
	}

	return cacheItem.Value, true
}

func (c *CacheService) Delete(key string) {
	c.items.Delete(key)
}

func (c *CacheService) startCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		c.items.Range(func(key, value interface{}) bool {
			item := value.(CacheItem)
			if time.Now().UnixNano() > item.Expiration {
				c.items.Delete(key)
			}
			return true
		})
	}
}
