package main

import (
	"log"
	"net/http"

	"os"

	"github.com/gin-gonic/gin"
)

var albums []SongData = make([]SongData, 0)

func Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"device": gin.H{"ID": "ID"}, "token": gin.H{"Token": "Token"}})
}

func Login(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "https://spotify.com")
}

func Logout(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func GetCurrentQueue(c *gin.Context) {
	c.JSON(http.StatusOK, albums)
}

type SongData struct {
	Name string `json:"name"`
}

func AddToCurrentQueue(c *gin.Context) {
	var song SongData

	if err := c.BindJSON(&song); err != nil {
		log.Println(err)
		return
	}

	albums = append(albums, song)

	c.Status(http.StatusOK)
}

func main() {
	host, ok := os.LookupEnv("HOST")

	if !ok {
		host = "127.0.0.1:8000"
	}

	r := gin.Default()
	r.GET("/", Root)
	r.GET("/login/", Login)
	r.GET("/logout/", Logout)
	r.GET("/current_queue/", GetCurrentQueue)
	r.POST("/current_queue/", AddToCurrentQueue)

	r.Run(host)
}
