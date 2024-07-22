package cache

import (
	"container/list"
	"sync"
	"time"
)

// AsyncQueryCache is a struct that manages an asynchronous query cache with LRU eviction policy.
type AsyncQueryCache struct {
	cache     map[string]*list.Element
	evictList *list.List
	capacity  int
	mutex     sync.RWMutex
}

// cacheItem is a struct that holds the cached item and its expiration.
type cacheItem struct {
	key        string
	value      string
	expiration time.Time
}

// NewAsyncQueryCache creates a new asynchronous query cache with LRU eviction policy.
func NewAsyncQueryCache(capacity int) *AsyncQueryCache {
	return &AsyncQueryCache{
		cache:     make(map[string]*list.Element),
		evictList: list.New(),
		capacity:  capacity,
	}
}

// Get is a method that retrieves a value from the cache.
func (aqc *AsyncQueryCache) Get(key string) (string, bool) {
	aqc.mutex.RLock()
	defer aqc.mutex.RUnlock()
	if elem, found := aqc.cache[key]; found {
		item := elem.Value.(*cacheItem)
		if !item.expiration.IsZero() && item.expiration.Before(time.Now()) {
			aqc.evictList.Remove(elem)
			delete(aqc.cache, item.key)
			return "", false
		}
		aqc.evictList.MoveToFront(elem)
		return item.value, true
	}
	return "", false
}

// Set is a method that sets a value in the cache asynchronously.
func (aqc *AsyncQueryCache) Set(key, value string) {
	go func() {
		aqc.mutex.Lock()
		defer aqc.mutex.Unlock()
		if elem, found := aqc.cache[key]; found {
			aqc.evictList.MoveToFront(elem)
			item := elem.Value.(*cacheItem)
			item.value = value
			return
		}
		item := &cacheItem{key: key, value: value}
		elem := aqc.evictList.PushFront(item)
		aqc.cache[key] = elem
		if aqc.evictList.Len() > aqc.capacity {
			lastElem := aqc.evictList.Back()
			if lastElem != nil {
				lastItem := lastElem.Value.(*cacheItem)
				delete(aqc.cache, lastItem.key)
				aqc.evictList.Remove(lastElem)
			}
		}
	}()
}

// SetWithExpiry is a method that sets a value in the cache with an expiration asynchronously.
func (aqc *AsyncQueryCache) SetWithExpiry(key, value string, duration time.Duration) {
	go func() {
		expiration := time.Now().Add(duration)
		aqc.mutex.Lock()
		defer aqc.mutex.Unlock()
		if elem, found := aqc.cache[key]; found {
			aqc.evictList.MoveToFront(elem)
			item := elem.Value.(*cacheItem)
			item.value = value
			item.expiration = expiration
			return
		}
		item := &cacheItem{key: key, value: value, expiration: expiration}
		elem := aqc.evictList.PushFront(item)
		aqc.cache[key] = elem
		if aqc.evictList.Len() > aqc.capacity {
			lastElem := aqc.evictList.Back()
			if lastElem != nil {
				lastItem := lastElem.Value.(*cacheItem)
				delete(aqc.cache, lastItem.key)
				aqc.evictList.Remove(lastElem)
			}
		}
	}()
}
