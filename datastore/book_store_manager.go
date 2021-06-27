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
var bookStoreManager BookStorer
var userActivityManager UserActivitier
var bookInitOnce sync.Once
var userInitOnce sync.Once
var errBookInit error
var errUserInit error

func BooksStoreFactory() (BookStorer, error) {
	bookInitOnce.Do(func() {
		bookStoreManager, errBookInit = NewElasticBookManager(elasticUrl)
	})
	return bookStoreManager, errBookInit
}

func UserActivityFactory() (UserActivitier, error) {
	userInitOnce.Do(func() {
		userActivityManager, errUserInit = NewRedisUserActivity(redisUrl)
	})
	return userActivityManager, errUserInit
}

type BookStorer interface {
	GetBook(id string) (*models.Book, error)
	AddBook(book models.Book) (*string, error)
	UpdateBook(id string, title string) (*string, error)
	DeleteBook(id string) error
	Search(title string, author string, priceRange string) ([]models.Book, error)
	GetStoreInfo() (*models.Store, error)
}

type UserActivitier interface {
	GetActivities(username string) ([]string, error)
	AddActivity(username string, methodName string, request string)
}
