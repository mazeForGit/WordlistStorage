package routers

import (

	"github.com/gin-gonic/gin"
	data "github.com/mazeForGit/WordlistStorage/model"
)

func ResultGET(c *gin.Context) {

	c.Header("Content-Type", "application/json")
		
	var s model.ResponseStatus
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
	
	if test == "" && domain != "" {
		wl, err := model.ResultOnSession(domain)
		if err != nil {
			s = model.ResponseStatus{Code: 422, Text: "unknwon domain = " + domain}
			c.JSON(422, s)
			return
		}
		
		c.JSON(200, wl)
	} else if test != "" && domain != "" {
		r, err := model.ResultOnSessionByTest(test, domain)
		if err != nil {
			s = model.ResponseStatus{Code: 422, Text: "unknwon domain = " + domain}
			c.JSON(422, s)
			return
		}
		
		c.JSON(200, r)
	} else {
		s = model.ResponseStatus{Code: 422, Text: "missing data"}
		c.JSON(422, s)
	}
}
