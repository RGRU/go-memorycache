package memorycache

// TODO поиск по ключу и значению (регулярка)
// Сортировка и вывод множественных значений
// Инкримент, дикримент
// Конкатенация строк
// Экспорт & импорт в файл

import (
	"errors"
	"strings"
	"sync"
	"testing"

	"time"
)

// Cache struct cache
type Cache struct {
	sync.RWMutex
	items             map[string]Item
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
}

// Item struct cache item
type Item struct {
	Value      interface{}
	Expiration int64
	Created    time.Time
	Duration   time.Duration
}

// New. Initializing a new memory cache
func New(defaultExpiration, cleanupInterval time.Duration) *Cache {

	items := make(map[string]Item)

	// cache item
	cache := Cache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	if cleanupInterval > 0 {
		cache.StartGC()
	}

	return &cache
}

// Set setting a cache by key
func (c *Cache) Set(key string, value interface{}, duration time.Duration) error {

	var expiration int64

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()

	defer c.Unlock()

	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
		Duration:   duration,
	}

	return nil

}

// Get getting a cache by key
func (c *Cache) Get(key string) (interface{}, bool) {

	c.RLock()

	item, found := c.items[key]

	// cache not found
	if !found {
		c.RUnlock()
		return nil, false
	}

	if item.Expiration > 0 {

		// cache expired
		if time.Now().UnixNano() > item.Expiration {
			c.RUnlock()
			return nil, false
		}

	}

	c.RUnlock()

	return item.Value, true
}

// GetItem getting item cache
// Second parameter returns false if cache not found or expired
func (c *Cache) GetItem(key string) (*Item, bool) {

	c.RLock()

	item, found := c.items[key]

	// cache not found
	if !found {
		c.RUnlock()
		return nil, false
	}

	if item.Expiration > 0 {

		// cache expired
		if time.Now().UnixNano() > item.Expiration {
			c.RUnlock()
			return nil, false
		}

	}

	c.RUnlock()

	return &item, true
}

// GetCount return count items, without expired
func (c *Cache) GetCount() int {

	var count int

	for _, i := range c.items {

		// if cache no expired
		if !i.Expire() {

			count++

		}

	}

	return count

}

// Delete cache by key
// Return false if key not found
func (c *Cache) Delete(key string) error {

	c.Lock()

	defer c.Unlock()

	if _, found := c.items[key]; !found {
		return errors.New("Key not found")
	}

	delete(c.items, key)

	return nil
}

// Exists check cache exist
func (c *Cache) Exists(key string) bool {

	c.RLock()

	defer c.RUnlock()

	if value, found := c.items[key]; found {
		return !value.Expire()
	}

	return false
}

// Expire check cache expire
// Return true if cache expired
func (i *Item) Expire() bool {

	if i == nil {
		return false
	}

	return time.Now().UnixNano() > i.Expiration
}

// FlushAll delete all keys in database
func (c *Cache) FlushAll() error {

	c.Lock()

	defer c.Unlock()

	for k := range c.items {
		delete(c.items, k)
	}

	return nil
}

// Rename rename key to newkey
// When renaming, the value with the key is deleted
// Return false in key eq newkey
// Return false if newkey is exist
func (c *Cache) Rename(key string, newKey string) error {

	if key == newKey {
		return errors.New("The new name can not be the same with the old one")
	}

	_, foundNewKey := c.items[newKey]

	if foundNewKey {
		return errors.New("A key with a new name already exists")
	}

	item, found := c.items[key]

	if !found {
		return errors.New("Key not found")
	}

	c.Lock()

	defer c.Unlock()

	c.items[newKey] = item

	delete(c.items, key)

	return nil
}

// Copy copying a value from key to a newkey
// Return error if key eq newkey
// Return error if key is empty
// Return error if newkey is exist
func (c *Cache) Copy(key string, newKey string) error {

	if key == newKey {
		return errors.New("The new name can not be the same with the old one")
	}

	_, foundNewKey := c.items[newKey]

	if foundNewKey {
		return errors.New("A key with a new name already exists")
	}

	c.Lock()

	defer c.Unlock()

	item, found := c.items[key]

	if !found {
		return errors.New("Key not found")
	}

	c.items[newKey] = item

	return nil
}

// GetLikeKey list of keys by mask
// findString% - start of string
// %findString - end of string
// %findString% - any occurrence of a string
// findString - full occurrence of a string
func (c *Cache) GetLikeKey(search string) ([]interface{}, bool) {

	var values []interface{}

	if len(c.items) == 0 {
		return nil, false
	}

	search, like := parseLikeString(search)

	ls := len(search)

	// full
	if like == 3 {

		item, found := c.Get(search)

		if !found {
			return nil, false
		}

		values = append(values, item)

		return values, true

	}

	c.RLock()

	for k, i := range c.items {

		// search string in key name
		index := strings.Index(k, search)

		if index > -1 {

			// if cache no expired
			if !i.Expire() {

				switch {
				case like == 0 && index == 0: // start
					values = append(values, i.Value)
					break
				case like == 1 && index == len(k)-ls: // end
					values = append(values, i.Value)
					break
				case like == 2: // middle
					values = append(values, i.Value)
					break
				}

			}

		}

	}

	if len(values) == 0 {
		c.RUnlock()
		return nil, false
	}

	c.RUnlock()

	return values, true

}

// parseLikeString parse string param like
// return required string, 0: start | 1: middle | 2: end | 3: full
func parseLikeString(s string) (string, int) {

	count := strings.Count(s, "%")

	// required string
	search := strings.Replace(s, "%", "", 2)

	switch count {
	case 0: // full
		return search, 3
	case 2: // middle
		return search, 1
	}

	i := strings.Index(s, "%")

	if i == 0 { // end
		return search, 2
	}

	// start
	return search, 0

}

// StartGC start Garbage Collection
func (c *Cache) StartGC() error {

	go c.GC()

	return nil

}

// GC Garbage Collection
func (c *Cache) GC() {

	// if c.cleanupInterval < 1 {
	//
	// 	return
	// }

	for {

		<-time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		// fmt.Println(c.expiredKeys())

		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)

		}

	}

}

// expiredKeys returns key list which are expired.
func (c *Cache) expiredKeys() (keys []string) {

	c.RLock()

	defer c.RUnlock()

	for k, i := range c.items {
		if i.Expire() {
			keys = append(keys, k)
		}
	}

	return
}

// clearItems removes all the items which key in keys.
func (c *Cache) clearItems(keys []string) {

	c.Lock()

	defer c.Unlock()

	for _, k := range keys {
		delete(c.items, k)
	}
}

// Get getting a cache by key without expire check
func (c *Cache) getWithOutExpire(key string) (interface{}, bool) {

	c.RLock()

	item, found := c.items[key]

	// cache not found
	if !found {
		c.RUnlock()
		return nil, false
	}

	c.RUnlock()

	return item.Value, true
}

// benchmarkGet benchmark set cahce
func (c *Cache) benchmarkSet(b *testing.B) {

	for n := 0; n < b.N; n++ {

		c.Set("testKey:"+string(b.N), "testValue"+string(b.N), 1*time.Minute)

	}

}

// benchmarkGet benchmark get cahce
func (c *Cache) benchmarkGet(b *testing.B) {

	for n := 0; n < b.N; n++ {

		c.Get("testKey:" + string(b.N))

	}

}

// benchmarkGet benchmark get like key
func (c *Cache) benchmarkGetLikeKey(b *testing.B) {

	for n := 0; n < b.N; n++ {

		c.GetLikeKey("testKey%")

	}

}
