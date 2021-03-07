package main

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

// global variable definitions
var (
	version string = "beta"
	author  string = "Wen Zhou"
	project string
	release bool
	esIndex string = "bilbo"
	r       *gin.Engine
)

type pageFiller struct {
	VersionPage   string
	ContactAuthor string
	EAProject     string
}

var render = pageFiller{
	VersionPage:   version,
	ContactAuthor: author,
	EAProject:     project,
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

func errorHandler(err error) {
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
		"convertFileJSONResp": convertFileJSONResp,
		"showStatusIcon":      showStatusIcon,
	})

	r.LoadHTMLGlob("template/**/*.tmpl")
	r.Use(static.Serve("/img", static.LocalFile("./html/img", true)))
	// r.Use(static.Serve("/css", static.LocalFile("./html/css", true)))
	r.RedirectFixedPath = true
	r.RedirectTrailingSlash = true

	r.GET("/", showIndexPage)

	bilbo := r.Group("/bilbo")
	{
		bilbo.GET("/health", healthBilboHandler)
		bilbo.GET("/create/:proj", createBilboHandler)
		bilbo.GET("/update/:proj", updateBilboHandler)
		bilbo.GET("/query/:proj", queryBilboHandler)
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
