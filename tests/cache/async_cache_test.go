package cache_test

import (
	"testing"
	"time"

	"github.com/faciam-dev/goquent-query-builder/internal/cache"
)

func TestAsyncQueryCache(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *cache.AsyncQueryCache
		key      string
		value    string
		expected string
		found    bool
		waitTime time.Duration
	}{
		{
			"Set and Get",
			func() *cache.AsyncQueryCache {
				cache := cache.NewAsyncQueryCache()
				cache.Set("SELECT * FROM users", "result")
				return cache
			},
			"SELECT * FROM users",
			"result",
			"result",
			true,
			100 * time.Millisecond,
		},
		{
			"Get Nonexistent Key",
			func() *cache.AsyncQueryCache {
				return cache.NewAsyncQueryCache()
			},
			"SELECT * FROM users",
			"",
			"",
			false,
			100 * time.Millisecond,
		},
		{
			"Set With Expiry",
			func() *cache.AsyncQueryCache {
				cache := cache.NewAsyncQueryCache()
				cache.SetWithExpiry("SELECT * FROM users", "result", 200*time.Millisecond)
				return cache
			},
			"SELECT * FROM users",
			"result",
			"result",
			true,
			100 * time.Millisecond,
		},
		{
			"Get Expired Key",
			func() *cache.AsyncQueryCache {
				cache := cache.NewAsyncQueryCache()
				cache.SetWithExpiry("SELECT * FROM users", "result", 200*time.Millisecond)
				return cache
			},
			"SELECT * FROM users",
			"",
			"",
			false,
			300 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cache := tt.setup()
			time.Sleep(tt.waitTime)
			result, found := cache.Get(tt.key)
			if found != tt.found {
				t.Errorf("expected found to be %v but got %v", tt.found, found)
			}
			if result != tt.expected {
				t.Errorf("expected '%s' but got '%s'", tt.expected, result)
			}
		})
	}
}
