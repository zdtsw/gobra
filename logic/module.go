package main

import (
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

// module for "support"

///////////////////////////////////data strcture /////////////////////////////////////////
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
		TasksRunning int `json:"TasksRunning"`
	} `json:"apps"`
}

type returnAppResp struct {
	Host    string
	URL     string
	Project string
	Live    int
}

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


func log4Load() {
	//TODO
}

/////////////////////////////////////////////DCOS functions /////////////////////////////////////////////////
func queryDCOS(appid string) []byte {
	endpoint := "http://admin-thor.dice.se:8080/v2/apps?id=" + appid + "&label=HAPROXY_GROUP&embed=apps.count"
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
			image := "dreeu/docker-jenkins-oss:latest"
			if app.Container.Docker.Image == "registry.gitlab.ea.com/"+image {
				for _, pValue := range app.Container.Docker.Parameters {
					if (pValue.Key == "hostname") && (strings.Contains(pValue.Value, projShort)) && (app.Labels.IsTest != "true") {
						n = append(n, returnAppResp{Host: pValue.Value, URL: app.Labels.VHOST, Project: projShort, Live: app.TasksRunning})
						continue
					}
				}
			}
		case "bilbo":
			image := "dre-cobra/bilbo:latest"
			if app.Container.Docker.Image == "registry.gitlab.ea.com/"+image {
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
