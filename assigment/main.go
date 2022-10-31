package main

import (
	"github.com/gin-gonic/gin"
)

const baseURL = "0.0.0.0:8000"

func main() {

	router := gin.Default()
	router.GET("/todo", HelloWorld)

	router.Run(baseURL)
}

func HelloWorld(c *gin.Context) {

	word := "Hello World"
	c.JSON(200, word)
}
