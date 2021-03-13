package main

import (
	"log"
	//"net/http"
	"github.com/gin-gonic/gin"
)

// module for "k8s"
func createServiceHandler(c *gin.Context) {
	log4Caller()
	log4Debug()
	log.Println("Load page in path: " + c.Request.URL.Path)
}

func listServiceHandler(c *gin.Context) {
	log4Caller()
	log4Debug()
	log.Println("Load page in path: " + c.Request.URL.Path)
}
