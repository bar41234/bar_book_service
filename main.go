package main

import (
	"github.com/gin-gonic/gin"
)

func main()  {
	gin := gin.Default()
	Routes(gin)
	gin.Run()
}