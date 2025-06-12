package memcache

import (
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

// MemCacheOptions is a functional option type for configuring MemCache
type MemCacheOptions func(*memCacheOptions)

// memCacheOptions holds the configuration options for MemCache
type memCacheOptions struct {
	numCounters              int64
	maxCost                  int64
	bufferItems              int64
	metrics                  bool
	thetlTickerDurationInSec int64
}

// defaultMemCacheOptions returns the default options for MemCache
func defaultMemCacheOptions() *memCacheOptions {
	return &memCacheOptions{
		numCounters:              1e4,
		maxCost:                  1 << 23,
		bufferItems:              64,
		metrics:                  false,
		thetlTickerDurationInSec: 10,
	}
}

/*
NumCounters determines the number of counters (keys) to keep that hold access frequency information. It's generally a good idea to have more counters than the max cache capacity, as this will improve eviction accuracy and subsequent hit ratios.

For example, if you expect your cache to hold 1,000,000 items when full, NumCounters should be 10,000,000 (10x). Each counter takes up roughly 3 bytes (4 bits for each counter * 4 copies plus about a byte per counter for the bloom filter). Note that the number of counters is internally rounded up to the nearest power of 2, so the space usage may be a little larger than 3 bytes * NumCounters.

We've seen good performance in setting this to 10x the number of items you expect to keep in the cache when full.
*/
func WithNumCounters(numCounters int64) MemCacheOptions {
	return func(opts *memCacheOptions) {
		opts.numCounters = numCounters
	}
}

/*
MaxCost is how eviction decisions are made. For example, if MaxCost is 100 and a new item with a cost of 1 increases total cache cost to 101, 1 item will be evicted.

MaxCost can be considered as the cache capacity, in whatever units you choose to use.

For example, if you want the cache to have a max capacity of 100MB, you would set MaxCost to 100,000,000 and pass an item's number of bytes as the `cost` parameter for calls to Set. If new items are accepted, the eviction process will take care of making room for the new item and not overflowing the MaxCost value.

MaxCost could be anything as long as it matches how you're using the cost values when calling Set.
*/
func WithMaxCost(maxCost int64) MemCacheOptions {
	return func(opts *memCacheOptions) {
		opts.maxCost = maxCost
	}
}

/*
BufferItems determines the size of Get buffers.

Unless you have a rare use case, using `64` as the BufferItems value results in good performance.

If for some reason you see Get performance decreasing with lots of contention (you shouldn't), try increasing this value in increments of 64. This is a fine-tuning mechanism and you probably won't have to touch this.
*/
func WithBufferItems(bufferItems int64) MemCacheOptions {
	return func(opts *memCacheOptions) {
		opts.bufferItems = bufferItems
	}
}

/*
Metrics is true when you want variety of stats about the cache. There is some overhead to keeping statistics, so you should only set this flag to true when testing or throughput performance isn't a major factor.
*/
func WithMetrics(metrics bool) MemCacheOptions {
	return func(opts *memCacheOptions) {
		opts.metrics = metrics
	}
}

/*
TtlTickerDurationInSec sets the value of time ticker for cleanup keys on TTL expiry.
*/
func WithTtlTickerDurationInSec(ttlTickerDurationInSec int64) MemCacheOptions {
	return func(opts *memCacheOptions) {
		opts.thetlTickerDurationInSec = ttlTickerDurationInSec
	}
}

// func NewWebhookLogger(webhookURL string, ServiceName string, loggeropts ...func(*LoggerOptions)) *zap.Logger {

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
func NewMemCache[V any](ttl time.Duration, opts ...MemCacheOptions) (Cache[V], error) {

	defaultOpts := defaultMemCacheOptions()
	for _, opt := range opts {
		opt(defaultOpts)
	}

	c, err := ristretto.NewCache(&ristretto.Config[string, V]{
		NumCounters:            defaultOpts.numCounters,
		MaxCost:                defaultOpts.maxCost,
		BufferItems:            defaultOpts.bufferItems,
		Metrics:                defaultOpts.metrics,
		TtlTickerDurationInSec: defaultOpts.thetlTickerDurationInSec,
	})

	if err != nil {
		return nil, err
	}

	return &MemCache[V]{cache: c, ttl: ttl}, nil
}

// Get retrieves a value from the cache by key. It returns the value and a boolean indicating if the key was found.
func (rc *MemCache[V]) Get(key string) (V, bool) {
	return rc.cache.Get(key)
}

// Set adds a value to the cache with a specified key. It returns true if the value was successfully set, false otherwise.
func (rc *MemCache[V]) Set(key string, value V) bool {
	result := rc.cache.SetWithTTL(key, value, 1, rc.ttl)
	rc.cache.Wait()
	return result
}
