package routers

import (
	"strconv"
	//"fmt"
	
	"github.com/gin-gonic/gin"
	model "github.com/mazeForGit/WordlistStorage/model"
)

func WordListStorageGET(c *gin.Context) {
	c.JSON(200, model.GlobalWordListStorage)
}
func WordListStoragePUT(c *gin.Context) {
	var s model.ResponseStatus
	
	removeCount := strconv.Itoa(len(model.GlobalWordListStorage))
	err := c.BindJSON(&model.GlobalWordListStorage)
	if err != nil {
		s = model.ResponseStatus{Code: 422, Text: "unprocessable entity"}
		c.JSON(422, s)
		return
	}
	
	model.RebuildWordListResult()
	replaceCount := strconv.Itoa(len(model.GlobalWordListStorage))
	s = model.ResponseStatus{Code: 200, Text: "replaced items = " + removeCount + " by " + replaceCount}
	c.JSON(200, s)
}
