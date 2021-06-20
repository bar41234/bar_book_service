package service

import (
	"github.com/bar41234/bar_book_service/models"
	"github.com/bar41234/bar_book_service/platform"
	"sync"
)

// Singleton
var ds platform.BookManipulator
var dsOnce sync.Once

// Singleton method
func GetDsInstance() platform.BookManipulator {
	dsOnce.Do(func() {
		ds = platform.GetDataStore()
	})
	return ds
}

func Get(id string) (*models.Book, error) {
	dataStore := GetDsInstance()
	return dataStore.Get(id)
}

func Put(book models.Book) (string, error) {
	dataStore := GetDsInstance()
	return dataStore.Put(book)
}

func Post(shortBook models.ShortBook) (*models.Book, error) {
	dataStore := GetDsInstance()
	return dataStore.Post(shortBook)
}

func Delete(id string) error {
	dataStore := GetDsInstance()
	return dataStore.Delete(id)
}

func Search(query models.BookQuery) ([]models.Book, error) {
	datastore := GetDsInstance()
	return datastore.Search(query)

}

func GetStore() (models.Store, error) {
	dataStore := GetDsInstance()
	return dataStore.GetStore()
}
