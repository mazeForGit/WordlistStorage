package routers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	data "github.com/mazeForGit/WordlistStorage/data"
)
func ConfigGET(c *gin.Context) {
	c.JSON(200, data.GlobalConfig)
}
func ConfigPOST(c *gin.Context) {
	var s data.ResponseStatus
	var err error
	
	fmt.Println("ConfigPOST .. before bind")
	err = c.BindJSON(&data.GlobalConfig)
	if err != nil {
		s = data.ResponseStatus{Code: 422, Text: "unprocessable entity"}
		c.JSON(422, s)
		return
	}
	
	fmt.Println("ConfigPOST .. before read")
	err = data.ReadGlobalWordlistFromRemote()
	if err != nil {
		s = data.ResponseStatus{Code: 422, Text: "can't read global wordlist"}
		c.JSON(200, s)
		return
	}
	
	fmt.Println("ConfigPOST .. ok")
	s = data.ResponseStatus{Code: 200, Text: "got global wordlist"}
	c.JSON(200, s)
}
