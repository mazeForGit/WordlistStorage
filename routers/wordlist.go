package routers

import (
	//"strconv"
	//"fmt"
	
	"github.com/gin-gonic/gin"
	model "github.com/mazeForGit/WordlistStorage/model"
)
func WordListGET(c *gin.Context) {	
	c.JSON(200, model.GlobalWordList)
}
func WordListPUT(c *gin.Context) {
	var s model.ResponseStatus

	var wl model.WordList
	err := c.BindJSON(&wl)
	if err != nil {
		s = model.ResponseStatus{Code: 422, Text: "unprocessable entity"}
		c.JSON(422, s)
		return
	}
	//fmt.Println(wl)
	model.AddWordListToStorage(wl)
	
	s = model.ResponseStatus{Code: 200, Text: "data received"}
	c.JSON(200, s)
}
