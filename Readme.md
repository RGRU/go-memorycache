# go-cache
Сache is a Go cache manager


## How to install?

	go get github.com/RGRU/go-cache


## How to use it?

First you must import it

	import (
		cache "github.com/RGRU/go-cache"
	)

Init a Cache (example with memory adapter)

	AppCache := cache.NewCache("memory", `{"expiration":10, "interval":10}`)
    or outside the function
    var AppCache *cache.Cache = cache.NewCache("memory", `{"expiration":10, "interval":10}`)

Use it like this:

	cache.AppCache.Memory.Set("myKey", "My value", 10 * time.Minute)
	cache.AppCache.Memory.Get("myKey")
	cache.AppCache.Memory.IsExist("myKey")
	cache.AppCache.Memory.Delete("myKey")


## Memory adapter

Configure memory adapter like this:

	{"expiration":10, "interval":10}
