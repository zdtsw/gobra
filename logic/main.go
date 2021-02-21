package main

import (
	"github.com/gin-gonic/gin"
)

// global variable definitions
var (
		version string = "Wen"
		author string = "Zhou"
		project string
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

// main function definition
func main() {

	r := gin.Default()
	r.LoadHTMLGlob("template/**/*")
	r.RedirectFixedPath = true
	r.RedirectTrailingSlash = true
	

//	initializeRoutes()
	r.GET("/", showIndexPage)

	bilbo := r.Group("/bilbo")
	{
		bilbo.PUT("/create", createBilboHandler)
		bilbo.POST("/update", updateBilboHandler)
		bilbo.GET("/query", queryBilboHandler)
	}
	jenkins := r.Group("/jenkins")
	{
		jenkins.GET("/info", jenkinsInfoHandler)
		jenkins.GET("/project/:project", jenkinsMainHandler)
	}
	k8s := r.Group("/k8s")
	{
		k8s.GET("/:action", createServiceHandler)
	}

	r.Run(":8888")
}

