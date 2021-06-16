package main

import (
	"github.com/gin-gonic/gin"
	"sync"
)

// Singleton
var bc bookCache
var cacheOnce sync.Once

// Singleton method
func GetCacheInstance() bookCache {
	cacheOnce.Do(func() {
		bc = GetCache()
	})
	return bc
}

func CacheGet(username string) ([]string, error) {
	curCache := GetCacheInstance()
	return curCache.CacheGet(username)
}

func CacheAdd(username string, ctx * gin.Context) {
	curCache := GetCacheInstance()
	curCache.CacheAdd(username, ctx)
}