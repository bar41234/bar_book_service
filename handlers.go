package main

import (
	"encoding/json"
	"github.com/bar41234/bar_book_service/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ERROR_MSG_ID_NOT_FOUND         = "Error: Id was not found!"
	ERROR_MSG_INVALID_PUT_REQUEST  = "Error: Invalid PUT request"
	ERROR_MSG_INVALID_POST_REQUEST = "Error: Invalid POST request"
)

func Ping(c *gin.Context) {
	username := c.Query("username")
	CacheAdd(username, c)

	c.JSON(http.StatusOK, map[string]string{
		"PING!": "PONG!",
	})
}

func GetBook(c *gin.Context) {
	username := c.Query("username")
	CacheAdd(username, c)

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusNotFound, ERROR_MSG_ID_NOT_FOUND)
		return
	}
	book, err := Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, book)
}

func PutBook(c *gin.Context) {
	username := c.Query("username")
	CacheAdd(username, c)

	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, ERROR_MSG_INVALID_PUT_REQUEST)
		return
	}
	book := models.Book{}
	unmarshallErr := json.Unmarshal(jsonData, &book)
	if unmarshallErr != nil {
		c.JSON(http.StatusBadRequest, ERROR_MSG_INVALID_PUT_REQUEST)
		return
	}
	id, putErr := Put(book)
	if putErr != nil {
		c.JSON(http.StatusBadRequest, putErr)
		return
	}
	c.String(http.StatusOK, id)
}

func PostBook(c *gin.Context) {
	username := c.Query("username")
	CacheAdd(username, c)

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusNotFound, ERROR_MSG_ID_NOT_FOUND)
		return
	}
	jsonData, err := c.GetRawData()
	// Note: if the user post with typo title (e.g "titlee"), the post will succeed with title = ""
	if err != nil {
		c.JSON(http.StatusBadRequest, ERROR_MSG_INVALID_POST_REQUEST)
		return
	}
	shortBook := models.ShortBook{}
	unmarshalErr := json.Unmarshal(jsonData, &shortBook)
	if unmarshalErr != nil {
		c.JSON(http.StatusBadRequest, ERROR_MSG_INVALID_POST_REQUEST)
		return
	}
	shortBook.Id = id
	book, postErr := Post(shortBook)
	if postErr != nil {
		c.JSON(http.StatusBadRequest, postErr)
		return
	}
	c.JSON(http.StatusOK, book)
}

func DeleteBook(c *gin.Context) {
	username := c.Query("username")
	CacheAdd(username, c)

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusNotFound, ERROR_MSG_ID_NOT_FOUND)
		return
	}
	err := Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, "Book was successfully deleted!")
}

func SearchBook(c *gin.Context) {
	username := c.Query("username")
	CacheAdd(username, c)

	title := c.Query("title")
	author := c.Query("author_name")
	priceRange := c.Query("price_range")
	bookQuery := models.BookQuery{
		Title:      title,
		AuthorName: author,
		PriceRange: priceRange,
	}
	books, searchErr := Search(bookQuery)
	if searchErr != nil {
		c.JSON(http.StatusBadRequest, searchErr)
		return
	}
	c.JSON(http.StatusOK, books)
}

func GetStoreInfo(c *gin.Context) {
	username := c.Query("username")
	CacheAdd(username, c)

	store, err := GetStore()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, store)
}

func Cache(c *gin.Context) {
	username := c.Query("username")
	actions, errCache := CacheGet(username)
	if errCache != nil {
		c.JSON(http.StatusBadRequest, errCache)
		return
	}
	c.JSON(http.StatusOK, actions)
}
