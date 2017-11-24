package memoryca

import (
	"testing"
	"time"
)

const (
	testKey        string = "cache:test"
	incrKey        string = "cache:incr"
	testKeyEmpty   string = "cache:empty"
	testValue      string = "Hello Test"
	testKeyExpired string = "cache:expired"
	testRenameKey  string = "cache:renamekey"
	testCopyKey    string = "cache:copykey"
)

// AppCache init new cache
var AppCache = New("testDB", 10*time.Minute, 1*time.Second)

// TestSet set cache
func TestSet(t *testing.T) {

	s := AppCache.Set(testKey, testValue, 1*time.Minute)

	// fmt.Println("TestSet:", s)

	if s != nil {
		t.Error("Error: ", "Could not set cache")
	}

}

// TestGet get cache by key
func TestGet(t *testing.T) {

	AppCache.Set(testKey, testValue, 1*time.Minute)

	value, found := AppCache.Get(testKey)

	if value != testValue {
		t.Error("Error: ", "The received value: do not correspond to the expectation:", value, testValue)
	}

	if found != true {
		t.Error("Error: ", "Could not get cache")
	}

	// get cache by key is empty
	value, found = AppCache.Get(testKeyEmpty)

	if value != nil || found != false {
		t.Error("Error: ", "Value does not exist and must be empty", value)
	}
}

// TestGetItem get cahce item
func TestGetItem(t *testing.T) {

	AppCache.Set(testKey, testValue, 1*time.Millisecond)

	// sleep to make the cache expired
	time.Sleep(2 * time.Millisecond)

	// get cache by key is empty
	value, found := AppCache.GetItem(testKey)

	if value != nil || found != false {
		t.Error("Error: ", "Value does not exist and must be empty", value)
	}

}

// TestGetCount get count
func TestGetCount(t *testing.T) {

	AppCache.FlushAll()

	AppCache.Set("one:count", testValue, 10*time.Second)
	AppCache.Set("two:count", testValue, 10*time.Second)
	AppCache.Set("three:count", testValue, 10*time.Second)

	count := AppCache.GetCount()

	if count != 3 {
		t.Error("Error: ", "Count items does not match the expectation", count)
	}

}

// TestGetEmpty get cache by key expired
func TestGetExpired(t *testing.T) {

	s := AppCache.Set(testKeyExpired, testValue, 1*time.Millisecond)

	if s != nil {
		t.Error("Error: ", "Could not set cache")
	}

	// sleep to make the cache expired
	time.Sleep(2 * time.Millisecond)

	value, found := AppCache.Get(testKeyExpired)

	if value != nil || found != false {
		t.Error("Error: ", "Cache Expired and must be empty", value)
	}

}

// TestDelete delete cache by key
func TestDelete(t *testing.T) {

	AppCache.Set(testKey, testValue, 1*time.Minute)

	error := AppCache.Delete(testKey)

	if error != nil {
		t.Error("Error: ", "Cache delete failed")
	}

	value, found := AppCache.Get(testKey)

	if found {
		t.Error("Error: ", "Should not be found because it was deleted")
	}

	if value != nil {
		t.Error("Error: ", "Value is not nil:", value)
	}

	// repeat deletion of an existing cache
	error = AppCache.Delete(testKeyEmpty)

	if error == nil {
		t.Error("Error: ", "An empty cache should return an error")
	}

}

// TestExists check exists cache
func TestExists(t *testing.T) {

	AppCache.Set(testKey, testValue, 1*time.Minute)

	ok := AppCache.Exists(testKey)

	if ok != true {
		t.Error("Error: ", "Value is empty")
	}

	// check exists cache is empty
	ok = AppCache.Exists(testKeyEmpty)

	if ok != false {
		t.Error("Error: ", "Value must be empty")
	}

}

// TestExpire check expire cache
func TestExpire(t *testing.T) {

	AppCache.Set(testKey, testValue, 2*time.Millisecond)

	value, found := AppCache.GetItem(testKey)

	if !found {
		t.Error("Error: ", "Cache not found")
	}

	expired := value.Expire()

	if expired == true {
		t.Error("Error: ", "The cache must be up-to-date")
	}

	// check empty key
	value, _ = AppCache.GetItem(testKeyEmpty)

	expired = value.Expire()

	if expired != false {
		t.Error("Error: ", "An empty cache should return an error")
	}

}

// TestFlushAll delete all the keys
func TestFlushAll(t *testing.T) {

	AppCache.Set(testKey, testValue, 5*time.Minute)

	error := AppCache.FlushAll()

	if error != nil {
		t.Error("Error: ", "Could not flush all cache")
	}

	value, found := AppCache.Get(testKey)

	if found {
		t.Error("Error: ", "Should not be found because it was deleted")
	}

	if value != nil {
		t.Error("Error: ", "Value is not nil:", value)
	}

}

// TestRename rename key to newkey
func TestRename(t *testing.T) {

	AppCache.Set(testKey, testValue, 5*time.Minute)

	error := AppCache.Rename(testKey, testRenameKey)

	if error != nil {
		t.Error("Error: ", "Error renaming key")
	}

	value, found := AppCache.Get(testRenameKey)

	if !found {
		t.Error("Error: ", "Cache not found")
	}

	if value != testValue {
		t.Error("Error: ", "The received value: do not correspond to the expectation:", value, testValue)
	}

	// with duplicate keys
	error = AppCache.Rename(testRenameKey, testRenameKey)

	if error == nil {
		t.Error("Error: ", "The name of the keys can not be the same")
	}

	// with empty key
	error = AppCache.Rename(testKeyEmpty, testRenameKey)

	if error == nil {
		t.Error("Error: ", "Cache with the specified key not found")
	}

	// copying two empty keys
	error = AppCache.Rename("empty:three", "empty:four")

	if error == nil {
		t.Error("Error: ", "Cache with the specified key not found")
	}

}

// TestCopy copy value key to newkey
func TestCopy(t *testing.T) {

	AppCache.Set(testKey, testValue, 5*time.Minute)

	error := AppCache.Copy(testKey, testCopyKey)

	if error != nil {
		t.Error("Error: ", "Error copy key")
	}

	value, found := AppCache.Get(testCopyKey)

	if !found {
		t.Error("Error: ", "Cache not found")
	}

	if value != testValue {
		t.Error("Error: ", "The received value: do not correspond to the expectation:", value, testValue)
	}

	// with duplicate keys
	error = AppCache.Copy(testCopyKey, testCopyKey)

	if error == nil {
		t.Error("Error: ", "The name of the keys can not be the same")
	}

	// with empty key
	error = AppCache.Copy(testKeyEmpty, testCopyKey)

	if error == nil {
		t.Error("Error: ", "Cache with the specified key not found")
	}

	// copying two empty keys
	error = AppCache.Copy("empty:one", "empty:two")

	if error == nil {
		t.Error("Error: ", "Copying two empty keys")
	}

}

func TestStartGC(t *testing.T) {

	AppCache.Set(testKey, testValue, 1*time.Second)

	// get cahce before run GC
	value, found := AppCache.getWithOutExpire(testKey)

	if value != testValue {
		t.Error("Error: ", "The received value: do not correspond to the expectation:", value, testValue)
	}

	if found != true {
		t.Error("Error: ", "Could not get cache")
	}

	// start garbage collector
	AppCache.StartGC()

	// sleep to make the cache expired
	time.Sleep(3 * time.Second)

	// get cahce after run GC
	value, found = AppCache.getWithOutExpire(testKey)

	if found {
		t.Error("Error: ", "Cache Expired and must be empty", value)
	}

}

// // BenchmarkSet benchmark Set
// func BenchmarkSet(b *testing.B) {
//
// 	AppCache.benchmarkSet(b)
//
// }
//
// // BenchmarkGet benchmark Get
// func BenchmarkGet(b *testing.B) {
//
// 	AppCache.benchmarkGet(b)
//
// }
