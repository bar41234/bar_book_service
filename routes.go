package main

import (
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) *gin.Engine {
	router.GET("/ping", Ping)
	book := router.Group("/book")
	{
		book.PUT("", PutBook)
		book.GET("/:id", GetBook)
		book.POST("/:id", PostBook)
		book.DELETE("/:id", DeleteBook)
	}
	router.GET("/search", SearchBook)
	router.GET("/store", GetStoreInfo)
	router.GET("/activity", Cache)

	return router
}
