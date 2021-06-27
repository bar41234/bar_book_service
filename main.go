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
	router := gin.Default()

	//Adding middlewares
	middlewares(router)

	//Registering the routes
	service.Routes(router)

	err = router.Run()
	if err != nil {
		panic(err)
	}
}

func middlewares(router *gin.Engine) {
	router.Use(service.Middleware)
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
