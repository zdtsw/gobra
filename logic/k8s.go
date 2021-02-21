package main
import (
	"log"
	//"net/http"
	"github.com/gin-gonic/gin"
)

// functions for "k8s"
func createServiceHandler(c *gin.Context) {
	log.Println("Calling: createServiceHandler")
	log.Println("Load page in path: " + c.Request.URL.Path)
}