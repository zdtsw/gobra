package main

// module for "jenkins"

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	//"fmt"
	"strings"
)

const (
	dstURLFileHead   string = "https://gitlab.ea.com/api/v4/projects/1410/repository/files%2Fsrc%2Fproject/"
	dstURLFolderHead string = "https://gitlab.ea.com/api/v4/projects/1410/repository/tree?path=src/project/"
	dstURLTail       string = "/raw?ref=master"
	token            string = "Bearer txt1Jvx1ywo6LfZ3qndi"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
type dreProject struct {
	Project   string `json:"project"`
	Projshort string `json:"short"`
	Studio    string `json:"studio"`
}
type gitlabResponseFolder []struct {
	//ID   string `json:"id"`  // useless field comment out
	// Mode string `json:"mode"` // useless field comment out
	Name string `json:"name"` // project shortname or filename
	Type string `json:"type"` // blob;folder
	// Path string `json:"path"`  // e.g src/project/dun/ useless comment out
}

type appResponse struct {
	Apps []struct {
		ID        string `json:"id"`
		Container struct {
			Docker struct {
				Image      string `json:"image"`
				Parameters []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"parameters"`
			} `json:"docker"`
		} `json:"container"`
		Env struct {
			JAVAOPTS        string `json:"JAVA_OPTS"`
			PROJECTSEEDROOT string `json:"PROJECT_SEED_ROOT"`
			NAME            string `json:"CLUSTER_NAME"`
		} `json:"env"`
		Labels struct {
			SEEDROOT string `json:"SEED_ROOT"`
			IsTest   string `json:"IS_TEST_INSTANCE"`
			VHOST    string `json:"HAPROXY_0_VHOST"`
		} `json:"labels"`
	} `json:"apps"`
}

type returnAppResp struct {
	Host    string
	URL     string
	Project string
}

////////////////////////////////////////////////////////////////////////////////////////////////////
var allProjects = []dreProject{
	{Project: "Kingston", Projshort: "kin", Studio: "DICE"},
	{Project: "Walrus", Projshort: "wal", Studio: "DICE"},
	{Project: "Casablanca", Projshort: "cas", Studio: "DICE"},
	{Project: "DiceUpgradeNext", Projshort: "dun", Studio: "DICE"},
	{Project: "FB1", Projshort: "fb1", Studio: "Frostbite"},
	{Project: "FB2021", Projshort: "fb2021", Studio: "Frostbite"},
	{Project: "Excalibur", Projshort: "exc", Studio: "Critiron"},
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

	// log.Println("WEN--query: " + endpoint)

	req, err := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Authorization", token)

	// resp, err := http.Get(endpoint); if we do not need token, this should be enough
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
	//	log.Println("WEN--RESPONSE:" + string([]byte(body)))
	return body
}

//////////////////////////////////JSON functins/////////////////////////////////////
func parsegitlabJSONResponse(body []byte) *gitlabResponseFolder {
	//	log.Println("WEN--running: parseJSONResponse()")
	//fmt.Println(string(body))
	var gr = new(gitlabResponseFolder)
	if err := json.Unmarshal(body, &gr); err != nil {
		errorHandler(err)
	}
	//	fmt.Println(&s)
	return gr
}

func getJenkinsMasters(Projshort string) *gitlabResponseFolder {
	jsonBody := queryGitlab(Projshort+"/masterSettings", false)
	return parsegitlabJSONResponse(jsonBody)
}

func getJenkinsBranches(Projshort string) *gitlabResponseFolder {
	jsonBody := queryGitlab(Projshort+"/branchSettings", false)
	return parsegitlabJSONResponse(jsonBody)
}

func getJenkinsMasterURL(Projshort string) []returnAppResp {
	jsonBody := queryDCOS("jenkins")
	//fmt.Println(string(jsonBody))
	return parseDCOSJSONResponse(jsonBody, "jenkins", Projshort)
}

/////////////////////////////template function//////////////////////////
/* convert return item from gitlab folder by strip file sufix and convert to lowercase*/
func convertFileJSONResp(n string) string {
	return strings.ToLower(strings.Split(n, ".")[0])
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////
// routine functions for /jenkins/info
func projectInfoHandler(c *gin.Context) {
	log.Println("Calling: jenkinsInfoHandler")
	log.Println("Load page in path: " + c.Request.URL.Path)
	allProjects := getAllProjects()
	renderResponse(c, gin.H{
		"payload": allProjects, //on json or xml response can retrieve more info than html
		"version": render.VersionPage,
		"author":  render.ContactAuthor,
	}, "jenkins/info.tmpl")
}

// routine functions for /jenkins/project/:proj
func jenkinsInstanceHandler(c *gin.Context) {
	log.Println("Calling: jenkinsMainHandler")
	log.Println("Load page in path: " + c.Request.URL.Path)
	projName := c.Param("proj")
	//allMasters := getJenkinsMasters(projName)
	allBranches := getJenkinsBranches(projName)
	allDCOS := getJenkinsMasterURL(projName)
	//if (gitlabResponseFolder{}) == allMasters  {log.Println("WEN-DEBIG: allMasters is empty, damn it ")}

	// c.HTML(http.StatusOK, "jenkins/main.tmpl", gin.H{
	// 	"version": render.VersionPage,
	// 	"author":  render.ContactAuthor,
	// 	"project": projName,
	// 	"payloadmaster": allMasters,
	// 	"payloadbranch": allBranches,
	//})
	renderResponse(c, gin.H{
		//"payloadmaster": allMasters,
		"payloadbranch": allBranches,
		"payloaddcos":   allDCOS,
		"version":       render.VersionPage,
		"author":        render.ContactAuthor,
		"project":       projName,
	}, "jenkins/main.tmpl")
}
