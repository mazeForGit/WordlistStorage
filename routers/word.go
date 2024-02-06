package routers

import (
	//"strconv"
	//"fmt"
	
	"github.com/gin-gonic/gin"
	model "github.com/mazeForGit/WordlistStorage/model"
)
func WordGET(c *gin.Context) {	
	var s model.ResponseStatus
	
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
	
	wrl, err := model.GetWordList(test, size)
	
	if err != nil {
		s = model.ResponseStatus{Code: 422, Text: "unprocessable entity"}
		c.JSON(422, s)
		return
	}
	
	c.JSON(200, wrl)
}

func WordPOST(c *gin.Context) {
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
	
	var wrds []model.Word
	err := c.BindJSON(&wrds)
	if err != nil {
		s = model.ResponseStatus{Code: 422, Text: "unprocessable entity"}
		c.JSON(422, s)
		return
	}
	//fmt.Println(wrds)
	model.AddWordsToStorage(test, domain, wrds)
	
	s = model.ResponseStatus{Code: 200, Text: "data received"}
	c.JSON(200, s)
}