Language: English | [Русский](https://github.com/RGRU/go-memorycache/blob/master/ReadmeRU.md)

# Go-memorycache [![Build Status](https://travis-ci.org/RGRU/go-memorycache.svg?branch=master)](https://travis-ci.org/RGRU/go-memorycache)
Manager memory key:value store/cache in Golang


## How to install?

	go get github.com/RGRU/go-memorycache


## How to use it?

First you must import it

	import (
		memorycache "github.com/RGRU/go-memorycache"
	)

Init a new Cache

	cache := memorycache.New("testDB", 10*time.Minute, 10*time.Minute)


Use it like this:

	// Set cache by key
	cache.Set("myKey", "My value", 5 * time.Minute)

	// Get cache by key
	cache.Get("myKey")

	// Check exist cache
	cache.Exist("myKey")

	// Delete cache by key
	cache.Delete("myKey")
