package service

import (
	"github.com/bar41234/bar_book_service/platform"
	"github.com/gin-gonic/gin"
	"sync"
)

// Singleton
var bc platform.BookCache
var cacheOnce sync.Once

// Singleton method
func GetCacheInstance() platform.BookCache {
	cacheOnce.Do(func() {
		bc = platform.GetCache()
	})
	return bc
}

func CacheGet(username string) ([]string, error) {
	curCache := GetCacheInstance()
	return curCache.CacheGet(username)
}

func CacheAdd(username string, ctx *gin.Context) {
	curCache := GetCacheInstance()
	curCache.CacheAdd(username, ctx)
}
