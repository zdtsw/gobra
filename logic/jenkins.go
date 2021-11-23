package main

// module for "jenkins"

import (
	"crypto/tls"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	dstURLFileHead   string = "https://gitlab.mycompany.com/api/v4/projects/1410/repository/files%2Fsrc%2Fproject/"
	dstURLFolderHead string = "https://gitlab.mycompany.com/api/v4/projects/1410/repository/tree?path=src/project/"
	dstURLTail       string = "/raw?ref=master"
	token            string = "Bearer txt1Jvx1ywo6LfZ3qndi"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
type dreProject struct {
	Project    string `json:"project"`
	Projshort  string `json:"short"`
	Studio     string `json:"studio"`
	Dashboard  string `json:"dashboard"`
	Googlepage string `json:"google"`
}
type gitlabResponseFolder []struct {
	//ID   string `json:"id"`  // useless field comment out
	// Mode string `json:"mode"` // useless field comment out
	Name string `json:"name"` // project shortname or filename
	Type string `json:"type"` // blob;folder
	// Path string `json:"path"`  // e.g src/project/dun/ useless comment out
}

////////////////////////////////////////////////////////////////////////////////////////////////////
var allProjects = []dreProject{
	{
		Project: "MyProject1", Projshort: "One", Studio: "StudioA, StudioB",
		Dashboard:  "https://wen-dashing.mycompany.com/Queue",
		Googlepage: "https://sites.google.com/mycompany.com/wen-myproject1/home",
	},
	{
		Project: "MyProject2", Projshort: "Two", Studio: "StudioA",
		Dashboard:  "https://wen-dashing.mycompany.com/Queue",
		Googlepage: "https://sites.google.com/mycompany.com/wen-myproject2/home",
	},
	{Project: "MyProject3", Projshort: "Tre", Studio: "StudioC"},
	{Project: "MyProject4", Projshort: "Si", Studio: "StudioC"},
	{
		Project: "MyProject5", Projshort: "Fem", Studio: "StudioB",
		Dashboard:  "https://excalibur-devblog.mycompany.com/dev/dashboard",
		Googlepage: "https://sites.google.com/mycompany.com/wen-myproject5/home",
	},
	{
		Project: "MyProject6", Projshort: "Six", Studio: "StudioA",
		Googlepage: "https://sites.google.com/mycompany.com/wen-myproject6/home",
	},
	{
		Project: "MyProject7", Projshort: "Sju", Studio: "StudioA",
		Googlepage: "https://sites.google.com/mycompany.com/wen-myproject7/home",
	},
	{Project: "MyProject8", Projshort: "Eight", Studio: "StudioB"},
}

////////////////////////////////////////////////////////////////////////////////////////////////////
func getAllProjects() []dreProject {
	return allProjects
}

/////////////////////////////////////////////GitLab functions///////////////////////////////////////////////////
func queryGitlab(path string, isFile bool) []byte {
	var endpoint string
	if isFile {
		endpoint = dstURLFileHead + url.PathEscape(path) + dstURLTail
	} else {
		endpoint = dstURLFolderHead + path
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		errorHandler(err)
	}
	req.Header.Add("Authorization", token)

	// resp, err := http.Get(endpoint); if we do not need token, this should be enough
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)

	if err != nil {
		errorHandler(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorHandler(err)
	}

	return body
}

//////////////////////////////////JSON functins/////////////////////////////////////
func parsegitlabJSONResponse(body []byte) *gitlabResponseFolder {
	var gr = new(gitlabResponseFolder)
	if err := json.Unmarshal(body, &gr); err != nil {
		errorHandler(err)
	}
	//	fmt.Println(&s)
	return gr
}

// func getJenkinsMasters(Projshort string) *gitlabResponseFolder {
// 	jsonBody := queryGitlab(Projshort+"/masterSettings", false)
// 	return parsegitlabJSONResponse(jsonBody)
// }

func getJenkinsBranches(Projshort string) *gitlabResponseFolder {
	jsonBody := queryGitlab(Projshort+"/branchSettings", false)
	return parsegitlabJSONResponse(jsonBody)
}

func getJenkinsMasterURL(Projshort string) []returnAppResp {
	jsonBody := queryDCOS("jenkins")
	return parseDCOSJSONResponse(jsonBody, "jenkins", Projshort)
}

/////////////////////////////template function//////////////////////////
/* convert return item from gitlab folder by strip file sufix and convert to lowercase*/
func convertFileJSONResp(n string) string {
	return strings.ToLower(strings.Split(n, ".")[0])
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////
// projectInfoHandler godoc
// @Summary Show a specific Jenkins controller info
// @Description get detail info. about controller
// @ID projectInfoHandler
// @Tags jenkins
// @Accept  json
// @Produce  json
// @Success 200 {string} todo
// @Router /jenkins/info [get]
func projectInfoHandler(c *gin.Context) {
	log4Caller()
	log4Debug()
	allProjects := getAllProjects()
	renderResponse(c, gin.H{
		"payload": allProjects, //on json or xml response can retrieve more info than html
		"version": render.VersionPage,
		"author":  render.ContactAuthor,
		"title":   "Jenkins Services",
	}, "jenkins/info.tmpl")
}


// jenkinsInstanceHandler godoc
// @Summary Show a specific Jenkins controller info
// @Description get detail info. about controller with branches
// @ID jenkinsInstanceHandler
// @Tags jenkins
// @Accept  json
// @Produce  json
// @Param proj path string true "Project Name"
// @Success 200 {string} todo
// @Router /jenkins/project/{proj} [get]
func jenkinsInstanceHandler(c *gin.Context) {
	log4Caller()
	log4Debug()
	log.Println("Load page in path: " + c.Request.URL.Path)
	projName := c.Param("proj")
	//allMasters := getJenkinsMasters(projName)
	allBranches := getJenkinsBranches(projName)
	allDCOS := getJenkinsMasterURL(projName)

	renderResponse(c, gin.H{
		//"payloadmaster": allMasters,
		"payloadbranch": allBranches,
		"payloaddcos":   allDCOS,
		"version":       render.VersionPage,
		"author":        render.ContactAuthor,
		"project":       projName,
		"title":         projName + " Jenkins",
	}, "jenkins/main.tmpl")
}
