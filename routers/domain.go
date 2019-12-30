package routers

import (
	//"strconv"
	//"fmt"
	
	"github.com/gin-gonic/gin"
	data "github.com/mazeForGit/WordlistStorage/data"
)
func DomainGET(c *gin.Context) {	
	c.JSON(200, data.GetDomains())
}