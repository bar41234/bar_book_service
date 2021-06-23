package datastore

import (
	"github.com/bar41234/bar_book_service/models"
	"sync"
)

const (
	elasticUrl = "http://es-search-7.fiverrdev.com:9200"
	redisUrl   = "localhost:6379"
)

// Singleton
var bookStoreManager BookStoreManager
var userActivityManager UserActivityManager
var bcfOnce sync.Once
var uafOnce sync.Once

// Singleton method
func BooksContainerFactory() (BookStoreManager, error) {
	var err error
	bcfOnce.Do(func() {
		bookStoreManager, err = NewElastic(elasticUrl)
	})
	return bookStoreManager, err
}

// Singleton method
func UserActivityFactory() (UserActivityManager, error) {
	var err error
	uafOnce.Do(func() {
		userActivityManager, err = NewRedis(redisUrl)
	})
	return userActivityManager, err
}

type BookStoreManager interface {
	Get(id string) (*models.Book, error)
	Add(book models.Book) (string, error)
	Update(id string, title string) (string, error)
	Delete(id string) error
	Search(title string, author string, priceRange string) ([]models.Book, error)
	GetStore() (*models.Store, error)
}

type UserActivityManager interface {
	Get(username string) ([]string, error)
	Add(username string, methodName string, request string)
}
