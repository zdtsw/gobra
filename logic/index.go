package logic

// module for "index page"

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ShowIndexPage godoc
// @Summary Show a list of services provided in Gobra
// @Description list jenkins, bilbo, k8s, aws and so on
// @ID ShowIndexPage
// @Tags main
// @Accept  json
// @Produce  html
// @Success 200 {string} todo
// @Router / [get]
func ShowIndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"version": version,
		"author":  author,
		"title":   "Welcome to the Love Of Team Wen",
	})
}
