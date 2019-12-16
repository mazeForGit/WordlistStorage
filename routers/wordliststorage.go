package routers

import (
	"strconv"
	//"fmt"
	
	"github.com/gin-gonic/gin"
	data "github.com/mazeForGit/WordlistStorage/data"
)

func WordListStorageGET(c *gin.Context) {
	c.JSON(200, data.GlobalWordListStorage)
}
func WordListStoragePUT(c *gin.Context) {
	var s data.ResponseStatus
	
	removeCount := strconv.Itoa(len(data.GlobalWordListStorage))
	err := c.BindJSON(&data.GlobalWordListStorage)
	if err != nil {
		s = data.ResponseStatus{Code: 422, Text: "unprocessable entity"}
		c.JSON(422, s)
		return
	}
	
	data.RebuildWordListResult()
	replaceCount := strconv.Itoa(len(data.GlobalWordListStorage))
	s = data.ResponseStatus{Code: 200, Text: "replaced items = " + removeCount + " by " + replaceCount}
	c.JSON(200, s)
}
