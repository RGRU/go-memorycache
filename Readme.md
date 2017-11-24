# go-memoryca
Manager memory key:value store/cache in Golang


## How to install?

	go get github.com/RGRU/go-memoryca


## How to use it?

First you must import it

	import (
		memoryca "github.com/RGRU/go-memoryca"
	)

Init a new Cache

	cache := memoryca.New("testDB", 10*time.Minute, 10*time.Minute)


Use it like this:

	cache.Set("myKey", "My value", 5 * time.Minute)
	cache.Get("myKey")
	cache.Exist("myKey")
	cache.Delete("myKey")
