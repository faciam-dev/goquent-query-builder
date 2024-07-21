package cache

import (
	"time"
)

// BlankQueryCache is a struct that manages a blank cache.
type BlankQueryCache struct {
}

// cacheItem holds the cached item and its expiration.
func NewBlankQueryCache() *BlankQueryCache {
	return &BlankQueryCache{}
}

// Get は、キャッシュから値を取得します。
func (aqc *BlankQueryCache) Get(key string) (string, bool) {
	return "", false
}

// Set は、キャッシュに値を非同期で設定します。
func (aqc *BlankQueryCache) Set(key, value string) {

}

// SetWithExpiry は、キャッシュに有効期限付きで値を非同期で設定します。
func (aqc *BlankQueryCache) SetWithExpiry(key, value string, duration time.Duration) {

}
