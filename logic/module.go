package main

import (
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// module for "support"

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
						n = append(n, returnAppResp{Host: pValue.Value, URL: app.Labels.VHOST, Project: projShort})
						continue
					}
				}
			}
		case "bilbo":
			image := "dre-cobra/bilbo:latest"
			if app.Container.Docker.Image == "registry.gitlab.ea.com/"+image {
				n = append(n, returnAppResp{Host: app.Env.NAME, URL: app.Labels.VHOST, Project: projShort})
				continue
			}
		}
	}
	//fmt.Println(n)
	return n
}
