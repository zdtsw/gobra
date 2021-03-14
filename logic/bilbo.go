package main

import (
	"log"
	"net/http"

	//"net/url"
	"encoding/json"
	"io/ioutil"

	es6 "github.com/elastic/go-elasticsearch/v6"
	"github.com/gin-gonic/gin"
	//es6api "github.com/elastic/go-elasticsearch/v6/esapi"
	// "fmt"
)

var bilbofqdn = "bilbo-fb1.dre.dice.se"
var bilboport = ":9200"
var escfg = es6.Config{
	Addresses: []string{"https://" + bilbofqdn + bilboport},
}
var docTypeList = [...]string{
	"code",
	"tnt_local",
	"ant_local",
	"drone",
	"symbols",
	"avalanchestate",
	"avalanchestate_autotest",
	"webexport",
	"frosty",
	"expressiondebugdata",
	"bundles",
	"clone_db",
	"build",
}

//////////////////////////////////////////////////////////////////////////////////
type healthResp struct {
	ClusterName       string `json:"cluster_name"`
	URL               string
	Status            string `json:"status"`
	Timedout          bool   `json:"timed_out"`
	NumberOfNodes     int    `json:"number_of_nodes"`
	NumberOfDataNodes int    `json:"number_of_data_nodes"`
}

type bilboOpsDropDown struct {
	value string
	name  string
}

////////////////////template functions/////////////////////////////
func showStatusIcon(status string) string {
	switch status {
	case "red":
		return "/img/redcross.svg"
	case "yellow":
		return "/img/yellowsign.svg"
	case "green":
		return "/img/greencheck.svg"
	default:
		return "/img/error.svg"
	}
}

// module for "bilbo"
func createBilboHandler(c *gin.Context) {
	log4Caller()
	log4Debug()
	log.Println("Load page in path: " + c.Request.URL.Path)
	// if ( myindex == "" ) {		myindex = esIndex	}
	// curl -XPUT  --header 'Content-Type: application/json' http://localhost:9200/myindex/_doc/1 -d '{
	// 	"school" : "Harvard"
	// }' add data json format
	// curl -u elastic:(password) -XPOST --header 'Content-Type: application/json' http://localhost:9200/myindex/_doc/1/_update -d '{
	// 	"doc" : {
	// 		"students": 50000
	// 	}
	// }' update data
	// curl -u elastic:(password) -XPOST --header 'Content-Type: application/json' http://localhost:9200/_reindex -d '{
	// 	"source": {
	// 	  "index": "myindex"
	// 	},
	// 	"dest": {
	// 	  "index": "myindex_backup"
	// 	}
	//   }' backup index

}

func updateBilboHandler(c *gin.Context) {
	log4Caller()
	log4Debug()
	log.Println("Load page in path: " + c.Request.URL.Path)
	//if ( myindex == "" ) {		myindex = esIndex	}
	//curl -u elastic:(password) -X DELETE 'http://localhost:9200/myindex'  delete index
}

func loadQueryPageHandler(c *gin.Context) {
	log4Caller()
	log4Debug()
	log.Println("Load page in path: " + c.Request.URL.Path)
}

func queryBilboHandler(c *gin.Context) {

	log4Caller()
	log4Debug()
	log.Println("Load page in path: " + c.Request.URL.Path)

	var es = new(es6.Client)
	es = initES()
	res, _ := es.Search(
		//es.Search.WithContext(context.Background()),
		es.Search.WithIndex("bilbo"),
		//es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	defer res.Body.Close()
	log.Println(res)
	// if ( myindex == "" ) {		myindex = esIndex	}
	// curl -u elastic:(password) -XGET 'http://localhost:9200/_cat/indices?pretty' list all indices
	// curl -u elastic:(password) -XGET 'http://localhost:9200/myindex/_search?pretty' list all documents in index
	// curl -u elastic:(password) -XGET http://localhost:9200/myindex list index mapping
	// curl -u elastic:(password) -XGET http://localhost:9200/myindex/_search?q=school:Harvard query when school=Harvard
	// http://bilbo-kin-uk.dre.dice.se/bilbo/_search?q=type:drone
	// curl -u elastic:(password) -XGET --header 'Content-Type: application/json' http://localhost:9200/myindex/_search -d '{
	// 	"query" : {
	// 		"match" : { "school": "Harvard"
	// 		}
	// 	}
	// }'  still query but with json

	// GET filebeat-7.6.2-2020.05.05-000001/_search
	// {
	// 	"_source": ["suricata.eve.timestamp","source.geo.region_name","event.created"],
	// 	"query":      {
	// 		"match" : { "source.geo.country_iso_code": "GR" }
	// 	}
	// }
	// GET filebeat-7.6.2-2020.05.05-000001/_search
	// {
	// 	"query": {
	// 		"range" : {
	// 			"event.created": {
	// 				"gte" : "now-7d/d"
	// 			}
	// 		}
	// 	}
	// }
}

func healthBilboHandler(c *gin.Context) {
	log4Caller()
	log4Debug()
	log.Println("Load page in path: " + c.Request.URL.Path)

	jsonBodyBilbo := queryDCOS("bilbo/") // to exclude bilbo-ui
	jsonBodyMetric := queryDCOS("es-metrics-1")
	//Projshort := "kin"
	result := parseDCOSJSONResponse(jsonBodyBilbo, "bilbo", "")
	result = append(result, parseDCOSJSONResponse(jsonBodyMetric, "metrics", "")...)

	var summary []healthResp
	var healthy healthResp
	for _, bilboInstance := range result {
		endpoint := "https://" + bilboInstance.URL + "/_cluster/health?pretty"
		if bilboInstance.Live == 0 {
			continue
		}
		resp, err := http.Get(endpoint)
		if err != nil {
			errorHandler(err)
		}
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			errorHandler(err)
		} else {
			if err := json.Unmarshal(body, &healthy); err != nil {
				errorHandler(err)
			}
			summary = append(summary, healthResp{ClusterName: healthy.ClusterName, URL: "https://" + bilboInstance.URL, Status: healthy.Status, NumberOfNodes: healthy.NumberOfNodes})
		}
	}
	//fmt.Println(summary)
	renderResponse(c, gin.H{
		"bilboSummary": summary,
		"version":      render.VersionPage,
		"author":       render.ContactAuthor,
		"title":        "Bilbo Services",
		"docTypeList":  docTypeList,
	}, "bilbo/health.tmpl")
}

func initES() *es6.Client {
	log4Debug()
	log4Caller()
	es, err := es6.NewDefaultClient()
	log.Println(es6.Version)

	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	// res, err := es.Info()
	// if err != nil {
	// log.Fatalf("Error getting response: %s", err)
	// }

	// defer res.Body.Close()
	// if res.IsError() {
	// 	log.Fatalf("Error: %s", res.String())
	// }
	// log.Println(res)

	return es
}
