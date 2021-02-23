package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// global variable definitions
var (
		version string = "Wen"
		author string = "Zhou"
		project string
		release bool
		r *gin.Engine
)

type pageFiller struct {
	VersionPage   string
	ContactAuthor string
	EAProject	string
}

var render = pageFiller{
	VersionPage:   version,
	ContactAuthor: author,
	EAProject: project,
}

func renderResponse(c *gin.Context, data gin.H, tmplFile string) {
	switch c.Request.Header.Get("Accept") {
		case "application/json":
		c.JSON(http.StatusOK, data["payload"])
		case "application/xml":
		c.XML(http.StatusOK, data["payload"])
		default:
		c.HTML(http.StatusOK, tmplFile, data)
		}  
}

func errorHandler(err error){
	//&(gin.Context).JSON(http.StatusBadRequest, gin.H{"Error: ": err})
	panic(err.Error())
}

// main function definition
func main() {

	if release { 
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	
	r.LoadHTMLGlob("template/**/*")
	r.RedirectFixedPath = true
	r.RedirectTrailingSlash = true	

	r.GET("/", showIndexPage)

	bilbo := r.Group("/bilbo")
	{
		bilbo.PUT("/create", createBilboHandler)
		bilbo.POST("/update", updateBilboHandler)
		bilbo.GET("/query", queryBilboHandler)
	}
	jenkins := r.Group("/jenkins")
	{
		jenkins.GET("/info", projectInfoHandler)
		jenkins.GET("/project/:proj", jenkinsInstanceHandler)
	}
	k8s := r.Group("/k8s")
	{
		k8s.GET("/:action", createServiceHandler)
	}

	r.Run(":8888")
}

