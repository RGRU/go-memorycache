# Go-memoryca [![Build Status](https://travis-ci.org/RGRU/go-memoryca.svg?branch=master)](https://travis-ci.org/RGRU/go-memoryca)
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

	// Set cache by key
	cache.Set("myKey", "My value", 5 * time.Minute)

	// Get cache by key
	cache.Get("myKey")

	// Check exist cache
	cache.Exist("myKey")

	// Delete cache by key
	cache.Delete("myKey")
