package main
import (
	"log"
	//"net/http"
	"github.com/gin-gonic/gin"
)



// functions for "bilbo"
func createBilboHandler (c *gin.Context) {
	log.Println("Calling: createBilboHandler")
	log.Println("Load page in path: " + c.Request.URL.Path)
}
func updateBilboHandler(c *gin.Context) {
	log.Println("Calling: updateBilboHandler")
	log.Println("Load page in path: " + c.Request.URL.Path)
}
func queryBilboHandler(c *gin.Context) {
	log.Println("Calling: queryBilboHandler")
	log.Println("Load page in path: " + c.Request.URL.Path)
}