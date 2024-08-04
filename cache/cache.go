package cache

import "time"

type Cache interface {
	Set(key string, value string)
	Get(key string) (string, bool)
	SetWithExpiry(key, value string, duration time.Duration)
}
