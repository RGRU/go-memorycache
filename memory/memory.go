package memory

import (
	"encoding/json"

	"time"

	cache "github.com/patrickmn/go-cache"
)

// Cache general cache structure
type Cache struct {
	Memory Memory
}

// Memory cache-in-memory interface
type Memory interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte, duration time.Duration) bool
	Delete(key string) error
	// Incr(key string) error
	// Decr(key string) error
	IsExist(key string) bool
	// ClearAll() error
	// StartAndGC(config string) error
}

// memoryCache cache-in-memory struct
type memoryCache struct {
	cache  *cache.Cache
	Memory Memory
}

// memoryCacheConf cache-in-memory config struct
type memoryCacheConf struct {
	Expiration int32 `json:"expiration"`
	Interval   int32 `json:"interval"`
}

// NewMemoryCache initializing a new memory cache
func NewMemoryCache(config string) *memoryCache {

	var conf memoryCacheConf

	json.Unmarshal([]byte(config), &conf)

	c := cache.New(time.Duration(conf.Expiration)*time.Minute, time.Duration(conf.Interval)*time.Minute)

	return &memoryCache{cache: c}

}

// Set setting a cache by key
func (c *memoryCache) Set(key string, value []byte, duration time.Duration) bool {

	c.cache.Set(key, value, duration)

	return true

}

// Get getting a cache by key
func (c *memoryCache) Get(key string) ([]byte, bool) {

	if value, found := c.cache.Get(key); found {

		return value.([]byte), true

	}

	return nil, false

}

// Delete delete a cache by key
func (c *memoryCache) Delete(key string) error {

	c.cache.Delete(key)

	return nil

}

// IsExist checking the cache for existence
func (c *memoryCache) IsExist(key string) bool {

	if _, found := c.cache.Get(key); found {

		return true

	}

	return false

}
