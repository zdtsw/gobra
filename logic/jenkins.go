package main
// module for "jenkins"

import (
	"log"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

const (
	dstURLFileHead string = "https://gitlab.ea.com/api/v4/projects/1410/repository/files%2Fsrc%2Fproject/"
	dstURLFolderHead string = "https://gitlab.ea.com/api/v4/projects/1410/repository/tree?path=src/project/"
	dstURLTail string = "/raw?ref=master"
	token string = "Bearer txt1Jvx1ywo6LfZ3qndi"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
type dreProject struct {
	Project  	string   	`json:"project"`
	Projshort	string 		`json:"short"`
	Studio		string		`json:"studio"`
}
type gitlabResponseFolder []struct {
	ID   string `json:"id"`  // useless field
	Mode string `json:"mode"` // useless field
	Name string `json:"name"` // project shortname or filename
	Type string `json:"type"`  // blob;folder
	Path string `json:"path"`  // e.g src/project/dun/
}

////////////////////////////////////////////////////////////////////////////////////////////////////
var allProjects = []dreProject{
	{Project: "Kingston", 			Projshort: "kin", 	Studio: "DICE"},
	{Project: "Walrus", 			Projshort: "wal",	Studio: "DICE"},
	{Project: "Casablanca", 		Projshort: "cas",	Studio: "DICE"},
	{Project: "DiceUpgradeNext", 	Projshort: "dun",	Studio: "DICE"},
	{Project: "FB1", 				Projshort: "fb1",	Studio: "Frostbite"},
	{Project: "FB2021", 			Projshort: "fb2021",Studio: "Frostbite"},
	{Project: "Excalibur", 			Projshort: "exc", 	Studio: "Critiron"},
}

////////////////////////////////////////////////////////////////////////////////////////////////////
func getAllProjects() []dreProject {
	return allProjects
}

func queryGitlab(path string, isFile bool) ([]byte) {
	var endpoint string
	if isFile {
		endpoint = dstURLFileHead+url.PathEscape(path)+dstURLTail
	} else {
		endpoint = dstURLFolderHead+path
	}
	
	log.Println("WEN--query: " + endpoint)

	req, err := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Authorization", token)

	// resp, err := http.Get(endpoint);
	client := &http.Client{}
    resp, err := client.Do(req)
    
	if err != nil {
		errorHandler(err)
	}
	defer resp.Body.Close()

   	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorHandler(err)
	}
	log.Println("WEN--RESPONSE:" + string([]byte(body)))
	return body
}

func parseJSONResponse(body []byte)(*gitlabResponseFolder) {
	log.Println("WEN--running: parseJSONResponse()")
	inputs := string(body)
	fmt.Println(inputs)
    var s = new(gitlabResponseFolder)
    if err := json.Unmarshal(body, &s); err != nil {
		errorHandler(err)
	}
	fmt.Println(&s)
	// log.Println("WEN--getting:" + s)
    return s
}

func getJenkinsMasters(Projshort string) (*gitlabResponseFolder) {
	body := queryGitlab(Projshort+"/masterSettings", false)
	return parseJSONResponse(body)
}

// routine functions for /jenkins/info
func projectInfoHandler(c *gin.Context) {
	log.Println("Calling: jenkinsInfoHandler")
	log.Println("Load page in path: " + c.Request.URL.Path)
	allProjects := getAllProjects()
	renderResponse(c, gin.H{
					"payload": allProjects,  //on json or xml response can retrieve more info than html 
					"version": render.VersionPage,
					"author": render.ContactAuthor,
	  				},"jenkins/info.tmpl")	
}

// routine functions for /jenkins/project/:proj
func jenkinsInstanceHandler(c *gin.Context) {
	log.Println("Calling: jenkinsMainHandler")
	log.Println("Load page in path: " + c.Request.URL.Path)
	projName := c.Param("proj")
	allMasters := getJenkinsMasters(projName)
	//if (gitlabResponseFolder{}) == allMasters  {log.Println("WEN-DEBIG: allMasters is empty, damn it ")}

	c.HTML(http.StatusOK, "jenkins/main.tmpl", gin.H{
		"version": render.VersionPage,
		"author":  render.ContactAuthor,
		"project": projName,
		"payload": allMasters,
	})
}
