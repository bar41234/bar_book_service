package service

import (
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) *gin.Engine {
	router.Use(middleware)
	router.GET("/ping", Ping)
	book := router.Group("/book")
	{
		book.PUT("", AddBook)
		book.GET("/:id", GetBook)
		book.POST("/:id", UpdateBook)
		book.DELETE("/:id", DeleteBook)
	}
	router.GET("/search", SearchBook)
	router.GET("/store", GetStoreInfo)
	router.GET("/activity", GetActivities)

	return router
}
