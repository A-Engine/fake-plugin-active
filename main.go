package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"time"
	"github.com/A-Engine/fake-plugin-active/client"
)

var count = 0

func download(){
	wporg.UpdateCheck()
	wporg.DownloadPlugin()

	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan bool)
	for {
		select {
		case  <- ticker.C:
			wporg.UpdateCheck()
			wporg.DownloadPlugin()
			count++
		case <- quit:
			ticker.Stop()
			return
		}
	}
}

func main() {

	go download()

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"count": count,
		})
	})

	router.Run(":" + port)
}
