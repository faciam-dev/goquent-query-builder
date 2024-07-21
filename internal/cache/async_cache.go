package cache

import (
	"sync"
	"time"
)

// AsyncQueryCache is a struct that manages an asynchronous query cache.
type AsyncQueryCache struct {
	cache map[string]cacheItem
	mutex sync.RWMutex
}

// cacheItem is a struct that holds the cached item and its expiration.
type cacheItem struct {
	value      string
	expiration time.Time
}

// NewAsyncQueryCache creates a new asynchronous query cache.
func NewAsyncQueryCache() *AsyncQueryCache {
	return &AsyncQueryCache{cache: make(map[string]cacheItem)}
}

// Get is a method that retrieves a value from the cache.
func (aqc *AsyncQueryCache) Get(key string) (string, bool) {
	aqc.mutex.RLock()
	defer aqc.mutex.RUnlock()
	item, found := aqc.cache[key]
	if !found || (!item.expiration.IsZero() && item.expiration.Before(time.Now())) {
		return "", false
	}
	return item.value, true
}

// Set is a method that sets a value in the cache asynchronously.
func (aqc *AsyncQueryCache) Set(key, value string) {
	go func() {
		aqc.mutex.Lock()
		defer aqc.mutex.Unlock()
		aqc.cache[key] = cacheItem{value: value}
	}()
}

// SetWithExpiry is a method that sets a value in the cache with an expiration asynchronously.
func (aqc *AsyncQueryCache) SetWithExpiry(key, value string, duration time.Duration) {
	go func() {
		expiration := time.Now().Add(duration)
		aqc.mutex.Lock()
		defer aqc.mutex.Unlock()
		aqc.cache[key] = cacheItem{value: value, expiration: expiration}
	}()
}
