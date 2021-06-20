package platform

import (
	"github.com/bar41234/bar_book_service/cache"
	"github.com/bar41234/bar_book_service/datastore"
	"github.com/bar41234/bar_book_service/models"
	"github.com/gin-gonic/gin"
)

type BookManipulator interface {
	Get(id string) (*models.Book, error)
	Put(book models.Book) (string, error)
	Post(shortBook models.ShortBook) (*models.Book, error)
	Delete(id string) error
	Search(query models.BookQuery) ([]models.Book, error)
	GetStore() (models.Store, error)
}

func GetDataStore() BookManipulator {
	return &datastore.ElasticDb{}
}

func GetCache() BookCache {
	return &cache.RedisCache{}
}

type BookCache interface {
	CacheGet(username string) ([]string, error)
	CacheAdd(username string, ctx *gin.Context)
}

// change name :(
