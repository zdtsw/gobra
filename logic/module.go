package logic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"
)

/////////////////////////////////////////////LOGGIN functions ///////////////////////////////////////////////
func log4Debug() {
	_, filename, lineNum, stat := runtime.Caller(0)
	if stat {
		log.Println("Calling: " + filename + " line: " + strconv.Itoa(lineNum))
	}
}

func log4Caller() {
	pc, _, _, stat := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if stat && details != nil {
		log.Println("Called from: " + details.Name())
	}
}

// func log4Load() {
// 	//TODO
// }

/////////////////////////////////////////////FORMAT functions ///////////////////////////////////////////////
//print out format "space" number of space then input string list loop
func spaceOutputs(space int, items ...string) {
	for _, item := range items {
		fmt.Println(strings.Repeat(" ", space) + item)
	}
}

/////////////////////////////////////////////DCOS functions /////////////////////////////////////////////////
func queryDCOS(appid string) []byte {
	endpoint := "http://admin-thor.mycompany.com:8080/v2/apps?id=" + appid + "&label=HAPROXY_GROUP&embed=apps.count"
	resp, err := http.Get(endpoint)
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

func parseDCOSJSONResponse(body []byte, instance string, projShort string) []returnAppResp {
	var ar appResponse
	var n []returnAppResp
	if err := json.Unmarshal(body, &ar); err != nil {
		errorHandler(err)
	}
	for _, app := range ar.Apps {
		// var image string
		switch instance {
		case "jenkins":
			image := "eu/docker-jenkins-oss:latest"
			if app.Container.Docker.Image == "registry.gitlab.mycompany.com/"+image {
				for _, pValue := range app.Container.Docker.Parameters {
					if (pValue.Key == "hostname") && (strings.Contains(pValue.Value, projShort)) && (app.Labels.IsTest != "true") {
						n = append(n, returnAppResp{Host: pValue.Value, URL: app.Labels.VHOST, Project: projShort, Live: app.TasksRunning})
						continue
					}
				}
			}
		case "bilbo":
			image := "wen/bilbo:latest"
			if app.Container.Docker.Image == "registry.gitlab.mycompany.com/"+image {
				n = append(n, returnAppResp{Host: app.Env.NAME, URL: app.Labels.VHOST, Project: projShort, Live: app.TasksRunning})
				continue
			}
		case "metrics":
			if app.Container.Docker.Image == "docker.elastic.co/elasticsearch/elasticsearch:6.2.2" {
				n = append(n, returnAppResp{Host: app.Env.NAME, URL: app.Labels.VHOST, Project: projShort, Live: app.TasksRunning})
				continue
			}
		}

	}
	//fmt.Println(n)
	return n
}


///////////////////////error handling////////////////////////
func errorHandler(err error) {
	//&(gin.Context).JSON(http.StatusBadRequest, gin.H{"Error: ": err})
	panic(err.Error())
}

//////////////////////common rednering//////////////////////////////
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
