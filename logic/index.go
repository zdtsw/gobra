package main

// module for "index page"

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func showIndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"version": version,
		"author":  author,
		"title":   "Welcome to the Love Of Team Wen",
	})
}
