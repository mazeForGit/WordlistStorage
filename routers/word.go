package routers

import (
	//"strconv"
	//"fmt"
	
	"github.com/gin-gonic/gin"
	data "github.com/mazeForGit/WordlistStorage/data"
)
func WordGET(c *gin.Context) {	
	var s data.ResponseStatus
	
	var vars map[string][]string
	vars = c.Request.URL.Query()
	var size string = "short"
	var test string = "Big Five"
	
	if _, ok := vars["size"]; ok {
		size = c.Request.URL.Query().Get("size")
	}
	if _, ok := vars["test"]; ok {
		test = c.Request.URL.Query().Get("test")
	}
	
	wrl, err := data.GetWordList(test, size)
	
	if err != nil {
		s = data.ResponseStatus{Code: 422, Text: "unprocessable entity"}
		c.JSON(422, s)
		return
	}
	
	c.JSON(200, wrl)
}

func WordPOST(c *gin.Context) {
	var s data.ResponseStatus
	
	var vars map[string][]string
	vars = c.Request.URL.Query()
	var domain string = ""
	
	if _, ok := vars["domain"]; ok {
		domain = c.Request.URL.Query().Get("domain")
	}
	
	var wrds []data.Word
	err := c.BindJSON(&wrds)
	if err != nil {
		s = data.ResponseStatus{Code: 422, Text: "unprocessable entity"}
		c.JSON(422, s)
		return
	}
	//fmt.Println(wrds)
	data.AddWordsToStorage(domain, wrds)
	
	s = data.ResponseStatus{Code: 200, Text: "data received"}
	c.JSON(200, s)
}