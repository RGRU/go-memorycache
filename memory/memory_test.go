package memory

import (
	"testing"
	"time"
)

const (
	testKey      string = "cache:test"
	testKeyEmpty string = "cache:test:empty"
	testValue    string = "Hello Test!"
)

// AppCache переменная которая ссылается на ресурс кеша
var AppCache *Cache = &Cache{
	NewMemoryCache(`{"expiration":10, "interval":10}`),
}

// TestSet set cache (testKey)
func TestSet(t *testing.T) {

	c := AppCache.Memory.Set(testKey, []byte(testValue), 5*time.Minute)

	if c != true {
		t.Error("Error: ", "Could not set cache")
	}

}

// TestGet get cache by key (testKey)
func TestGet(t *testing.T) {

	v, ok := AppCache.Memory.Get(testKey)

	if string(v) != testValue {
		t.Error("Error: ", "The received value: "+string(v)+" do not correspond to the expectation:", testValue)
	}

	if ok != true {
		t.Error("Error: ", "Could not get cache")
	}

}

// TestIsExist checking the cache for existence
func TestIsExist(t *testing.T) {

	ve := AppCache.Memory.IsExist(testKeyEmpty)

	if ve != false {
		t.Error("Error: ", "There is a cache that should not be")
	}

	v := AppCache.Memory.IsExist(testKey)

	if v != true {
		t.Error("Error: ", "Cache not found")
	}

}

// TestDelete delete cache by key (testKey)
func TestDelete(t *testing.T) {

	AppCache.Memory.Delete(testKey)

	v := AppCache.Memory.IsExist(testKey)

	if v != false {
		t.Error("Error: ", "Cache removal failed")
	}

}
