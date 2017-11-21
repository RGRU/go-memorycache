package cache

import (
	"lovefrontend/models/cache/memory"
)

// Cache general cache structure
type Cache struct {
	Memory memory.Memory
	// Redis Redis
}

// NewCache initializing a new cache
func NewCache(adapter string, config string) *Cache {

	if adapter == "memory" {

		c := memory.NewMemoryCache(config)

		return &Cache{c}

	}

	return nil

}
