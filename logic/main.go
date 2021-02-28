package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/static"
	"net/http"
	"html/template"
)

// global variable definitions
var (
		version string = "beta"
		author string = "WenZhou"
		project string
		release bool
		esIndex string = "bilbo"
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

	// register some functions
	r.SetFuncMap(template.FuncMap{
		"formatJSONResp":formatJSONResp,
    })
	
	r.LoadHTMLGlob("template/**/*.tmpl")
	r.Use(static.Serve("/img", static.LocalFile("./html/img", true)))
	// r.Use(static.Serve("/css", static.LocalFile("./html/css", true)))
	r.RedirectFixedPath = true
	r.RedirectTrailingSlash = true	

	r.GET("/", showIndexPage)

	bilbo := r.Group("/bilbo")
	{
		bilbo.GET("/", healthBilboHandler)
		bilbo.GET("/create", createBilboHandler)  
		bilbo.GET("/update", updateBilboHandler)
		bilbo.GET("/query", queryBilboHandler)
	}
	jenkins := r.Group("/jenkins")
	{
		jenkins.GET("/info", projectInfoHandler)
		jenkins.GET("/project/:proj", jenkinsInstanceHandler)
	}
	k8s := r.Group("/k8s")
	{
		k8s.GET("/list", listServiceHandler)
		k8s.GET("/action/:action", createServiceHandler)
	}

	r.Run(":8888")
}

