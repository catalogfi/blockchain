package memcache

import (
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

// Cache is a generic interface for caching operations
type Cache[V any] interface {
	Get(key string) (V, bool)
	Set(key string, value V) bool
}

// MemCache is a generic wrapper around ristretto.Cache
type MemCache[V any] struct {
	cache *ristretto.Cache[string, V]
	ttl   time.Duration
}

// NewMemCache creates a new memory cache with the specified TTL
func NewMemCache[V any](ttl time.Duration) (Cache[V], error) {
	c, err := ristretto.NewCache(&ristretto.Config[string, V]{
		NumCounters:            1e4,
		MaxCost:                1 << 23,
		BufferItems:            32,
		Metrics:                false,
		TtlTickerDurationInSec: 10,
	})

	if err != nil {
		return nil, err
	}

	return &MemCache[V]{cache: c, ttl: ttl}, nil
}

func (rc *MemCache[V]) Get(key string) (V, bool) {
	return rc.cache.Get(key)
}

func (rc *MemCache[V]) Set(key string, value V) bool {
	result := rc.cache.SetWithTTL(key, value, 1, rc.ttl)
	rc.cache.Wait()
	return result
}
