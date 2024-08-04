package cache_test

import (
	"testing"
	"time"

	"github.com/faciam-dev/goquent-query-builder/cache"
)

func TestAsyncQueryCache(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() cache.Cache
		key         string
		expected    string
		shouldExist bool
		waitTime    time.Duration
	}{
		{
			"Set and Get",
			func() cache.Cache {
				cache := cache.NewAsyncQueryCache(100)
				cache.Set("SELECT * FROM users", "result")
				return cache
			},
			"SELECT * FROM users",
			"result",
			true,
			100 * time.Millisecond,
		},
		{
			"Get Nonexistent Key",
			func() cache.Cache {
				return cache.NewAsyncQueryCache(100)
			},
			"SELECT * FROM users",
			"",
			false,
			100 * time.Millisecond,
		},
		{
			"Set With Expiry",
			func() cache.Cache {
				cache := cache.NewAsyncQueryCache(100)
				cache.SetWithExpiry("SELECT * FROM users", "result", 200*time.Millisecond)
				return cache
			},
			"SELECT * FROM users",
			"result",
			true,
			100 * time.Millisecond,
		},
		{
			"Get Expired Key",
			func() cache.Cache {
				cache := cache.NewAsyncQueryCache(100)
				cache.SetWithExpiry("SELECT * FROM users", "result", 200*time.Millisecond)
				return cache
			},
			"SELECT * FROM users",
			"",
			false,
			300 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := tt.setup()
			time.Sleep(tt.waitTime)
			result, found := cache.Get(tt.key)
			if found != tt.shouldExist {
				t.Errorf("expected found to be %v but got %v", tt.shouldExist, found)
			}
			if result != tt.expected {
				t.Errorf("expected '%s' but got '%s'", tt.expected, result)
			}
		})
	}
}
