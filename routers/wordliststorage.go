package routers

import (
	//"strconv"
	//"fmt"
	
	"github.com/gin-gonic/gin"
	data "github.com/mazeForGit/WordlistStorage/data"
)

func WordListStorageGET(c *gin.Context) {
	c.JSON(200, data.GlobalWordListStorage)
}
