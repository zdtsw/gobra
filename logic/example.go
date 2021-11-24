package logic

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// logic function definitions
func example() {
	r := gin.Default()
	r.Static("/assets", "./assets")
	r.StaticFS("/more_static", http.Dir("my_file_system"))
	r.StaticFile("/favicon.ico", "./resources/favicon.ico")

	r.GET("/project/:bilbo", func(c *gin.Context) {
		bilbo := c.Param("bilbo")
		c.String(http.StatusOK, "Using %s's bilbo", bilbo)
	})

	r.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)

	})

	// router.POST("/somePost", posting)
	// router.PUT("/somePut", putting)
	// router.DELETE("/someDelete", deleting)
	// router.PATCH("/somePatch", patching)
	// router.HEAD("/someHead", head)
	// router.OPTIONS("/someOptions", options)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}

func renderFormat(c *gin.Context, data gin.H, templateName string) {

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		c.XML(http.StatusOK, data["payload"])
	default:
		c.HTML(http.StatusOK, templateName, data)
	}

}

/*
<script>
function loadImage(){
    console.log("lololololo");
  //  element.setAttribute("src", "/img/redcheck.svg");


function changeImage(element, ClusterName)
{
    var x = document.getElementById(ClusterName);
    console.log(x);
    console.log(x.getAttribute("alt"));
   //  var y = document.getElementsByTagName("P");
   // console.log(y)
    if (x == "red" ) {
        v ="/img/redcheck.svg";
        x.setAttribute("src", v);
    }
}
</script>
*/
