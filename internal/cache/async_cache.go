package cache

import (
	"sync"
	"time"
)

// AsyncQueryCache は、非同期キャッシュを管理する構造体です。
type AsyncQueryCache struct {
	cache map[string]cacheItem
	mutex sync.RWMutex
}

// cacheItem は、キャッシュされたアイテムとその有効期限を保持します。
type cacheItem struct {
	value      string
	expiration time.Time
}

// NewAsyncQueryCache は、新しい非同期キャッシュを作成します。
func NewAsyncQueryCache() *AsyncQueryCache {
	return &AsyncQueryCache{cache: make(map[string]cacheItem)}
}

// Get は、キャッシュから値を取得します。
func (aqc *AsyncQueryCache) Get(key string) (string, bool) {
	aqc.mutex.RLock()
	defer aqc.mutex.RUnlock()
	item, found := aqc.cache[key]
	if !found || (!item.expiration.IsZero() && item.expiration.Before(time.Now())) {
		return "", false
	}
	return item.value, true
}

// Set は、キャッシュに値を非同期で設定します。
func (aqc *AsyncQueryCache) Set(key, value string) {
	go func() {
		aqc.mutex.Lock()
		defer aqc.mutex.Unlock()
		aqc.cache[key] = cacheItem{value: value}
	}()
}

// SetWithExpiry は、キャッシュに有効期限付きで値を非同期で設定します。
func (aqc *AsyncQueryCache) SetWithExpiry(key, value string, duration time.Duration) {
	go func() {
		expiration := time.Now().Add(duration)
		aqc.mutex.Lock()
		defer aqc.mutex.Unlock()
		aqc.cache[key] = cacheItem{value: value, expiration: expiration}
	}()
}
