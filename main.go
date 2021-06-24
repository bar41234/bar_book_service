package main

import (
	"github.com/bar41234/bar_book_service/datastore"
	"github.com/bar41234/bar_book_service/service"
	"github.com/gin-gonic/gin"
)

func main() {
	err := setup()
	if err != nil {
		panic(err)
	}
	gin := gin.Default()
	service.Routes(gin)
	gin.Run()
}

func setup() error {
	_, err := datastore.BooksStoreFactory()
	if err != nil {
		return err
	}
	_, err = datastore.UserActivityFactory()
	if err != nil {
		return err
	}
	return nil
}
