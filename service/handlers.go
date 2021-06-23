package service

import (
	"encoding/json"
	"github.com/bar41234/bar_book_service/datastore"
	"github.com/bar41234/bar_book_service/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	errorMsgIdNotFound         = "Error: Id was not found!"
	errorMsgInvalidPutRequest  = "Error: Invalid PUT request"
	errorMsgInvalidPostRequest = "Error: Invalid POST request"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{
		"message": "Ping Pong!",
	})
}

func GetBook(c *gin.Context) {
	bookStoreContainer, _ := datastore.BooksContainerFactory()
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusNotFound, errorMsgIdNotFound)
		return
	}
	book, err := bookStoreContainer.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, book)
}

func AddBook(c *gin.Context) {
	bookStoreContainer, _ := datastore.BooksContainerFactory()
	book := models.Book{}
	err := c.ShouldBind(&book)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorMsgInvalidPutRequest)
		return
	}
	id, err := bookStoreContainer.Add(book)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.String(http.StatusOK, id)
}

func UpdateBook(c *gin.Context) {
	bookStoreContainer, _ := datastore.BooksContainerFactory()
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusNotFound, errorMsgIdNotFound)
		return
	}
	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, errorMsgInvalidPostRequest)
		return
	}
	var book models.Book
	err = json.Unmarshal(jsonData, &book)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorMsgInvalidPostRequest)
		return
	}
	bookId, err := bookStoreContainer.Update(id, book.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, "Book "+bookId+" was successfully updated")
}

func DeleteBook(c *gin.Context) {
	bookStoreContainer, _ := datastore.BooksContainerFactory()
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusNotFound, errorMsgIdNotFound)
		return
	}
	err := bookStoreContainer.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, "Book was successfully deleted!")
}

func SearchBook(c *gin.Context) {
	bookStoreContainer, _ := datastore.BooksContainerFactory()
	title := c.Query("title")
	author := c.Query("author_name")
	priceRange := c.Query("price_range")
	books, err := bookStoreContainer.Search(title, author, priceRange)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, books)
}

func GetStoreInfo(c *gin.Context) {
	bookStoreContainer, _ := datastore.BooksContainerFactory()
	store, err := bookStoreContainer.GetStore()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, store)
}

func GetActivities(c *gin.Context) {
	userActivityManager, _ := datastore.UserActivityFactory()
	username := c.Query("username")
	actions, err := userActivityManager.Get(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, actions)
}

func middleware(c *gin.Context) {
	userActivityManager, _ := datastore.UserActivityFactory()
	username := c.Query("username")
	userActivityManager.Add(username, c.Request.Method, c.Request.RequestURI)
	c.Next()
}
