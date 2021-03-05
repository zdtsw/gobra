package main

// module for "index page"

import (
	//log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func showIndexPage(c *gin.Context) {
	//log.Println("HELLO INDEX")
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"version": version,
		"author":  author,
		"title":   "Welcome to the Love Of Gobra",
	})
}
