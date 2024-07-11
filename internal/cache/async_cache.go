package cache

import (
	"sync"
	"time"
)

type AsyncQueryCache struct {
	cache map[string]string
	mutex sync.RWMutex
}

func NewAsyncQueryCache() *AsyncQueryCache {
	return &AsyncQueryCache{cache: make(map[string]string)}
}

func (aqc *AsyncQueryCache) Get(query string) (string, bool) {
	aqc.mutex.RLock()
	result, found := aqc.cache[query]
	aqc.mutex.RUnlock()
	return result, found
}

func (aqc *AsyncQueryCache) Set(query, result string) {
	go func() {
		aqc.mutex.Lock()
		aqc.cache[query] = result
		aqc.mutex.Unlock()
	}()
}

func (aqc *AsyncQueryCache) SetWithExpiry(query, result string, duration time.Duration) {
	go func() {
		aqc.mutex.Lock()
		aqc.cache[query] = result
		aqc.mutex.Unlock()
		time.Sleep(duration)
		aqc.mutex.Lock()
		delete(aqc.cache, query)
		aqc.mutex.Unlock()
	}()
}
