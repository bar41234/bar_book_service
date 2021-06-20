package main

import (
	"github.com/bar41234/bar_book_service/service"
	"github.com/gin-gonic/gin"
)

func main() {
	gin := gin.Default()
	service.Routes(gin)
	gin.Run()
}
