package routers

import (
	//"strconv"
	//"fmt"
	
	"github.com/gin-gonic/gin"
	data "github.com/mazeForGit/WordlistStorage/data"
)
func WordListGET(c *gin.Context) {	
	c.JSON(200, data.GlobalWordList)
}
func WordListPUT(c *gin.Context) {
	var s data.ResponseStatus

	var wl data.WordList
	err := c.BindJSON(&wl)
	if err != nil {
		s = data.ResponseStatus{Code: 422, Text: "unprocessable entity"}
		c.JSON(422, s)
		return
	}
	//fmt.Println(wl)
	data.AddWordListToStorage(wl)
	
	s = data.ResponseStatus{Code: 200, Text: "data received"}
	c.JSON(200, s)
}
