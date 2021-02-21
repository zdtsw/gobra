package main
// functions for "jenkins"

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

type jenkinsMaster struct {
	Project  	string   	`json:"project"`
	Instance	int			`json:"number"`
	Studio		string		`json:"studio"`
}
var allJenkins = []jenkinsMaster{
	{Instance: 4, Project: "Kin", Studio: "DICE"},
	{Instance: 1, Project: "wal", Studio: "DICE"},
	{Instance: 1, Project: "cas", Studio: "DICE"},
	{Instance: 3, Project: "dun", Studio: "DICE"},
	{Instance: 1, Project: "fb1", Studio: "Frostbite"},
	{Instance: 1, Project: "fb2021", Studio: "Frostbite"},
	{Instance: 2, Project: "Exc", Studio: "Critiron"},
  }

// routine functions for /jenkins/info
func jenkinsInfoHandler(c *gin.Context) {
	allJenkins := getAllJenkins()
	log.Println("Calling: jenkinsInfoHandler")
	log.Println("Load page in path: " + c.Request.URL.Path)
	c.HTML(http.StatusOK, "jenkins/info.tmpl", gin.H{
		"payload": allJenkins,
		"version": render.VersionPage,
		"author":  render.ContactAuthor,
	})
}

// routine functions for /jenkins/project/:project
func jenkinsMainHandler(c *gin.Context) {
	log.Println("Calling: jenkinsMainHandler")
	log.Println("Load page in path: " + c.Request.URL.Path)
	c.HTML(http.StatusOK, "jenkins/main.tmpl", gin.H{
		"version": render.VersionPage,
		"author":  render.ContactAuthor,
		"project": render.EAProject,
	})
}

// Return a list of all the articles
func getAllJenkins() []jenkinsMaster {
	return allJenkins
  }