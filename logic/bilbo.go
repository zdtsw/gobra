package main

import (
	"log"
	"net/http"
	//"net/url"
	"encoding/json"
	es6 "github.com/elastic/go-elasticsearch/v6"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	//es6api "github.com/elastic/go-elasticsearch/v6/esapi"
)

var bilbofqdn = "bilbo-fb1.dre.dice.se"
var bilboport = ":9200"
var escfg = es6.Config{
	Addresses: []string{"https://" + bilbofqdn + bilboport},
}

//////////////////////////////////////////////////////////////////////////////////
type healthResp struct {
	ClusterName       string `json:"cluster_name"`
	Status            string `json:"status"`
	Timedout          bool   `json:"timed_out"`
	NumberOfNodes     int    `json:"number_of_nodes"`
	NumberOfDataNodes int    `json:"number_of_data_nodes"`
}

// module for "bilbo"
func createBilboHandler(c *gin.Context) {
	log.Println("Calling: createBilboHandler")
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
	log.Println("Calling: updateBilboHandler")
	log.Println("Load page in path: " + c.Request.URL.Path)
	//if ( myindex == "" ) {		myindex = esIndex	}
	//curl -u elastic:(password) -X DELETE 'http://localhost:9200/myindex'  delete index
}

func queryBilboHandler(c *gin.Context) {
	log.Println("Calling: queryBilboHandler")
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
	log.Println("Calling: healthBilboHandler")
	log.Println("Load page in path: " + c.Request.URL.Path)
	projName := "kin"
	endpoint := "http://bilbo-" + projName + ".dre.dice.se/_cluster/health?pretty"
	resp, err := http.Get(endpoint)
	if err != nil {
		errorHandler(err)
	}
	defer resp.Body.Close()
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		errorHandler(err)
	} else {
		var healthy healthResp
		if err := json.Unmarshal(body, &healthy); err != nil {
			errorHandler(err)
		}
		renderResponse(c, gin.H{
			"payloadBilboHealth": healthy,
			"version":            render.VersionPage,
			"author":             render.ContactAuthor,
			"project":            projName,
		}, "bilbo/health.tmpl")

	}

	// curl -u elastic:(password)  -H 'Content-Type: application/json' -XGET https://locahost:9243/_cluster/health?pretty healthy check
}

func initES() *es6.Client {
	log.Println("Calling: initES")
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
