package main
// module for "index page"

import (
	//log"
	"net/http"
	"github.com/gin-gonic/gin"
)

func showIndexPage(c *gin.Context){
	//log.Println("HELLO INDEX")
	c.HTML(http.StatusOK, "index.html", gin.H {
		"version": version,
		"author": author,
		"title": "Welcome to the Love Of Gobra",
	})
}