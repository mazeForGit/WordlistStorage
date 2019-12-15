package routers

import (
	//"strconv"
	//"fmt"
	
	"github.com/gin-gonic/gin"
	data "github.com/mazeForGit/WordlistStorage/data"
)

func ResultGET(c *gin.Context) {
	var s data.ResponseStatus
	var vars map[string][]string
	vars = c.Request.URL.Query()
	var test string = ""
	var domain string = ""
	
	if _, ok := vars["test"]; ok {
		test = c.Request.URL.Query().Get("test")
	} 
	if _, ok := vars["domain"]; ok {
		domain = c.Request.URL.Query().Get("domain")
	}
	
	if test != "" && domain != "" {
		c.JSON(200, data.ResultOnSession(test, domain))
	} else {
		s = data.ResponseStatus{Code: 422, Text: "missing data"}
		c.JSON(422, s)
	}
}
