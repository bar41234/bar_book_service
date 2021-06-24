package service

import (
	"encoding/json"
	"fmt"
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
	bookStore, _ := datastore.BooksStoreFactory()
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusNotFound, errorMsgIdNotFound)
		return
	}
	book, err := bookStore.GetBook(id)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, book)
}

func AddBook(c *gin.Context) {
	bookStore, _ := datastore.BooksStoreFactory()
	book := models.Book{}
	err := c.ShouldBind(&book)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorMsgInvalidPutRequest)
		return
	}
	id, err := bookStore.AddBook(book)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, id)
}

func UpdateBook(c *gin.Context) {
	bookStore, _ := datastore.BooksStoreFactory()
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
	bookId, err := bookStore.UpdateBook(id, book.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, fmt.Sprintf("Book %s was successfully updated", bookId))
}

func DeleteBook(c *gin.Context) {
	bookStore, _ := datastore.BooksStoreFactory()
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusNotFound, errorMsgIdNotFound)
		return
	}
	err := bookStore.DeleteBook(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, "Book was successfully deleted!")
}

func SearchBook(c *gin.Context) {
	bookStore, _ := datastore.BooksStoreFactory()
	title := c.Query("title")
	author := c.Query("author_name")
	priceRange := c.Query("price_range")
	books, err := bookStore.Search(title, author, priceRange)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, books)
}

func GetStoreInfo(c *gin.Context) {
	bookStore, _ := datastore.BooksStoreFactory()
	store, err := bookStore.GetStoreInfo()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, store)
}

func GetActivities(c *gin.Context) {
	userActivity, _ := datastore.UserActivityFactory()
	username := c.Query("username")
	actions, err := userActivity.GetActivities(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, actions)
}

func middleware(c *gin.Context) {
	userActivity, _ := datastore.UserActivityFactory()
	username := c.Query("username")
	if username == "" {
		return
	}
	userActivity.AddActivity(username, c.Request.Method, c.Request.RequestURI)
	c.Next()
}
