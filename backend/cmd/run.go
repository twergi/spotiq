package main

import (
	"net/http"

	"os"

	"github.com/gin-gonic/gin"
)

func pingGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func main() {
	host, ok := os.LookupEnv("HOST")

	if !ok {
		host = "127.0.0.1:8000"
	}

	r := gin.Default()
	r.GET("/ping", pingGet)
	r.Run(host)
}
